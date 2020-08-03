package providers

import (
	"github.com/go-masonry/mortar/constructors"
	"github.com/go-masonry/mortar/constructors/partial"
	"github.com/go-masonry/mortar/handlers"
	"github.com/go-masonry/mortar/http/server"
	"github.com/go-masonry/mortar/middleware/interceptors/client"
	"go.uber.org/fx"
)

// CreateEntireWebServiceDependencyGraph creates the entire dependency graph
// and registers all provided fx.LifeCycle hooks
func CreateEntireWebServiceDependencyGraph() fx.Option {
	return fx.Invoke(constructors.Service)
}

// HttpClientBuildersFxOption adds both (GRPC, REST) partial http clients to the graph
func HttpClientBuildersFxOption() fx.Option {
	return fx.Provide(
		partial.HttpClientBuilder,
		partial.GRPCClientConnectionBuilder,
	)
}

// HttpServerFxOption adds Default Http Server builder which later injected to the Service Invoke option
// by calling CreateEntireWebServiceDependencyGraph fx.Invoke option to the graph
func HttpServerFxOption() fx.Option {
	return fx.Provide(partial.HttpServerBuilder)
}

// InternalDebugHandlersFxOption adds Internal Debug Handlers to the graph
func InternalDebugHandlersFxOption() fx.Option {
	return fx.Provide(fx.Annotated{
		Group:  partial.FxGroupInternalHttpHandlers + ",flatten",
		Target: handlers.InternalDebugHandlers,
	})
}

// InternalProfileHandlerFunctionsFxOption adds Internal Profiler Handler to the graph
func InternalProfileHandlerFunctionsFxOption() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group:  partial.FxGroupInternalHttpHandlerFunctions + ",flatten",
			Target: handlers.InternalProfileHandlerFunctions,
		})
}

// SelfHandlersFxOption adds Internal Self Build and Config information Handler to the graph
func SelfHandlersFxOption() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group:  partial.FxGroupInternalHttpHandlers + ",flatten",
			Target: handlers.SelfHandlers,
		})
}

// GRPCGatewayMetadataTraceCarrierFxOption adds GRPCGatewayMuxOption that will inject trace into the context.Context
// Make sure to understand what it does by reading server.MetadataTraceCarrierOption code and explanation
func GRPCGatewayMetadataTraceCarrierFxOption() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group:  partial.FxGroupGRPCGatewayMuxOptions,
			Target: server.MetadataTraceCarrierOption,
		})
}

// CopyGRPCHeadersClientInterceptorFxOption adds grpc Client Interceptor that copies values from grpc Incoming to Outgoing metadata
func CopyGRPCHeadersClientInterceptorFxOption() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group:  partial.FxGroupGRPCUnaryClientInterceptors,
			Target: client.CopyGRPCHeadersClientInterceptor,
		})
}
