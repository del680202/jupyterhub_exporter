

This is a prometheus exporter for monitoring jupyterHub metrics.

This exporter use os command to fetch user metrics, it need to be installed on jupyterHub server.

It provide metrics as below

|Metrics Name|Description|
|:---|:---|
|user_total|Total users in jupyterhub database|
|process_count|Process number per each user|
|cpu_usage|CPU usage per each jupyterhub user|
|memory_usage|Memory usage per each jupyterhub user|
|disk_usage|Disk usage per each jupyterhub user|

# Parameter

|Name|Default|Description|
|:---|:---|:---|
|--web.listen-address|:9527|Listen port|
|--web.telemetry-path|/metrics|Prometheus endpoint|
|--jupyter.api-token|*(Required)*|JupyterHub REST API admin token|
|--jupyter.api-url|http://127.0.0.1:8081/hub/api|JupyterHub REST API URL|
|--jupyter.notebook-dir|/home|Jupyter notebook root directory|
|-h|X|Help Message|

# Prerequisites

* Go 1.8+
* CentOS 7

# Prepare

Setup Go environment

```
$ yum install go

$ mkdir $HOME/go
$ vim ~/.bashrc
...
export GOROOT=/usr/lib/golang
export GOPATH=$HOME/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

$ source ~/.bashrc
```

# Deploy

```
# Download source code
$ env GIT_TERMINAL_PROMPT=1 go get github.com/del680202/jupyterhub_exporter
$ cd $GOPATH/src/github.com/del680202/jupyterhub_exporter

# Install go building tool
$ go get github.com/mitchellh/gox 
$ go get -v github.com/golang/dep
$ go install -v github.com/golang/dep/cmd/dep

# Download dependency
$ dep init
$ dep ensure
```

# Test

```
$ go test
```

# Test Run

```
$ go run main.go --jupyter.api-token=YOUR_API_TOKEN --jupyter.notebook-dir=YOUR_NOTEBOOK_HOME
```

# Build
```
gox --osarch "linux/amd64"  --output release/jupyterhub_exporter
```

# Run
```
./release/jupyterhub_exporter --jupyter.api-token=YOUR_API_TOKEN --jupyter.notebook-dir=YOUR_NOTEBOOK_HOME
```
