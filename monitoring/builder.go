package monitoring

import (
	"container/list"

	"github.com/go-masonry/mortar/interfaces/monitor"
)

type monitorConfig struct {
	tags       monitor.Tags
	extractors []monitor.ContextExtractor
	onError    func(err error)
}

// WrapperBuilder is a helper builder to define internal Mortar monitoring wrapper
type WrapperBuilder interface {
	Build(bricksBuilder monitor.Builder) monitor.Reporter
	DoOnError(onError func(error)) WrapperBuilder
	AddExtractors(extractors ...monitor.ContextExtractor) WrapperBuilder
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
		cfg.tags = tags
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

	return newMortarReporter(bricksBuilder, cfg)
}

var _ WrapperBuilder = (*wrapperBuilder)(nil)
