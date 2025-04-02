package main

import (
	"flag"
	"html/template"
	stdlog "log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/DazWilkin/fly-exporter/collector"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	namespace string = "fly"
	subsystem string = "exporter"
	version   string = "v0.0.1"
)

const (
	rootTemplate string = `
{{- define "content" }}
<!DOCTYPE html>
<html lang="en-US">
<head>
<title>Prometheus Exporter for Fly.io</title>
<style>
body {
  font-family: Verdana;
}
</style>
</head>
<body>
	<h2>Prometheus Exporter for Fly.io</h2>
	<hr/>
	<ul>
	<li><a href="{{ .MetricsPath }}">metrics</a></li>
	<li><a href="/healthz">healthz</a></li>
	</ul>
</body>
</html>
{{- end}}
`
)

var (
	// GitCommit is the git commit value and is expected to be set during build
	GitCommit string
	// GoVersion is the Golang runtime version
	GoVersion = runtime.Version()
	// OSVersion is the OS version (uname --kernel-release) and is expected to be set during build
	OSVersion string
	// StartTime is the start time of the exporter represented as a UNIX epoch
	StartTime = time.Now().Unix()
)
var (
	endpoint    = flag.String("endpoint", "0.0.0.0:8080", "The endpoint of the HTTP server")
	metricsPath = flag.String("path", "/metrics", "The path on which Prometheus metrics will be served")
)
var (
	log logr.Logger
)
var (
	token = os.Getenv("TOKEN")
)

type Content struct {
	MetricsPath string
}

func handleHealthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("ok")); err != nil {
		log.Error(err, "unable to write response")
	}
}
func handleRoot(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	t := template.Must(template.New("content").Parse(rootTemplate))
	if err := t.ExecuteTemplate(w, "content", Content{MetricsPath: *metricsPath}); err != nil {
		log.Error(err, "unable to execute template")
	}
}
func main() {
	log = stdr.NewWithOptions(stdlog.New(os.Stderr, "", stdlog.LstdFlags), stdr.Options{LogCaller: stdr.All})
	log = log.WithName("main")

	if token == "" {
		msg := "environment variable `TOKEN` is required (use `flyctl auth token`)"
		log.Info(msg)
		panic(msg)
	}

	flag.Parse()
	if *endpoint == "" {
		msg := "expected flag `--endpoint"
		log.Info(msg)
		panic(msg)
	}

	registry := prometheus.NewRegistry()

	s := collector.System{
		Namespace: namespace,
		Subsystem: subsystem,
		Version:   version,
	}

	b := collector.Build{
		OsVersion: OSVersion,
		GoVersion: GoVersion,
		GitCommit: GitCommit,
		StartTime: StartTime,
	}

	registry.MustRegister(collector.NewExporterCollector(s, b, log))
	registry.MustRegister(collector.NewFlyCollector(s, token, log))

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(handleRoot))
	mux.Handle("/healthz", http.HandlerFunc(handleHealthz))
	mux.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	log.Info("Server starting",
		"endpoint", *endpoint,
	)
	log.Error(http.ListenAndServe(*endpoint, mux), "unable to start server")
}
