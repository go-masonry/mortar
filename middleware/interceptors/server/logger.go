package server

import (
	"context"
	"github.com/go-masonry/mortar/interfaces/cfg"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/go-masonry/mortar/mortar"
	"github.com/go-masonry/mortar/utils"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"time"
)

type loggerInterceptorDeps struct {
	fx.In

	Config cfg.Config
	Logger log.Logger
}

// LoggerGRPCInterceptor logging interceptor, it will log grpc server call with request/response if configured
func LoggerGRPCInterceptor(deps loggerInterceptorDeps) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()
		resp, err = handler(ctx, req)
		if logLevel := deps.Config.Get(mortar.ServerGRPCLogLevel); logLevel.IsSet() {
			level := log.ParseLevel(logLevel.String())
			entry := deps.Logger.
				WithError(err).
				WithField("api", info.FullMethod).
				WithField("start", start).
				WithField("duration", time.Since(start))
			if d, ok := ctx.Deadline(); ok {
				entry = entry.WithField("deadline", d)
			}
			if deps.Config.Get(mortar.MiddlewareServerGRPCLogIncludeRequest).Bool() { // if absent will return false anyway
				entry = addBodyToLogger(entry, "request", req)
			}
			if deps.Config.Get(mortar.MiddlewareServerGRPCLogIncludeResponse).Bool() {
				entry = addBodyToLogger(entry, "response", resp)
			}
			entry.Custom(ctx, level, "%s finished", info.FullMethod)
		}
		return
	}
}

func addBodyToLogger(entry log.Fields, name string, i interface{}) log.Fields {
	if bytes, err := utils.MarshalMessageBody(i); err == nil {
		return entry.WithField(name, bytes)
	}
	return entry
}
