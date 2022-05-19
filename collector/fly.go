package collector

import (
	"context"

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
			"apps_count",
			"Number of Apps",
			[]string{"org_name"},
			nil,
		),
	}
}

// Collect implements Prometheus' Collector interface and is used to collect metrics
func (c *FlyCollector) Collect(ch chan<- prometheus.Metric) {
	log := c.Log.WithName("Collect")

	api.SetBaseURL(url)
	client := api.NewClient(c.Token, "foo", "foo", terminal.DefaultLogger)

	ctx := context.Background()
	role := ""
	apps, err := client.GetApps(ctx, &role)
	if err != nil {
		log.Error(err, "unable to get apps")
	}

	for _, app := range apps {
		log.Info(app.ID, app.Name, app.Organization, app.Regions)
		ch <- prometheus.MustNewConstMetric(
			c.Count,
			prometheus.CounterValue,
			1.0,
			app.Organization.Name,
		)
	}

}

// Describe implements Prometheus' Collector interface and is used to describe metrics
func (c *FlyCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Count
}
