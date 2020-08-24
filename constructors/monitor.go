package constructors

import (
	"context"

	"github.com/go-masonry/mortar/interfaces/cfg"
	"github.com/go-masonry/mortar/interfaces/monitor"
	"github.com/go-masonry/mortar/monitoring"
	"github.com/go-masonry/mortar/mortar"
	"go.uber.org/fx"
)

const (
	// FxGroupMonitorContextExtractors defines group name
	FxGroupMonitorContextExtractors = "monitorContextExtractors"
)

type monitorDeps struct {
	fx.In

	LifeCycle         fx.Lifecycle
	Config            cfg.Config
	MonitorBuilder    monitor.Builder
	ContextExtractors []monitor.ContextExtractor `group:"monitorContextExtractors"`
}

// DefaultMonitor is a constructor that will create a Metrics reporter based on values from the Config Map
// such as
//
// 	- Tags: we will look for default tags using mortar.MonitorTagsKey within the configuration map
//
func DefaultMonitor(deps monitorDeps) monitor.Metrics {
	tags := deps.Config.Get(mortar.MonitorTagsKey).StringMapString() // can be empty

	reporter := monitoring.NewMortarReporter(deps.MonitorBuilder, tags, deps.ContextExtractors...)
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
