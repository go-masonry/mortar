package monitor

import (
	"context"
	"time"
)

//go:generate mockgen -source=interfaces.go -destination=mock/mock.go

// Buckets must be ordered in increasing order
type Buckets []float64

// Tags or Labels, key->value
type Tags map[string]string

// TagsAwareCounter defines a counter with the ability override tags value either explicitly or from context (by extractors)
type TagsAwareCounter interface {
	Counter
	WithTags(tags Tags) TagsAwareCounter
	WithContext(ctx context.Context) TagsAwareCounter
}

// Counter is a Metric that represents a single numerical value that only ever goes up
type Counter interface {
	// Inc increments the counter by 1
	Inc()
	// Add adds the given value to the counter, negative values are not advised
	Add(v float64)
}

// TagsAwareGauge defines a gauge with the ability to override tags value either explicitly or from context (by extractors)
type TagsAwareGauge interface {
	Gauge
	WithTags(tags Tags) TagsAwareGauge
	WithContext(ctx context.Context) TagsAwareGauge
}

//Gauge is a Metric that represents a single numerical value that can arbitrarily go up and down
type Gauge interface {
	// Set sets Gauge value
	Set(v float64)
	// Add adds (or substracts if negative) from previously set value
	Add(v float64)
	// Inc adds 1
	Inc()
	// Dec adds -1
	Dec()
}

// TagsAwareHistogram defines a histogram with the ability to override tags value either explicitly or from context (by extractors)
type TagsAwareHistogram interface {
	Histogram
	WithTags(tags Tags) TagsAwareHistogram
	WithContext(ctx context.Context) TagsAwareHistogram
}

// A Histogram counts individual observations from an event or sample stream in configurable buckets
type Histogram interface {
	// Record value
	Record(v float64)
}

// Timer is used to track duration of operations
type Timer interface {
	Record(d time.Duration)
}

// --- Auxiliary types ---

// TagsAwareTimer defines a timer with the ability to override tags value either explicitly or from context (by extractors)
type TagsAwareTimer interface {
	Timer
	WithTags(tags Tags) TagsAwareTimer
	WithContext(ctx context.Context) TagsAwareTimer
}

// Metrics defines various monitoring capabilities
//
// It is expected that each Metric is unique, uniqueness is calculated by combining
// 	- name
//	- tag key names
type Metrics interface {
	// Counter creates a counter with possible predefined tags
	Counter(name, desc string) TagsAwareCounter
	// Gauge creates a gauge with possible predefined tags
	Gauge(name, desc string) TagsAwareGauge
	// Histogram creates a histogram with possible predefined tags
	Histogram(name, desc string, buckets Buckets) TagsAwareHistogram
	// Timer creates a timer with possible predefined tags
	Timer(name, desc string) TagsAwareTimer
	// WithTags sets custom tags to be included if possible in every Metric
	WithTags(tags Tags) Metrics
}

// Reporter defines Metrics reporter with Connect/Close options on demand and not on creation
type Reporter interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	Metrics() Metrics
}

// ContextExtractor is a function that will extract values from the context and return them as Tags
// Make sure that this function returns fast and is "thread safe"
type ContextExtractor func(ctx context.Context) Tags
