package constructors

import (
	"context"

	"github.com/go-masonry/mortar/interfaces/cfg"
	confkeys "github.com/go-masonry/mortar/interfaces/cfg/keys"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/go-masonry/mortar/interfaces/monitor"
	"github.com/go-masonry/mortar/monitoring"
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
	Logger            log.Logger
	MonitorBuilder    monitor.Builder
	ContextExtractors []monitor.ContextExtractor `group:"monitorContextExtractors"`
}

// DefaultMonitor is a constructor that will create a Metrics reporter based on values from the Config Map
// such as
//
//   - Tags: we will look for default tags using mortar.MonitorTagsKey within the configuration map
func DefaultMonitor(deps monitorDeps) monitor.Metrics {
	tags := deps.Config.Get(confkeys.MonitorTags).StringMapString() // can be empty
	reporter := monitoring.Builder().SetTags(tags).AddExtractors(deps.ContextExtractors...).DoOnError(func(err error) {
		deps.Logger.WithError(err).Custom(nil, log.WarnLevel, 2, "monitoring error")
	}).Build(deps.MonitorBuilder)

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
