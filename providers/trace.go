package providers

import (
	"github.com/go-masonry/mortar/constructors/partial"
	"github.com/go-masonry/mortar/middleware/interceptors/trace"
	"go.uber.org/fx"
)

// TracerGRPCClientInterceptorFxOption adds grpc trace client interceptor to the graph
func TracerGRPCClientInterceptorFxOption() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group:  partial.FxGroupGRPCUnaryClientInterceptors,
			Target: trace.TracerGRPCClientInterceptor,
		})
}

// TracerRESTClientInterceptorFxOption adds REST trace client interceptor to the graph
func TracerRESTClientInterceptorFxOption() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group:  partial.FxGroupRESTClientInterceptors,
			Target: trace.TracerRESTClientInterceptor,
		})
}

// GRPCTracingUnaryServerInterceptorFxOption adds grpc trace unary server interceptor to the graph
func GRPCTracingUnaryServerInterceptorFxOption() fx.Option {
	return fx.Provide(fx.Annotated{
		Group:  partial.FxGroupUnaryServerInterceptors,
		Target: trace.GRPCTracingUnaryServerInterceptor,
	})
}
