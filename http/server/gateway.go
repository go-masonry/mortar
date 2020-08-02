package server

import (
	"context"
	"github.com/go-masonry/mortar/utils"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/metadata"
	"net/http"
)

// TODO add grpc-gateway custom header mapper

// MetadataTraceCarrierOption is a nice trick to avoid creating an additional server span from REST to GRPC
// it will extract trace context from Headers and put that information into the Context without creating new Span
func MetadataTraceCarrierOption(tracer opentracing.Tracer) runtime.ServeMuxOption {
	return runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
		var md = metadata.New(nil)
		if tracer != nil {
			spanContext, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
			if err == nil {
				// we ignore error here, since we assume that a new span will be open by gRPC Trace Interceptor anyway
				tracer.Inject(spanContext, opentracing.HTTPHeaders, utils.MDTraceCarrier(md))
			}
		}
		return md
	})
}
