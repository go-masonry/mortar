package monitoring

import (
	"sort"
	"strings"
	"sync"

	"github.com/go-masonry/mortar/interfaces/monitor"
)

type externalRegistry struct {
	cm       sync.RWMutex
	gm       sync.RWMutex
	tm       sync.RWMutex
	hm       sync.RWMutex
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

func (r *externalRegistry) loadOrStoreCounter(name, desc string, keys ...string) (bricksCounter monitor.BricksCounter, err error) {
	ID := calcID(name, keys...)
	if known, ok := r.counters.Load(ID); ok {
		return known.(monitor.BricksCounter), nil
	}
	r.cm.Lock()
	defer r.cm.Unlock()
	bricksCounter, err = r.external.Counter(name, desc, keys...)
	if err == nil {
		cacheValue, _ := r.counters.LoadOrStore(ID, bricksCounter) // if a previous duplicate is already there (was created by other go routine)
		bricksCounter = cacheValue.(monitor.BricksCounter)
	} else {
		if known, ok := r.counters.Load(ID); ok { // perhaps it's already there (was created by other go routine) and the underlying impl have a dup check
			bricksCounter, err = known.(monitor.BricksCounter), nil
			return
		}
	}
	return
}

func (r *externalRegistry) loadOrStoreGauge(name, desc string, keys ...string) (bricksGauge monitor.BricksGauge, err error) {
	ID := calcID(name, keys...)
	if known, ok := r.gauges.Load(ID); ok {
		return known.(monitor.BricksGauge), nil
	}
	r.gm.Lock()
	defer r.gm.Unlock()
	bricksGauge, err = r.external.Gauge(name, desc, keys...)
	if err == nil {
		cacheValue, _ := r.gauges.LoadOrStore(ID, bricksGauge) // if a previous duplicate is already there (was created by other go routine)
		bricksGauge = cacheValue.(monitor.BricksGauge)
	} else {
		if known, ok := r.gauges.Load(ID); ok { // perhaps it's already there (was created by other go routine) and the underlying impl have a dup check
			bricksGauge, err = known.(monitor.BricksGauge), nil
			return
		}
	}
	return
}

func (r *externalRegistry) loadOrStoreHistogram(name, desc string, buckets monitor.Buckets, keys ...string) (bricksHistogram monitor.BricksHistogram, err error) {
	ID := calcID(name, keys...)
	if known, ok := r.histograms.Load(ID); ok {
		return known.(monitor.BricksHistogram), nil
	}
	r.hm.Lock()
	defer r.hm.Unlock()
	bricksHistogram, err = r.external.Histogram(name, desc, buckets, keys...)
	if err == nil {
		cacheValue, _ := r.histograms.LoadOrStore(ID, bricksHistogram) // if a previous duplicate is already there (was created by other go routine)
		bricksHistogram = cacheValue.(monitor.BricksHistogram)
	} else {
		if known, ok := r.histograms.Load(ID); ok { // perhaps it's already there (was created by other go routine) and the underlying impl have a dup check
			bricksHistogram, err = known.(monitor.BricksHistogram), nil
			return
		}
	}
	return
}

func (r *externalRegistry) loadOrStoreTimer(name, desc string, keys ...string) (bricksTimer monitor.BricksTimer, err error) {
	ID := calcID(name, keys...)
	if known, ok := r.timers.Load(ID); ok {
		return known.(monitor.BricksTimer), nil
	}
	r.tm.Lock()
	defer r.tm.Unlock()
	bricksTimer, err = r.external.Timer(name, desc, keys...)
	if err == nil {
		cacheValue, _ := r.timers.LoadOrStore(ID, bricksTimer) // if a previous duplicate is already there (was created by other go routine)
		bricksTimer = cacheValue.(monitor.BricksTimer)
	} else {
		if known, ok := r.timers.Load(ID); ok { // perhaps it's already there (was created by other go routine) and the underlying impl have a dup check
			bricksTimer, err = known.(monitor.BricksTimer), nil
			return
		}
	}
	return
}

func calcID(name string, keys ...string) (ID string) {
	if len(keys) > 0 {
		var stringsSet = make(map[string]struct{}, len(keys))
		var parts = make([]string, 0, len(keys)+1) // preallocate slice with extra space
		for _, key := range keys {
			if _, ok := stringsSet[key]; !ok {
				stringsSet[key] = struct{}{}
				parts = append(parts, key)
			}
		}
		sort.Strings(parts)
		// avoid allocation and prepend name
		parts = append(parts, "")  // add empty string to the end -> len++
		copy(parts[1:], parts[0:]) // shift
		parts[0] = name            // name is the first string now
		return strings.Join(parts, "_")
	}
	return name
}
