package monitor

import (
	"context"
	"time"
)

//go:generate mockgen -source=interfaces.go -destination=mock/mock.go

// Metrics defines various monitoring capabilities
type Metrics interface {
	// Gauge measures the value of a metric at a particular time.
	Gauge(ctx context.Context, name string, value float64) error
	// Count tracks how many times something happened per second.
	Count(ctx context.Context, name string, value int64) error
	// Histogram tracks the statistical distribution of a set of values on each host.
	Histogram(ctx context.Context, name string, value float64) error
	// Distribution tracks the statistical distribution of a set of values across your infrastructure.
	Distribution(ctx context.Context, name string, value float64) error
	// Decr is just Count of -1
	Decr(ctx context.Context, name string) error
	// Incr is just Count of 1
	Incr(ctx context.Context, name string) error
	// Set counts the number of unique elements in a group. 'value' is an element in a final SET (["one", "two", "three"])
	Set(ctx context.Context, name string, value string) error
	// Timing sends timing information
	Timing(ctx context.Context, name string, value time.Duration) error
	// Add custom tags for this metric
	AddTag(name, value string) Metrics
	// Set custom rate for this metric
	SetRate(rate float64) Metrics
	// Implementation returns the actual lib/struct that is responsible for the above logic
	Implementation() interface{}
}

type Reporter interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	Metrics() Metrics
}

// ContextExtractor is a function that will extract values from the context and return them as Tags to be added
// Make sure that this function returns fast and is "thread safe"
type ContextExtractor func(ctx context.Context) map[string]string

type Builder interface {
	SetAddress(hostPort string) Builder
	SetPrefix(prefix string) Builder
	SetTags(tags map[string]string) Builder
	AddContextExtractors(extractors ...ContextExtractor) Builder
	Build() Reporter
}
