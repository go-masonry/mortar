package providers

import (
	"github.com/go-masonry/mortar/constructors/partial"
	"github.com/go-masonry/mortar/middleware/interceptors/server"
	"go.uber.org/fx"
)

// MonitorGRPCInterceptorFxOption adds grpc metric interceptor to the graph
func MonitorGRPCInterceptorFxOption() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group:  partial.FxGroupUnaryServerInterceptors,
			Target: server.MonitorGRPCInterceptor,
		})
}
