package monitoring

import (
	"sort"
	"strings"
	"sync"

	"github.com/go-masonry/mortar/interfaces/monitor"
)

type externalRegistry struct {
	external monitor.BricksMetrics
	// TODO perhaps change this to self evicting cache that will remove metrics if unused for a long time to save space
	counters   *sync.Map
	gauges     *sync.Map
	histograms *sync.Map
	timers     *sync.Map
}

func newRegistry(externalMetrics monitor.BricksMetrics) *externalRegistry {
	return &externalRegistry{
		external:   externalMetrics,
		counters:   new(sync.Map),
		gauges:     new(sync.Map),
		histograms: new(sync.Map),
		timers:     new(sync.Map),
	}
}

func (r *externalRegistry) loadOrStoreCounter(name, desc string, keys ...string) (monitor.BricksCounter, error) {
	ID := calcID(name, keys...)
	if known, ok := r.counters.Load(ID); ok {
		return known.(monitor.BricksCounter), nil
	}
	bricksCounter, err := r.external.Counter(name, desc, keys...)
	if err == nil {
		r.counters.Store(ID, bricksCounter)
	}
	return bricksCounter, err
}

func (r *externalRegistry) loadOrStoreGauge(name, desc string, keys ...string) (monitor.BricksGauge, error) {
	ID := calcID(name, keys...)
	if known, ok := r.gauges.Load(ID); ok {
		return known.(monitor.BricksGauge), nil
	}
	bricksGauge, err := r.external.Gauge(name, desc, keys...)
	if err == nil {
		r.gauges.Store(ID, bricksGauge)
	}
	return bricksGauge, err
}

func (r *externalRegistry) loadOrStoreHistogram(name, desc string, buckets monitor.Buckets, keys ...string) (monitor.BricksHistogram, error) {
	ID := calcID(name, keys...)
	if known, ok := r.histograms.Load(ID); ok {
		return known.(monitor.BricksHistogram), nil
	}
	bricksHistogram, err := r.external.Histogram(name, desc, buckets, keys...)
	if err == nil {
		r.histograms.Store(ID, bricksHistogram)
	}
	return bricksHistogram, err
}

func (r *externalRegistry) loadOrStoreTimer(name, desc string, keys ...string) (monitor.BricksTimer, error) {
	ID := calcID(name, keys...)
	if known, ok := r.timers.Load(ID); ok {
		return known.(monitor.BricksTimer), nil
	}
	bricksTimer, err := r.external.Timer(name, desc, keys...)
	if err == nil {
		r.timers.Store(ID, bricksTimer)
	}
	return bricksTimer, err
}

func calcID(name string, keys ...string) (ID string) {
	if len(keys) > 0 {
		var stringsSet = make(map[string]struct{}, len(keys))
		var filteredKeys = make([]string, 0, len(keys)) // preallocate slice
		for _, key := range keys {
			if _, ok := stringsSet[key]; !ok {
				stringsSet[key] = struct{}{}
				filteredKeys = append(filteredKeys, key)
			}
		}
		sort.Strings(filteredKeys)
		return name + "_" + strings.Join(filteredKeys, "_")
	}
	return name
}
