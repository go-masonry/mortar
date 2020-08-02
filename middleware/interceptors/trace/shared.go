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

func (d tracingDeps) extractIncomingCarrier(ctx context.Context) mdRW {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	return mdRW(md.Copy()) // make a copy since this map is not thread safe
}

func (d tracingDeps) extractOutgoingCarrier(ctx context.Context) mdRW {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	return mdRW(md.Copy()) // make a copy since this map is not thread safe
}

type mdRW metadata.MD

func (md mdRW) Set(key, value string) {
	metadata.MD(md).Set(key, value)
}

func (md mdRW) ForeachKey(handler func(key, value string) error) error {
	for k, vv := range md {
		for _, v := range vv {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

var grpcTag = opentracing.Tag{Key: string(ext.Component), Value: "gRPC"}
var restTag = opentracing.Tag{Key: string(ext.Component), Value: "REST"}
