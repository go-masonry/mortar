package trace

import (
	"context"
	"github.com/go-masonry/mortar/interfaces/cfg"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/go-masonry/mortar/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	traceLog "github.com/opentracing/opentracing-go/log"
	"go.uber.org/fx"
	"google.golang.org/grpc/metadata"
)

type tracingDeps struct {
	fx.In

	Logger log.Logger
	Config cfg.Config
	Tracer opentracing.Tracer `optional:"true"`
}

func addBodyToSpan(span opentracing.Span, name string, msg interface{}) {
	bytes, err := utils.MarshalMessageBody(msg)
	if err == nil {
		span.LogFields(traceLog.String(name, string(bytes))) // TODO: can exceed length limit, introduce option
	} else {
		// If marshaling failed let's try to log msg.ToString()
		span.LogKV(name, msg)
	}
}

func (d tracingDeps) extractIncomingCarrier(ctx context.Context) utils.MDTraceCarrier {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	return utils.MDTraceCarrier(md.Copy()) // make a copy since this map is not thread safe
}

func (d tracingDeps) extractOutgoingCarrier(ctx context.Context) utils.MDTraceCarrier {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	return utils.MDTraceCarrier(md.Copy()) // make a copy since this map is not thread safe
}

var grpcTag = opentracing.Tag{Key: string(ext.Component), Value: "gRPC"}
var restTag = opentracing.Tag{Key: string(ext.Component), Value: "REST"}
