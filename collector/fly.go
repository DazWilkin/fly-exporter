package collector

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/superfly/flyctl/api"
	"github.com/superfly/flyctl/terminal"
)

const (
	url string = "https://api.fly.io"
)

// FlyCollector collects metrics
type FlyCollector struct {
	System System
	Token  string
	Log    logr.Logger
	Apps   *prometheus.Desc
}

// NewFlyCollector returns a new FlyCollector
func NewFlyCollector(s System, token string, log logr.Logger) *FlyCollector {
	return &FlyCollector{
		System: s,
		Token:  token,
		Log:    log,

		Apps: prometheus.NewDesc(
			prometheus.BuildFQName(s.Namespace, s.Subsystem, "apps"),
			"Total Number of Apps",
			[]string{"id", "name", "org_slug", "status", "deployed"},
			nil,
		),
	}
}

// Collect implements Prometheus' Collector interface and is used to collect metrics
func (c *FlyCollector) Collect(ch chan<- prometheus.Metric) {
	log := c.Log.WithName("Collect")

	api.SetBaseURL(url)
	name := fmt.Sprintf("%s_%s", c.System.Namespace, c.System.Subsystem)
	client := api.NewClient(c.Token, name, c.System.Version, terminal.DefaultLogger)

	ctx := context.Background()
	role := ""
	apps, err := client.GetApps(ctx, &role)
	if err != nil {
		log.Error(err, "unable to get apps")
	}

	log.Info("Retrieved apps",
		"number", len(apps),
	)

	for _, app := range apps {
		log.Info("Details",
			"app", app,
		)
		ch <- prometheus.MustNewConstMetric(
			c.Apps,
			prometheus.CounterValue,
			1.0,
			app.ID, app.Name, app.Organization.Slug, app.Status, strconv.FormatBool(app.Deployed),
		)
	}

}

// Describe implements Prometheus' Collector interface and is used to describe metrics
func (c *FlyCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Apps
}
