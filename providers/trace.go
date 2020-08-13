package providers

import (
	"github.com/go-masonry/mortar/http/server"
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

// TracerGRPCClientInterceptor is a constructor that creates gRPC Unary Client Interceptor
// This interceptor will report a client span to the trace server
//
//	*Note* normally this dependency is part of a group. If you want to create it as a standalone
// 	dependency, remember that there can be only one of this kind in the graph.
//
// Consider using TracerGRPCClientInterceptorFxOption if you only want to provide it.
var TracerGRPCClientInterceptor = trace.TracerGRPCClientInterceptor

// TracerRESTClientInterceptorFxOption adds REST trace client interceptor to the graph
func TracerRESTClientInterceptorFxOption() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group:  groups.RESTClientInterceptors,
			Target: trace.TracerRESTClientInterceptor,
		})
}

// TracerRESTClientInterceptor is a constructor that creates REST HTTP Client Interceptor
// This interceptor will report a client span to the trace server
//
//	*Note* normally this dependency is part of a group. If you want to create it as a standalone
// 	dependency, remember that there can be only one of this kind in the graph.
//
// Consider using TracerRESTClientInterceptorFxOption if you only want to provide it.
var TracerRESTClientInterceptor = trace.TracerRESTClientInterceptor

// GRPCTracingUnaryServerInterceptorFxOption adds grpc trace unary server interceptor to the graph
func GRPCTracingUnaryServerInterceptorFxOption() fx.Option {
	return fx.Provide(fx.Annotated{
		Group:  groups.UnaryServerInterceptors,
		Target: trace.GRPCTracingUnaryServerInterceptor,
	})
}

// GRPCTracingUnaryServerInterceptor is a constructor that creates gRPC Unary Server Interceptor
// This interceptor will report a server span to the trace server
//
//	*Note* normally this dependency is part of a group. If you want to create it as a standalone
// 	dependency, remember that there can be only one of this kind in the graph.
//
// Consider using GRPCTracingUnaryServerInterceptorFxOption if you only want to provide it.
var GRPCTracingUnaryServerInterceptor = trace.GRPCTracingUnaryServerInterceptor

// GRPCGatewayMetadataTraceCarrierFxOption adds GRPCGatewayMuxOption that will inject trace into the context.Context
// Make sure to understand what it does by reading server.MetadataTraceCarrierOption code and explanation
func GRPCGatewayMetadataTraceCarrierFxOption() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group:  groups.GRPCGatewayMuxOptions,
			Target: server.MetadataTraceCarrierOption,
		})
}

// MetadataTraceCarrierOption creates a special metadata.MD carrier for the tracer.
// Make sure to understand what it does by reading server.MetadataTraceCarrierOption code and explanation
//
//	*Note* normally this dependency is part of a group. If you want to create it as a standalone
// 	dependency, remember that there can be only one of this kind in the graph.
//
// Consider using GRPCGatewayMetadataTraceCarrierFxOption if you only want to provide it.
var MetadataTraceCarrierOption = server.MetadataTraceCarrierOption
