package providers

import (
	"github.com/go-masonry/mortar/constructors"
	"github.com/go-masonry/mortar/constructors/partial"
	"github.com/go-masonry/mortar/handlers"
	"github.com/go-masonry/mortar/http/server"
	"github.com/go-masonry/mortar/middleware/interceptors/client"
	"github.com/go-masonry/mortar/providers/groups"
	"go.uber.org/fx"
)

// BuildMortarWebService creates the entire dependency graph
// and registers all provided fx.LifeCycle hooks
func BuildMortarWebService() fx.Option {
	return fx.Invoke(constructors.Service)
}

// HttpClientBuildersFxOption adds both (GRPC, REST) partial http clients to the graph
func HttpClientBuildersFxOption() fx.Option {
	return fx.Provide(
		partial.HttpClientBuilder,
		partial.GRPCClientConnectionBuilder,
	)
}

// HttpServerBuilderFxOption adds Default Http Server builder which later injected to the Service Invoke option
// by calling BuildMortarWebService fx.Invoke option to the graph
func HttpServerBuilderFxOption() fx.Option {
	return fx.Provide(partial.HttpServerBuilder)
}

// InternalDebugHandlersFxOption adds Internal Debug Handlers to the graph
func InternalDebugHandlersFxOption() fx.Option {
	return fx.Provide(fx.Annotated{
		Group:  groups.InternalHttpHandlers + ",flatten",
		Target: handlers.InternalDebugHandlers,
	})
}

// InternalProfileHandlerFunctionsFxOption adds Internal Profiler Handler to the graph
func InternalProfileHandlerFunctionsFxOption() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group:  groups.InternalHttpHandlerFunctions + ",flatten",
			Target: handlers.InternalProfileHandlerFunctions,
		})
}

// InternalSelfHandlersFxOption adds Internal Self Build and Config information Handler to the graph
func InternalSelfHandlersFxOption() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group:  groups.InternalHttpHandlers + ",flatten",
			Target: handlers.SelfHandlers,
		})
}

// GRPCGatewayMetadataTraceCarrierFxOption adds GRPCGatewayMuxOption that will inject trace into the context.Context
// Make sure to understand what it does by reading server.MetadataTraceCarrierOption code and explanation
func GRPCGatewayMetadataTraceCarrierFxOption() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group:  groups.GRPCGatewayMuxOptions,
			Target: server.MetadataTraceCarrierOption,
		})
}

// CopyGRPCHeadersClientInterceptorFxOption adds grpc Client Interceptor that copies values from grpc Incoming to Outgoing metadata
func CopyGRPCHeadersClientInterceptorFxOption() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group:  groups.GRPCUnaryClientInterceptors,
			Target: client.CopyGRPCHeadersClientInterceptor,
		})
}
