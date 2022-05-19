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
	namespace string = "fly"
	subsystem string = "exporter"
	version   string = "v0.0.1"
)
const (
	url string = "https://api.fly.io"
)

var (
	name string = fmt.Sprintf("%s_%s", namespace, subsystem)
)

// FlyCollector collects metrics
type FlyCollector struct {
	Token string
	Log   logr.Logger

	Count *prometheus.Desc
}

// NewFlyCollector returns a new FlyCollector
func NewFlyCollector(token string, log logr.Logger) *FlyCollector {
	return &FlyCollector{
		Token: token,
		Log:   log,

		Count: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "apps"),
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
	client := api.NewClient(c.Token, name, version, terminal.DefaultLogger)

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
			c.Count,
			prometheus.CounterValue,
			1.0,
			app.ID, app.Name, app.Organization.Slug, app.Status, strconv.FormatBool(app.Deployed),
		)
	}

}

// Describe implements Prometheus' Collector interface and is used to describe metrics
func (c *FlyCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Count
}
