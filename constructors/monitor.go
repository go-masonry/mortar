package constructors

import (
	"context"
	"github.com/go-masonry/mortar/interfaces/cfg"
	"github.com/go-masonry/mortar/interfaces/monitor"
	"github.com/go-masonry/mortar/mortar"
	"go.uber.org/fx"
)

const (
	FxGroupMonitorContextExtractors = "monitorContextExtractors"
)

type monitorDeps struct {
	fx.In

	LifeCycle         fx.Lifecycle
	Config            cfg.Config
	MonitorBuilder    monitor.Builder
	ContextExtractors []monitor.ContextExtractor `group:"monitorContextExtractors"`
}

// DefaultMonitor is a constructor that will create a Metrics client based on values from the Config Map
// such as
//
// 	- Address: we will look for a key mortar.MonitorAddressKey within the configuration map
// 	- Prefix: we will look for a key mortar.MonitorPrefixKey within the configuration map
// 	- Tags: we will look for default tags using mortar.MonitorTagsKey within the configuration map
//
func DefaultMonitor(deps monitorDeps) monitor.Metrics {
	address := deps.Config.Get(mortar.MonitorAddressKey).String()
	builder := deps.MonitorBuilder.SetAddress(address)
	if tags := deps.Config.Get(mortar.MonitorTagsKey); tags.IsSet() {
		builder = builder.SetTags(tags.StringMapString())
	}
	if prefix := deps.Config.Get(mortar.MonitorPrefixKey); prefix.IsSet() {
		builder = builder.SetPrefix(prefix.String())
	}
	reporter := builder.AddContextExtractors(deps.ContextExtractors...).Build()

	deps.LifeCycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return reporter.Connect(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return reporter.Close(ctx)
		},
	})
	return reporter.Metrics()
}
