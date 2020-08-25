package monitoring

import (
	"context"

	"github.com/go-masonry/mortar/interfaces/monitor"
)

type mortarReporter struct {
	externalMetrics monitor.BricksMetrics
	cfg             *monitorConfig
}

// NewMortarReporter creates a new mortar monitoring reporter which is a wrapper to support
// 	- ContextExtractors
// 	- Default Tags, for example: {"version":"v1.0.1", "service":"awesome"}
//
// Meaning, it is possible to also extract tag values from the context, this is useful when the value is set per request/call within the context.Context:
// 	- Canary release https://martinfowler.com/bliki/CanaryRelease.html identifier
// 	- Authentication Token values, but avoid using high cardinality values such as UserID
//
func newMortarReporter(cfg *monitorConfig) monitor.Reporter {
	return &mortarReporter{
		externalMetrics: cfg.reporter.Metrics(),
		cfg:             cfg,
	}
}

func (r *mortarReporter) Connect(ctx context.Context) error {
	return r.cfg.reporter.Connect(ctx)
}

func (r *mortarReporter) Close(ctx context.Context) error {
	return r.cfg.reporter.Close(ctx)
}

func (r *mortarReporter) Metrics() monitor.Metrics {
	return r
}

// Counter creates a counter with possible predefined tags
func (r *mortarReporter) Counter(name string, desc string) monitor.TagsAwareCounter {
	return newMetric(r.externalMetrics, r.cfg).WithTags(r.cfg.tags).Counter(name, desc)
}

// Gauge creates a gauge with possible predefined tags
func (r *mortarReporter) Gauge(name string, desc string) monitor.TagsAwareGauge {
	return newMetric(r.externalMetrics, r.cfg).WithTags(r.cfg.tags).Gauge(name, desc)
}

// Histogram creates a histogram with possible predefined tags
func (r *mortarReporter) Histogram(name string, desc string, buckets monitor.Buckets) monitor.TagsAwareHistogram {
	return newMetric(r.externalMetrics, r.cfg).WithTags(r.cfg.tags).Histogram(name, desc, buckets)
}

// Timer creates a timer with possible predefined tags
func (r *mortarReporter) Timer(name string, desc string, buckets monitor.Buckets) monitor.TagsAwareTimer {
	return newMetric(r.externalMetrics, r.cfg).WithTags(r.cfg.tags).Timer(name, desc, buckets)
}

// WithTags sets custom tags to be included if possible in every Metric
func (r *mortarReporter) WithTags(tags monitor.Tags) monitor.Metrics {
	return newMetric(r.externalMetrics, r.cfg).
		WithTags(r.cfg.tags). // first apply default tags
		WithTags(tags)        // then apply custom ones
}
