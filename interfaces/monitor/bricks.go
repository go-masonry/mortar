package monitor

import (
	"context"
)

//go:generate mockgen -source=interfaces.go -destination=mock/mock.go

// ### Bricks* interface
//
// All Brick* interfaces are to be implemented by the brick library for example
// 	- Prometheus
//	- Datadog
//
// This allows a better seperation of responsibility and also adds clarity.
//
// Within mortar we want to be able extract tag values from context, but the implementation shouldn't be
// aware of this. The only thing it needs to implement is a support for monitor.Tags

// BrickMetric is a marker interface
type BrickMetric interface{}

// BricksCounter defines a counter to be implemented by external wrapper
type BricksCounter interface {
	BrickMetric
	WithTags(tags map[string]string) Counter
}

// BricksGauge defines a gauge to be implemented by external wrapper
type BricksGauge interface {
	BrickMetric
	WithTags(tags map[string]string) Gauge
}

// BricksHistogram defines a histogram to be implemented by external wrapper
type BricksHistogram interface {
	BrickMetric
	WithTags(tags map[string]string) Histogram
}

// BricksMetrics defines various monitoring capabilities to be implemented by external wrapper
type BricksMetrics interface {
	// Counter creates a counter with predefined tag key names.
	// This will allow to set their values right before using Counter methods
	Counter(name, desc string, tagKeys ...string) (BricksCounter, error)
	// Gauge creates a gauge with predefined tag key names.
	// This will allow to set their values right before using Gauge methods
	Gauge(name, desc string, tagKeys ...string) (BricksGauge, error)
	// Histogram creates a histogram with predefined tag key names.
	// This will allow to set their values right before using Histogram methods
	Histogram(name, desc string, buckets []float64, tagKeys ...string) (BricksHistogram, error)
	// Remove removes this metric from external registry, if applicable
	Remove(metric BrickMetric) error
}

// BricksReporter defines Metrics reporter  to be implemented by external wrapper
type BricksReporter interface {
	// Connect, if applicable connect to the agent only when this function is called
	Connect(ctx context.Context) error
	// Close, if applicable disconnect/close connection to the agent
	Close(ctx context.Context) error
	// Metrics returns implementation of BricksMetrics
	Metrics() BricksMetrics
}

// Builder defines a simple BricksReporter builder
type Builder interface {
	Build() BricksReporter
}
