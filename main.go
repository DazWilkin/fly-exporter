package main

import (
	"flag"
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
	endpoint = flag.String("endpoint", "0.0.0.0:8080", "The endpoint of the HTTP server")
)
var (
	log logr.Logger
)
var (
	token = os.Getenv("TOKEN")
)

func handleHealthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
func handleRoot(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<a href=/metrics>metrics</a>"))
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
	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	log.Info("Server starting",
		"endpoint", *endpoint,
	)
	log.Error(http.ListenAndServe(*endpoint, mux), "unable to start server")
}
