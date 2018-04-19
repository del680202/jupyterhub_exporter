package exporter

import (
	. "git.rakuten-it.com/DSD/jupyterhub_exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
)

type Exporter struct {
	Parameters   map[string]string
	Metrics      map[string]prometheus.Gauge
	LabelMetrics map[string]*prometheus.GaugeVec
}

func NewExporter(namespace string, parameters map[string]string) *Exporter {
	metrics := make(map[string]prometheus.Gauge)
	labelMetrics := make(map[string]*prometheus.GaugeVec)

	metrics["user_total"] = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "user_total",
		Help:      "Total user number in jupyterhub database"})

	labelMetrics["process_count"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "process_count",
		Help:      "Process number per each jupyterhub user"},
		[]string{"user"})

	labelMetrics["cpu_usage"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "cpu_usage",
		Help:      "CPU usage per each jupyterhub user"},
		[]string{"user"})

	labelMetrics["memory_usage"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "memory_usage",
		Help:      "Memory per each jupyterhub user"},
		[]string{"user"})

	labelMetrics["disk_usage"] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "disk_usage",
		Help:      "Disk usage per each jupyterhub user"},
		[]string{"user"})

	return &Exporter{
		Parameters:   parameters,
		Metrics:      metrics,
		LabelMetrics: labelMetrics,
	}
}

type UserResult struct {
	name         string
	processCount float64
	cpuUsage     float64
	memoryUsage  float64
	diskUsage    float64
}

func collectJupterHubMetrics(e *Exporter) {

	users := FetchUserList(e.Parameters)
	processes := FetchProcessInfoList()
	e.Metrics["user_total"].Set(float64(len(users)))

	jobLength := len(users)
	jobs := make(chan UserResult, jobLength)

	for _, user := range users {
		go func(user User, e *Exporter, result chan UserResult) {
			processCount := make(chan float64)
			cpuUsage := make(chan float64)
			memoryUsage := make(chan float64)
			diskUsage := make(chan float64)
			go FetchProcessCount(user, processes, e.Parameters, processCount)
			go FetchCpuUsage(user, processes, e.Parameters, cpuUsage)
			go FetchMemoryUsage(user, processes, e.Parameters, memoryUsage)
			go FetchDiskUsage(user, e.Parameters, diskUsage)

			result <- UserResult{
				name:         user.Name,
				processCount: <-processCount,
				cpuUsage:     <-cpuUsage,
				memoryUsage:  <-memoryUsage,
				diskUsage:    <-diskUsage,
			}
		}(user, e, jobs)
	}
	//waiting for all jobs done
	for i := 0; i < jobLength; i++ {
		userResult := <-jobs
		e.LabelMetrics["process_count"].WithLabelValues(userResult.name).Set(userResult.processCount)
		e.LabelMetrics["cpu_usage"].WithLabelValues(userResult.name).Set(userResult.cpuUsage)
		e.LabelMetrics["memory_usage"].WithLabelValues(userResult.name).Set(userResult.memoryUsage)
		e.LabelMetrics["disk_usage"].WithLabelValues(userResult.name).Set(userResult.diskUsage)
	}
}

// Prometheus function: Collect
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	//Collection logic
	collectJupterHubMetrics(e)

	//used for prometheus
	for _, metric := range e.Metrics {
		metric.Collect(ch)
	}
	for _, metric := range e.LabelMetrics {
		metric.Collect(ch)
	}
}

// Prometheus function: Describe
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range e.Metrics {
		metric.Describe(ch)
	}
	for _, metric := range e.LabelMetrics {
		metric.Describe(ch)
	}
}
