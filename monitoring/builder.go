package monitoring

import (
	"container/list"
	"log"

	"github.com/go-masonry/mortar/interfaces/monitor"
)

type monitorConfig struct {
	tags       monitor.Tags
	extractors []monitor.ContextExtractor
	onError    func(err error)
	reporter   monitor.BricksReporter
}

// WrapperBuilder is a helper builder to define internal Mortar monitoring wrapper
type WrapperBuilder interface {
	// Build builds monitor.Reporter
	Build(bricksBuilder monitor.Builder) monitor.Reporter
	// DoOnError is a helper function to act when receiving an error during Metric creation
	DoOnError(onError func(error)) WrapperBuilder
	// AddExtractors adds ContextExtractors that might override tag values when calling metric functions
	AddExtractors(extractors ...monitor.ContextExtractor) WrapperBuilder
	// SetTags saves defaults tags, these tags will always be included in every metric
	SetTags(tags monitor.Tags) WrapperBuilder
}

type wrapperBuilder struct {
	ll *list.List
}

// Builder creates a WrapperBuilder
func Builder() WrapperBuilder {
	return &wrapperBuilder{
		ll: list.New(),
	}
}

func (b *wrapperBuilder) SetTags(tags monitor.Tags) WrapperBuilder {
	b.ll.PushBack(func(cfg *monitorConfig) {
		if tags != nil {
			cfg.tags = tags // make sure tags are always empty, not nil
		}
	})
	return b
}

func (b *wrapperBuilder) AddExtractors(extractors ...monitor.ContextExtractor) WrapperBuilder {
	b.ll.PushBack(func(cfg *monitorConfig) {
		cfg.extractors = append(cfg.extractors, extractors...)
	})
	return b
}

func (b *wrapperBuilder) DoOnError(onError func(error)) WrapperBuilder {
	b.ll.PushBack(func(cfg *monitorConfig) {
		cfg.onError = onError
	})
	return b
}

func (b *wrapperBuilder) Build(bricksBuilder monitor.Builder) monitor.Reporter {
	cfg := new(monitorConfig)
	for e := b.ll.Front(); e != nil; e = e.Next() {
		f := e.Value.(func(*monitorConfig))
		f(cfg)
	}
	if cfg.onError == nil {
		cfg.onError = func(err error) {
			log.Printf("WARNING: monitoring error, %v", err)
		}
	}
	if cfg.tags == nil {
		cfg.tags = monitor.Tags{}
	}
	cfg.reporter = bricksBuilder.Build()
	return newMortarReporter(cfg)
}

var _ WrapperBuilder = (*wrapperBuilder)(nil)
