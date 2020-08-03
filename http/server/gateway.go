package server

import (
	"context"
	"github.com/go-masonry/mortar/utils"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/fx"
	"google.golang.org/grpc/metadata"
	"net/http"
)

// TODO add grpc-gateway custom header mapper
type grpcGatewayMuxOptionsDeps struct {
	fx.In

	Tracer opentracing.Tracer `optional:"true"`
}

// MetadataTraceCarrierOption is a nice trick to avoid creating an additional server span from REST to GRPC
// it will extract trace context from Headers and put that information into the Context without creating a new Span
// However if you would like to create a new Span on the REST layer, you should read how to do it here
// https://grpc-ecosystem.github.io/grpc-gateway/docs/customizingyourgateway.html scroll to "OpenTracing Support"
func MetadataTraceCarrierOption(deps grpcGatewayMuxOptionsDeps) runtime.ServeMuxOption {
	return runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
		var md = metadata.New(nil)
		if deps.Tracer != nil {
			spanContext, err := deps.Tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
			if err == nil {
				// we ignore error here, since we assume that a new span will be open by gRPC Trace Interceptor anyway
				deps.Tracer.Inject(spanContext, opentracing.HTTPHeaders, utils.MDTraceCarrier(md))
			}
		}
		return md
	})
}
