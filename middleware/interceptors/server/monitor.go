package server

import (
	"context"
	"time"

	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/go-masonry/mortar/interfaces/monitor"
	"github.com/go-masonry/mortar/utils"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const (
	statusTagName = "status"
)

type gRPCMetricInterceptorsDeps struct {
	fx.In

	Logger  log.Logger
	Metrics monitor.Metrics `optional:"true"`
}

// MonitorGRPCInterceptor sends gRPC method invocation metrics to the configured Metrics server (Prometheus, Datadog)
func MonitorGRPCInterceptor(deps gRPCMetricInterceptorsDeps) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()
		resp, err = handler(ctx, req)
		if deps.Metrics != nil {
			var metric = deps.Metrics
			if err != nil {
				metric = metric.AddTag(statusTagName, statusTag(err))
			}
			_, methodName := utils.SplitMethodAndPackage(info.FullMethod)
			monitoringError := metric.Timing(ctx, methodName, time.Since(start)) // nothing to do with the error here
			if monitoringError != nil {
				deps.Logger.WithError(monitoringError).WithField("method", methodName).Info(ctx, "failed to send grpc timing metric")
			}
		}
		return
	}
}

func statusTag(err error) string {
	s, _ := status.FromError(err)
	return s.Code().String()
}
