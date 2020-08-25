package monitoring

import (
	"context"
	"sync"
	"time"

	"github.com/go-masonry/mortar/interfaces/monitor"
)

// *******************************************************************
// *                             Counter                             *
// *******************************************************************
type counter struct {
	*tagsMetric
	bricksCounter monitor.BricksCounter
	extractors    []monitor.ContextExtractor
}

func (c *counter) Inc() {
	c.bricksCounter.WithTags(c.tags).Inc()
}

func (c *counter) Add(v float64) {
	c.bricksCounter.WithTags(c.tags).Add(v)
}

func (c *counter) WithTags(tags monitor.Tags) monitor.TagsAwareCounter {
	c.withTags(tags)
	return c
}

func (c *counter) WithContext(ctx context.Context) monitor.TagsAwareCounter {
	c.withContext(ctx, c.extractors)
	return c
}

// *******************************************************************
// *                             Gauge                               *
// *******************************************************************
type gauge struct {
	*tagsMetric
	bricksGauge monitor.BricksGauge
	extractors  []monitor.ContextExtractor
}

// Set sets Gauge value
func (g *gauge) Set(v float64) {
	g.bricksGauge.WithTags(g.tags).Set(v)
}

// Add adds (or substracts if negative) from previously set value
func (g *gauge) Add(v float64) {
	g.bricksGauge.WithTags(g.tags).Add(v)
}

// Inc adds 1
func (g *gauge) Inc() {
	g.bricksGauge.WithTags(g.tags).Inc()
}

// Dec adds -1
func (g *gauge) Dec() {
	g.bricksGauge.WithTags(g.tags).Dec()
}

func (g *gauge) WithTags(tags monitor.Tags) monitor.TagsAwareGauge {
	g.withTags(tags)
	return g
}

func (g *gauge) WithContext(ctx context.Context) monitor.TagsAwareGauge {
	g.withContext(ctx, g.extractors)
	return g
}

// *******************************************************************
// *                             histogram                           *
// *******************************************************************
type histogram struct {
	*tagsMetric
	bricksHistogram monitor.BricksHistogram
	extractors      []monitor.ContextExtractor
}

// Record value
func (h *histogram) Record(v float64) {
	h.bricksHistogram.WithTags(h.tags).Record(v)
}

func (h *histogram) WithTags(tags monitor.Tags) monitor.TagsAwareHistogram {
	h.withTags(tags)
	return h
}

func (h *histogram) WithContext(ctx context.Context) monitor.TagsAwareHistogram {
	h.withContext(ctx, h.extractors)
	return h
}

// *******************************************************************
// *                             timer                               *
// *******************************************************************
type timer struct {
	*tagsMetric
	bricksHistogram monitor.BricksHistogram
	extractors      []monitor.ContextExtractor
}

// Record uses Histogram to record timed duration
// Since Histogram accepts float64 we will take the d.Seconds() which returns float64
func (t *timer) Record(d time.Duration) {
	t.bricksHistogram.WithTags(t.tags).Record(d.Seconds())
}

func (t *timer) WithTags(tags monitor.Tags) monitor.TagsAwareTimer {
	t.withTags(tags)
	return t
}

func (t *timer) WithContext(ctx context.Context) monitor.TagsAwareTimer {
	t.withContext(ctx, t.extractors)
	return t
}

// *******************************************************************
// *                             tags helper                         *
// *******************************************************************
type tagsMetric struct {
	sync.Mutex
	tags   monitor.Tags
	copied bool
}

func (tm *tagsMetric) withTags(tags monitor.Tags) {
	tm.Lock()
	defer tm.Unlock()
	if !tm.copied {
		var tagsCopy = monitor.Tags{}
		for k, v := range tm.tags {
			tagsCopy[k] = v
		}
		tm.tags = tagsCopy
		tm.copied = true
	}
	for k, v := range tags {
		tm.tags[k] = v
	}
}

func (tm *tagsMetric) withContext(ctx context.Context, extractors []monitor.ContextExtractor) {
	for _, extractor := range extractors {
		extractedTags := extractor(ctx)
		tm.withTags(extractedTags)
	}
}

// Metric Constructors

func newCounterWithTags(bricksCounter monitor.BricksCounter, predefinedTags monitor.Tags, extractors []monitor.ContextExtractor) monitor.TagsAwareCounter {
	return &counter{
		tagsMetric:    &tagsMetric{tags: predefinedTags},
		bricksCounter: bricksCounter,
		extractors:    extractors,
	}
}

func newGaugeWithTags(bricksGauge monitor.BricksGauge, predefinedTags monitor.Tags, extractors []monitor.ContextExtractor) monitor.TagsAwareGauge {
	return &gauge{
		tagsMetric:  &tagsMetric{tags: predefinedTags},
		bricksGauge: bricksGauge,
		extractors:  extractors,
	}
}

func newHistogramWithTags(bricksHistogram monitor.BricksHistogram, predefinedTags monitor.Tags, extractors []monitor.ContextExtractor) monitor.TagsAwareHistogram {
	return &histogram{
		tagsMetric:      &tagsMetric{tags: predefinedTags},
		bricksHistogram: bricksHistogram,
		extractors:      extractors,
	}
}

func newTimerWithTags(bricksHistogram monitor.BricksHistogram, predefinedTags monitor.Tags, extractors []monitor.ContextExtractor) monitor.TagsAwareTimer {
	return &timer{
		tagsMetric:      &tagsMetric{tags: predefinedTags},
		bricksHistogram: bricksHistogram,
		extractors:      extractors,
	}
}
