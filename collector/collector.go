package collector

import (
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/process" // https://godoc.org/github.com/shirou/gopsutil/process
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type User struct {
	Kind          string `json:kind`
	Name          string `json:name`
	Admin         string
	Groups        string
	Server        string
	Pending       string
	Last_activity string
}

type Process struct {
	Pid        int32
	Username   string
	CPUPercent float64
	RSS        uint64
}

func FetchUserList(parameters map[string]string) []User {
	userApiUrl := fmt.Sprintf("%s/users", parameters["apiUrl"])
	req, _ := http.NewRequest("GET", userApiUrl, nil)
	req.Header.Set("Authorization", "token "+parameters["apiToken"])
	client := new(http.Client)
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	var users []User
	json.Unmarshal(byteArray, &users)
	return users
}

func FetchProcessInfoList() []Process {
	pids, err := process.Pids()
	if err != nil {
		return []Process{}
	}
	var processList []Process
	for _, pid := range pids {
		p, _ := process.NewProcess(pid)
		user, _ := p.Username()
		cpu, _ := p.CPUPercent()
		mem, err := p.MemoryInfo()
		rss := uint64(0)
		if err == nil {
			rss = mem.RSS
		}
		isRunning, _ := p.IsRunning()
		if isRunning {
			processList = append(processList, Process{
				Pid:        pid,
				Username:   user,
				CPUPercent: cpu,
				RSS:        rss,
			})
		}
	}
	return processList
}

func FetchProcessCount(user User, processes []Process, parameters map[string]string, ch chan<- float64) {
	count := 0
	for _, p := range processes {
		if p.Username == user.Name {
			count++
		}
	}
	ch <- float64(count)
}

func FetchCpuUsage(user User, processes []Process, parameters map[string]string, ch chan<- float64) {
	cpu := float64(0)
	for _, p := range processes {
		if p.Username == user.Name {
			cpu += p.CPUPercent
		}
	}
	ch <- cpu
}

func FetchMemoryUsage(user User, processes []Process, parameters map[string]string, ch chan<- float64) {
	mem := uint64(0)
	for _, p := range processes {
		if p.Username == user.Name {
			mem += p.RSS
		}
	}
	ch <- float64(mem)
}

func FetchDiskUsage(user User, parameters map[string]string, ch chan<- float64) {
	userPath := fmt.Sprintf("%s/%s", parameters["notebookDir"], user.Name)
	exist, _ := exists(userPath)
	if exist {
		size, _ := dirSize(userPath)
		ch <- float64(size)
	} else {
		ch <- 0
	}
}

func dirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
