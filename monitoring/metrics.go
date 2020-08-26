package monitoring

import (
	"fmt"
	"sort"

	"github.com/go-masonry/mortar/interfaces/monitor"
)

type mortarMetric struct {
	*tagsMetric
	registry *externalRegistry
	cfg      *monitorConfig
}

func newMetric(registry *externalRegistry, cfg *monitorConfig) monitor.Metrics {
	return &mortarMetric{
		registry:   registry,
		cfg:        cfg,
		tagsMetric: &tagsMetric{tags: cfg.tags},
	}
}

// Counter creates a counter with possible predefined tags
func (mm *mortarMetric) Counter(name, desc string) monitor.TagsAwareCounter {
	bricksCounter, err := mm.registry.loadOrStoreCounter(name, desc, mm.extractTagKeys()...)
	if err != nil {
		mm.cfg.onError(fmt.Errorf("error registering counter [%s:%s] metric with %v tags, %w", name, desc, mm.extractTagKeys(), err))
		bricksCounter = newNoopCounter(name, desc, err, mm.cfg.onError)
	}

	return newCounterWithTags(bricksCounter, mm.tags, mm.cfg.extractors, mm.cfg.onError)
}

// Gauge creates a gauge with possible predefined tags
func (mm *mortarMetric) Gauge(name, desc string) monitor.TagsAwareGauge {
	bricksGauge, err := mm.registry.loadOrStoreGauge(name, desc, mm.extractTagKeys()...)
	if err != nil {
		mm.cfg.onError(fmt.Errorf("error registering gauge [%s:%s] metric with %v tags, %w", name, desc, mm.extractTagKeys(), err))
		bricksGauge = newNoopGauge(name, desc, err, mm.cfg.onError)
	}

	return newGaugeWithTags(bricksGauge, mm.tags, mm.cfg.extractors, mm.cfg.onError)
}

// Histogram creates a histogram with possible predefined tags
func (mm *mortarMetric) Histogram(name, desc string, buckets monitor.Buckets) monitor.TagsAwareHistogram {
	bricksHistogram, err := mm.registry.loadOrStoreHistogram(name, desc, buckets, mm.extractTagKeys()...)
	if err != nil {
		mm.cfg.onError(fmt.Errorf("error registering histogram [%s:%s] metric with %v tags, %w", name, desc, mm.extractTagKeys(), err))
		bricksHistogram = newNoopHistogram(name, desc, err, mm.cfg.onError)
	}

	return newHistogramWithTags(bricksHistogram, mm.tags, mm.cfg.extractors, mm.cfg.onError)
}

// Timer creates a timer with possible predefined tags
func (mm *mortarMetric) Timer(name, desc string) monitor.TagsAwareTimer {
	bricksTimer, err := mm.registry.loadOrStoreTimer(name, desc, mm.extractTagKeys()...)
	if err != nil {
		mm.cfg.onError(fmt.Errorf("error registering timer [%s:%s] metric with %v tags, %w", name, desc, mm.extractTagKeys(), err))
		bricksTimer = newNoopTimer(name, desc, err, mm.cfg.onError)
	}

	return newTimerWithTags(bricksTimer, mm.tags, mm.cfg.extractors, mm.cfg.onError)
}

// WithTags sets custom tags to be included if possible in every Metric
func (mm *mortarMetric) WithTags(tags monitor.Tags) monitor.Metrics {
	mm.withTags(tags)
	return mm
}

func (mm *mortarMetric) extractTagKeys() (keys []string) {
	for k := range mm.tags {
		keys = append(keys, k)
	}
	if len(keys) > 0 {
		sort.Strings(keys)
	}
	return
}
