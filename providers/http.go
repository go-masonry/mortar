package providers

import (
	"github.com/go-masonry/mortar/constructors"
	"github.com/go-masonry/mortar/constructors/partial"
	"github.com/go-masonry/mortar/middleware/interceptors/client"
	"github.com/go-masonry/mortar/providers/groups"
	"go.uber.org/fx"
)

// BuildMortarWebServiceFxOption creates the entire dependency graph
// and registers all provided fx.LifeCycle hooks
func BuildMortarWebServiceFxOption() fx.Option {
	return fx.Invoke(constructors.Service)
}

// BuildMortarWebService is a constructor that creates and registers fx.LifeCycle hooks for Mortar web services
//
// Consider using BuildMortarWebServiceFxOption if you only want to invoke it.
var BuildMortarWebService = constructors.Service

// HTTPServerBuilderFxOption adds Default Http Server builder which later injected to the Service Invoke option
// by calling BuildMortarWebServiceFxOption fx.Invoke option to the graph
func HTTPServerBuilderFxOption() fx.Option {
	return fx.Provide(partial.HTTPServerBuilder)
}

// HTTPServerBuilder is a constructor that creates a partial Mortars HTTP Server Builder
//
// Consider using HTTPServerBuilderFxOption if you only want to provide it.
var HTTPServerBuilder = partial.HTTPServerBuilder

// HTTPClientBuildersFxOption adds both (GRPC, REST) partial http clients to the graph
func HTTPClientBuildersFxOption() fx.Option {
	return fx.Provide(
		partial.HTTPClientBuilder,
		partial.GRPCClientConnectionBuilder,
	)
}

// HTTPClientBuilder is a constructor that creates HTTP Client builder
//
// Consider using HTTPClientBuildersFxOption if you only want to provide it.
var HTTPClientBuilder = partial.HTTPClientBuilder

// GRPCClientConnectionBuilder is a constructor that creates gRPC Connection for client builder
//
// Consider using HTTPClientBuildersFxOption if you only want to provide it.
var GRPCClientConnectionBuilder = partial.GRPCClientConnectionBuilder

// CopyGRPCHeadersClientInterceptorFxOption adds grpc Client Interceptor that copies values from grpc Incoming to Outgoing metadata
func CopyGRPCHeadersClientInterceptorFxOption() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group:  groups.GRPCUnaryClientInterceptors,
			Target: client.CopyGRPCHeadersClientInterceptor,
		})
}

// CopyGRPCHeadersClientInterceptor is a constructor that creates gRPC Unary Client Interceptor
//	*Note* normally this dependency is part of a group. If you want to create it as a standalone
// 	dependency, remember that there can be only one of this kind in the graph.
//
// Consider using CopyGRPCHeadersClientInterceptorFxOption if you only want to provide it.
var CopyGRPCHeadersClientInterceptor = client.CopyGRPCHeadersClientInterceptor
