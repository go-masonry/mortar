package monitoring

import (
	"github.com/go-masonry/mortar/interfaces/monitor"
)

type mortarMetrics struct{}

// Counter creates a counter with possible predefined tags
func (mm *mortarMetrics) Counter(name, desc string) monitor.TagsAwareCounter {
	panic("implement me")
}

// Gauge creates a gauge with possible predefined tags
func (mm *mortarMetrics) Gauge(name, desc string) monitor.TagsAwareGauge {
	panic("implement me")
}

// Histogram creates a histogram with possible predefined tags
func (mm *mortarMetrics) Histogram(name, desc string, buckets monitor.Buckets) monitor.TagsAwareHistogram {
	panic("implement me")
}

// Timer creates a timer with possible predefined tags
func (mm *mortarMetrics) Timer(name, desc string) monitor.TagsAwareTimer {
	panic("implement me")
}

// WithTags sets custom tags to be included if possible in every Metric
func (mm *mortarMetrics) WithTags(tags monitor.Tags) monitor.Metrics {
	panic("implement me")
}
