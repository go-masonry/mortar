package trace

import (
	"context"
	"github.com/go-masonry/mortar/mortar"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
)

// GRPCTracingUnaryServerInterceptor is a grpc unary server interceptor that adds trace information of the invoked grpc method and starts a new span
func GRPCTracingUnaryServerInterceptor(deps tracingDeps) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if deps.Tracer == nil {
			return handler(ctx, req)
		}
		var span opentracing.Span
		span, ctx = deps.newServerSpan(ctx, info.FullMethod)
		defer span.Finish()

		// log request if needed
		if deps.Config.Get(mortar.MiddlewareServerGRPCTraceIncludeRequest).Bool() {
			addBodyToSpan(span, "request", req)
		}
		// call handler
		resp, err = handler(ctx, req)
		if err != nil {
			ext.LogError(span, err)
		} else {
			// log response if needed
			if deps.Config.Get(mortar.MiddlewareServerGRPCTraceIncludeResponse).Bool() {
				addBodyToSpan(span, "response", resp)
			}
		}

		return resp, err
	}
}

func (d tracingDeps) newServerSpan(ctx context.Context, methodName string) (opentracing.Span, context.Context) {
	spanContext, extractError := d.Tracer.Extract(opentracing.HTTPHeaders, d.extractIncomingCarrier(ctx))
	if extractError != nil && extractError != opentracing.ErrSpanContextNotFound {
		d.Logger.WithError(extractError).Debug(ctx, "failed extracting trace info") // really low level information in my opinion
	}
	return opentracing.StartSpanFromContextWithTracer(ctx, d.Tracer, methodName, ext.RPCServerOption(spanContext), ext.SpanKindRPCServer, grpcTag)
}
