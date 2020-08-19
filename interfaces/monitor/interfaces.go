package monitor

import (
	"context"
	"time"
)

//go:generate mockgen -source=interfaces.go -destination=mock/mock.go

// Metrics defines various monitoring capabilities
type Metrics interface {
	// Counter returns the Counter object corresponding to the name.
	Counter(ctx context.Context, name string) Counter
	// Gauge returns the Gauge object corresponding to the name.
	Gauge(ctx context.Context, name string) Gauge
	// Timer returns the Timer object corresponding to the name.
	Timer(ctx context.Context, name string) Timer
	// Histogram returns the Histogram object corresponding to the name.
	Histogram(ctx context.Context, name string) Histogram
	// AddTag adds a tag to the metric.
	AddTag(name, value string) Metrics
	// Implementation returns the actual lib/struct that is responsible for the above logic
	Implementation() interface{}
}

// Counter is the interface for emitting counter type metrics.
type Counter interface {
	// Inc increments the counter by a delta.
	Inc(delta int64) error
}

// Gauge is the interface for emitting gauge metrics.
type Gauge interface {
	// Update sets the gauges absolute value.
	Update(value float64) error
}

// Timer is the interface for emitting timer metrics.
type Timer interface {
	// Record a specific duration directly.
	Record(value time.Duration) error
	// Start gives you back a specific point in time to report via Stop.
	Start() Stopwatch
}

type Stopwatch interface {
	// Stop reports time elapsed since the stopwatch start to the recorder.
	Stop() error
}

// Histogram is the interface for emitting histogram metrics
type Histogram interface {
	// RecordValue records a specific value directly.
	// Will use the configured value buckets for the histogram.
	RecordValue(value float64)
	// RecordDuration records a specific duration directly.
	// Will use the configured duration buckets for the histogram.
	RecordDuration(value time.Duration)
	// Start gives you a specific point in time to then record a duration.
	// Will use the configured duration buckets for the histogram.
	Start() Stopwatch
}

// Reporter defines Metrics reporter
type Reporter interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	Metrics() Metrics
}

// ContextExtractor is a function that will extract values from the context and return them as Tags to be added
// Make sure that this function returns fast and is "thread safe"
type ContextExtractor func(ctx context.Context) map[string]string

// Builder defines Monitor builder options
type Builder interface {
	SetAddress(hostPort string) Builder
	SetPrefix(prefix string) Builder
	SetTags(tags map[string]string) Builder
	AddContextExtractors(extractors ...ContextExtractor) Builder
	Build() Reporter
}
