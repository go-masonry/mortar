package providers

import (
	"github.com/go-masonry/mortar/constructors"
	"github.com/go-masonry/mortar/constructors/partial"
	"github.com/go-masonry/mortar/middleware/interceptors/server"
	"go.uber.org/fx"
)

// MonitorFxOption adds default metric client to the graph
func MonitorFxOption() fx.Option {
	return fx.Provide(constructors.DefaultMonitor)
}

// MonitorGRPCInterceptorFxOption adds grpc metric interceptor to the graph
func MonitorGRPCInterceptorFxOption() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group:  partial.FxGroupUnaryServerInterceptors,
			Target: server.MonitorGRPCInterceptor,
		})
}
