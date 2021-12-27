package groups

import (
	"github.com/go-masonry/mortar/constructors"
	"github.com/go-masonry/mortar/constructors/partial"
)

const (
	// RESTClientInterceptors - REST Client Interceptors group. This group will help you configure Mortar default REST Client builder
	RESTClientInterceptors = partial.FxGroupRESTClientInterceptors

	// GRPCUnaryClientInterceptors - GRPC Unary Client Interceptors group. This group will help you configure Mortar default GRPC Client builder
	GRPCUnaryClientInterceptors = partial.FxGroupGRPCUnaryClientInterceptors

	// GRPCServerAPIs - Mortar GRPC Service APIs group. This group is responsible on registering your gRPC server implementation
	GRPCServerAPIs = partial.FxGroupGRPCServerAPIs

	// GRPCGatewayGeneratedHandlers - GRPC Gateway generated handlers group. This group is responsible on registering your reverse-proxy over gRPC
	GRPCGatewayGeneratedHandlers = partial.FxGroupGRPCGatewayGeneratedHandlers

	// GRPCGatewayMuxOptions - GRPC Gateway Mux Options group. Use this group if you want to provide different GRPC-GW Mux options
	// https://grpc-ecosystem.github.io/grpc-gateway/docs/customizingyourgateway.html
	GRPCGatewayMuxOptions = partial.FxGroupGRPCGatewayMuxOptions

	// ExternalHTTPHandlers - External Http Handlers group, add your custom external HTTP Handlers
	ExternalHTTPHandlers = partial.FxGroupExternalHTTPHandlers

	// ExternalHTTPHandlerFunctions - External Http Handlers function group, add your custom external HTTP Handler Functions
	ExternalHTTPHandlerFunctions = partial.FxGroupExternalHTTPHandlerFunctions

	// ExternalHTTPInterceptors - External Http Interceptors group, add your custom external HTTP interceptors
	ExternalHTTPInterceptors = partial.FxGroupExternalHTTPInterceptors

	// UnaryServerInterceptors - GRPC Unary Server Interceptors group. Register different gRPC server interceptors
	UnaryServerInterceptors = partial.FxGroupUnaryServerInterceptors

	// InternalHTTPHandlers - Internal Http Handlers group. Mortar comes with several internal handlers, you can add yours.
	InternalHTTPHandlers = partial.FxGroupInternalHTTPHandlers

	//InternalHTTPHandlerFunctions - Internal Http Handler Functions group. Similar toInternalHttpHandlers but for functions
	InternalHTTPHandlerFunctions = partial.FxGroupInternalHTTPHandlerFunctions

	// InternalHTTPInterceptors - Internal Http Interceptors group, add your custom external HTTP interceptors
	InternalHTTPInterceptors = partial.FxGroupInternalHTTPInterceptors

	// LoggerContextExtractors -  Default Logger Context extractors group. Add custom extractors to enrich your logs from context.Context
	LoggerContextExtractors = constructors.FxGroupLoggerContextExtractors

	// MonitorContextExtractors - Monitor Context extractors group. Add different tags from context to each metric
	MonitorContextExtractors = constructors.FxGroupMonitorContextExtractors
)
