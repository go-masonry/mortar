package client

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-masonry/mortar/interfaces/http/client"
	"github.com/go-masonry/mortar/interfaces/monitor"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

// Names
const (
	ClientTimerMetric            = "client_calls_duration"
	ClientTimerMetricDescription = "Monitor external HTTP client calls"
	TargetTag                    = "target"
	PathTag                      = "path"
	ErrorTag                     = "err"
	TypeTag                      = "ctype"
	TypeGRPC                     = "grpc"
	TypeREST                     = "rest"
)

type monitorDeps struct {
	fx.In

	Metrics monitor.Metrics `optional:"true"`
}

// MonitorGRPCClientCallsInterceptor create a new GRPC Unary Client interceptor that monitor all external client calls
func MonitorGRPCClientCallsInterceptor(deps monitorDeps) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		start := time.Now()
		err = invoker(ctx, method, req, reply, cc, opts...)

		if deps.Metrics != nil {
			tags := prepareTags(cc.Target(), method, TypeGRPC, fmt.Sprintf("%t", err != nil))
			deps.Metrics.
				WithTags(tags).
				Timer(ClientTimerMetric, ClientTimerMetricDescription).
				WithContext(ctx).
				Record(time.Since(start))
		}
		return
	}
}

// MonitorRESTClientCallsInterceptor create a new REST Client interceptor that monitor all external client calls
func MonitorRESTClientCallsInterceptor(deps monitorDeps) client.HTTPClientInterceptor {
	return func(req *http.Request, handler client.HTTPHandler) (resp *http.Response, err error) {
		start := time.Now()
		resp, err = handler(req)

		if deps.Metrics != nil {
			tags := prepareTags(req.Host, req.URL.Path, TypeREST, fmt.Sprintf("%t", err != nil || resp.StatusCode >= http.StatusBadRequest))
			deps.Metrics.
				WithTags(tags).
				Timer(ClientTimerMetric, ClientTimerMetricDescription).
				WithContext(req.Context()).
				Record(time.Since(start))
		}
		return
	}
}

func prepareTags(host, path, clientType, err string) monitor.Tags {
	host = strings.Trim(host, ":") // remove trailing port if exists
	return monitor.Tags{
		TargetTag: host,
		PathTag:   path,
		ErrorTag:  err,
		TypeTag:   clientType,
	}
}
