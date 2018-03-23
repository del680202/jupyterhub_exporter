package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"

	"github.com/del680202/jupyterhub_exporter/exporter"
)

const (
	APP_NAME  = "jupyterhub_exporter"
	NAMESPACE = "jupyterhub"
)

func main() {

	var (
		listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9527").String()
		metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
		apiToken      = kingpin.Flag("jupyter.api-token", "Admin token of jupyterhub rest api.").Required().String()
		apiUrl        = kingpin.Flag("jupyter.api-url", "URL of jupyterhub rest api.").Default("http://127.0.0.1:8081/hub/api").String()
		notebookDir   = kingpin.Flag("jupyter.notebook-dir", "notebook home directory").Default("/home").String()
	)
	log.AddFlags(kingpin.CommandLine)
	kingpin.Version(version.Print(APP_NAME))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	parameters := map[string]string{
		"apiToken":    *apiToken,
		"apiUrl":      *apiUrl,
		"notebookDir": *notebookDir,
	}
	fmt.Println("Start Exporter")
	fmt.Println("Parameters:", parameters)

	// Register exporter
	exporter := exporter.NewExporter(NAMESPACE, parameters)
	prometheus.MustRegister(exporter)

	// Launch http service
	http.Handle(*metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>` + APP_NAME + `</title></head>
             <body>
             <h1>` + APP_NAME + `</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})
	fmt.Println(http.ListenAndServe(*listenAddress, nil))
}
