package trace

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/go-masonry/mortar/interfaces/http/client"
	"github.com/go-masonry/mortar/mortar"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// TracerGRPCClientInterceptor is a grpc tracing client interceptor, it can log req/resp if needed
func TracerGRPCClientInterceptor(deps tracingDeps) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if deps.Tracer == nil {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		var span opentracing.Span
		span, ctx = deps.newClientSpanForGRPC(ctx, method)
		defer span.Finish()
		// log request if needed
		if deps.Config.Get(mortar.MiddlewareClientGRPCTraceIncludeRequest).Bool() {
			addBodyToSpan(span, "request", req)
		}
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			ext.LogError(span, err)
		} else {
			// log response if needed
			if deps.Config.Get(mortar.MiddlewareClientGRPCTraceIncludeResponse).Bool() {
				addBodyToSpan(span, "response", reply)
			}
		}
		return err
	}
}

// TracerRESTClientInterceptor is a REST tracing client interceptor, it can log req/resp if needed
func TracerRESTClientInterceptor(deps tracingDeps) client.HTTPClientInterceptor {
	return func(req *http.Request, handler client.HTTPHandler) (resp *http.Response, err error) {
		if deps.Tracer == nil {
			return handler(req)
		}
		span, ctx := deps.newClientSpanForREST(req)
		defer span.Finish()

		req = req.WithContext(ctx)
		resp, err = handler(req)
		if err != nil {
			ext.LogError(span, err)
		} else {
			ext.HTTPStatusCode.Set(span, uint16(resp.StatusCode))
			if deps.Config.Get(mortar.MiddlewareClientRESTTraceIncludeResponse).Bool() {
				if respDump, dumpErr := httputil.DumpResponse(resp, true); dumpErr == nil {
					addBodyToSpan(span, "response", respDump)
				} else {
					deps.Logger.WithError(dumpErr).Debug(ctx, "failed to dump response")
				}
			}
		}
		return
	}
}

func (d tracingDeps) newClientSpanForGRPC(ctx context.Context, methodName string) (opentracing.Span, context.Context) {
	span, clientContext := opentracing.StartSpanFromContextWithTracer(ctx, d.Tracer, methodName, ext.SpanKindRPCClient, grpcTag)
	carrier := d.extractOutgoingCarrier(clientContext)
	if err := d.Tracer.Inject(span.Context(), opentracing.HTTPHeaders, carrier); err != nil {
		d.Logger.WithError(err).Warn(ctx, "failed injecting trace info")
	}
	clientContext = metadata.NewOutgoingContext(clientContext, metadata.MD(carrier))
	return span, clientContext
}

func (d tracingDeps) newClientSpanForREST(req *http.Request) (opentracing.Span, context.Context) {
	var ctx = context.Background()
	if req.Context() != nil {
		ctx = req.Context()
	}
	span, clientContext := opentracing.StartSpanFromContextWithTracer(ctx, d.Tracer, req.URL.Path, ext.SpanKindRPCClient, restTag)
	if d.Config.Get(mortar.MiddlewareClientRESTTraceIncludeRequest).Bool() {
		if reqDump, dumpErr := httputil.DumpRequestOut(req, true); dumpErr == nil {
			addBodyToSpan(span, "request", reqDump)
		} else {
			d.Logger.WithError(dumpErr).Debug(ctx, "failed to dump request")
		}
	}
	ext.HTTPUrl.Set(span, fmt.Sprintf("%v", req.URL))
	ext.HTTPMethod.Set(span, req.Method)
	if err := d.Tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header)); err != nil {
		d.Logger.WithError(err).Warn(ctx, "failed injecting trace info")
	}
	return span, clientContext
}
