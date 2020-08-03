package providers

import (
	"github.com/go-masonry/mortar/constructors"
	"go.uber.org/fx"
)

// MonitorFxOption adds default metric client to the graph
func MonitorFxOption() fx.Option {
	return fx.Provide(constructors.DefaultMonitor)
}
