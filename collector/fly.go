package collector

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/DazWilkin/fly-exporter/terminal"

	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/superfly/flyctl/api"
)

const (
	url string = "https://api.fly.io"
)

// FlyCollector collects metrics
type FlyCollector struct {
	System System
	Token  string
	Log    logr.Logger

	App  *prometheus.Desc
	Cert *prometheus.Desc
}

// NewFlyCollector returns a new FlyCollector
func NewFlyCollector(s System, token string, log logr.Logger) *FlyCollector {
	return &FlyCollector{
		System: s,
		Token:  token,
		Log:    log,

		App: prometheus.NewDesc(
			prometheus.BuildFQName(s.Namespace, s.Subsystem, "app_info"),
			"Info about Applications",
			[]string{"id", "name", "org_slug", "status", "deployed"},
			nil,
		),
		Cert: prometheus.NewDesc(
			prometheus.BuildFQName(s.Namespace, s.Subsystem, "cert_info"),
			"Info about Certificates",
			[]string{"app_id", "app_name", "status"},
			nil,
		),
	}
}

// Collect implements Prometheus' Collector interface and is used to collect metrics
func (c *FlyCollector) Collect(ch chan<- prometheus.Metric) {
	log := c.Log.WithName("Collect")

	api.SetBaseURL(url)
	name := fmt.Sprintf("%s_%s", c.System.Namespace, c.System.Subsystem)

	// Replaced flyctl/terminal with own implementation
	// flyctl/terminal takes a dependency on flyctl/internal
	// And this caused Modules issues requiring the use of a replace
	// github.com/loadsmart/calver-go => github.com/ndarilek/calver-go v0.0.0-20230710153822-893bbd83a936
	client := api.NewClient(c.Token, name, c.System.Version, terminal.New(c.Log))

	ctx := context.Background()
	role := ""
	apps, err := client.GetApps(ctx, &role)
	if err != nil {
		log.Error(err, "unable to get apps")
		return
	}

	log.Info("Retrieved apps",
		"number", len(apps),
	)

	var wg sync.WaitGroup
	for _, app := range apps {
		wg.Add(1)
		go func(app api.App) {
			defer wg.Done()
			log := log.WithValues("app", app.Name)
			log.Info("Details")
			ch <- prometheus.MustNewConstMetric(
				c.App,
				prometheus.CounterValue,
				1.0,
				app.ID, app.Name, app.Organization.Slug, app.Status, strconv.FormatBool(app.Deployed),
			)

			// Collect app certificates
			certs, err := client.GetAppCertificates(ctx, app.Name)
			if err != nil {
				log.Error(err, "unable to get app certificates")
			}

			log.Info("Retrieved app's certificates",
				"number", len(certs),
			)

			for _, cert := range certs {
				log := log.WithValues("cert", cert.ClientStatus)
				log.Info("Details")
				ch <- prometheus.MustNewConstMetric(
					c.Cert,
					prometheus.CounterValue,
					1.0,
					app.ID, app.Name, cert.ClientStatus,
				)
			}
		}(app)
		wg.Wait()
	}
}

// Describe implements Prometheus' Collector interface and is used to describe metrics
func (c *FlyCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.App
}
