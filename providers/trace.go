package providers

import (
	"github.com/go-masonry/mortar/middleware/interceptors/trace"
	"github.com/go-masonry/mortar/providers/groups"
	"go.uber.org/fx"
)

// TracerGRPCClientInterceptorFxOption adds grpc trace client interceptor to the graph
func TracerGRPCClientInterceptorFxOption() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group:  groups.GRPCUnaryClientInterceptors,
			Target: trace.TracerGRPCClientInterceptor,
		})
}

// TracerRESTClientInterceptorFxOption adds REST trace client interceptor to the graph
func TracerRESTClientInterceptorFxOption() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group:  groups.RESTClientInterceptors,
			Target: trace.TracerRESTClientInterceptor,
		})
}

// GRPCTracingUnaryServerInterceptorFxOption adds grpc trace unary server interceptor to the graph
func GRPCTracingUnaryServerInterceptorFxOption() fx.Option {
	return fx.Provide(fx.Annotated{
		Group:  groups.UnaryServerInterceptors,
		Target: trace.GRPCTracingUnaryServerInterceptor,
	})
}
