package server

import (
	"context"
	"time"

	"github.com/go-masonry/mortar/interfaces/cfg"
	confkeys "github.com/go-masonry/mortar/interfaces/cfg/keys"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/go-masonry/mortar/utils"
	"go.uber.org/fx"
	"google.golang.org/grpc"
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
		if logLevel := deps.Config.Get(confkeys.MiddlewareLogLevel); logLevel.IsSet() {
			level := log.ParseLevel(logLevel.String())
			if err != nil {
				if logErrorLevel := deps.Config.Get(confkeys.MiddlewareLogErrorLevel); logErrorLevel.IsSet() {
					level = log.ParseLevel(logErrorLevel.String())
				}
			}
			entry := deps.Logger.
				WithError(err).
				WithField("api", info.FullMethod).
				WithField("start", start).
				WithField("duration", time.Since(start).String())
			if d, ok := ctx.Deadline(); ok {
				entry = entry.WithField("deadline", d)
			}
			// log request if needed
			if deps.Config.Get(confkeys.MiddlewareLogIncludeRequest).Bool() {
				entry = addBodyToLogger(entry, "request", req)
			}
			// log response if needed
			if deps.Config.Get(confkeys.MiddlewareLogIncludeResponse).Bool() {
				entry = addBodyToLogger(entry, "response", resp)
			}
			entry.Custom(ctx, level, 0, "gRPC call finished")
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
