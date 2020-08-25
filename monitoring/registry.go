package monitoring

type cachedMetric struct{}

type registry struct {
	counters   map[string]*cachedMetric
	gauges     map[string]*cachedMetric
	histograms map[string]*cachedMetric
}

func (r *registry) getOrAddNew() {
	panic("implement me")
}
