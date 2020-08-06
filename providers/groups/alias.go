package groups

import (
	"github.com/go-masonry/mortar/constructors"
	"github.com/go-masonry/mortar/constructors/partial"
)

const (
	// REST Client Interceptors group. This group will help you configure Mortar default REST Client builder
	RESTClientInterceptors = partial.FxGroupRESTClientInterceptors

	// GRPC Unary Client Interceptors group. This group will help you configure Mortar default GRPC Client builder
	GRPCUnaryClientInterceptors = partial.FxGroupGRPCUnaryClientInterceptors

	// Mortar GRPC Service APIs group. This group is responsible on registering your gRPC server implementation
	GRPCServerAPIs = partial.FxGroupGRPCServerAPIs

	// GRPC Gateway generated handlers group. This group is responsible on registering your reverse-proxy over gRPC
	GRPCGatewayGeneratedHandlers = partial.FxGroupGRPCGatewayGeneratedHandlers

	// GRPC Gateway Mux Options group. Use this group if you want to provide different GRPC-GW Mux options
	// https://grpc-ecosystem.github.io/grpc-gateway/docs/customizingyourgateway.html
	GRPCGatewayMuxOptions = partial.FxGroupGRPCGatewayMuxOptions

	// GRPC Unary Server Interceptors group. Register different gRPC server interceptors
	UnaryServerInterceptors = partial.FxGroupUnaryServerInterceptors

	// Internal Http Handlers group. Mortar comes with several internal handlers, you can add yours.
	InternalHttpHandlers = partial.FxGroupInternalHttpHandlers

	// Internal Http Handler Functions group. Similar toInternalHttpHandlers but for functions
	InternalHttpHandlerFunctions = partial.FxGroupInternalHttpHandlerFunctions

	// Default Logger Context extractors group. Add custom extractors to enrich your logs from context.Context
	LoggerContextExtractors = constructors.FxGroupLoggerContextExtractors

	// Monitor Context extractors group. Add different tags from context to each metric
	MonitorContextExtractors = constructors.FxGroupMonitorContextExtractors
)
