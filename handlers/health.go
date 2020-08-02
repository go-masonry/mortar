package handlers

import (
	"context"
	"fmt"
	"github.com/go-masonry/mortar/constructors/partial"
	"github.com/go-masonry/mortar/interfaces/cfg"
	"github.com/go-masonry/mortar/interfaces/http/client"
	"github.com/go-masonry/mortar/interfaces/http/server"
	"github.com/go-masonry/mortar/mortar"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net/http"
)

func HealthHandlerOption() fx.Option {
	return fx.Provide(fx.Annotated{
		Group:  partial.FxGroupInternalHttpHandlers,
		Target: HealthHandler,
	})
}

type healthHandlerDeps struct {
	fx.In

	WebServer     server.WebService
	Config        cfg.Config
	ClientBuilder client.GRPCClientConnectionBuilder
}

func HealthHandler(deps healthHandlerDeps) partial.HttpHandlerPatternPair {
	serveError := func(w http.ResponseWriter, err error) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	wrapper, addr, err := deps.getGrpcConnWrapper()
	var handler http.HandlerFunc = func(w http.ResponseWriter, req *http.Request) {
		if err == nil {
			var ctx = context.Background()
			if timeout := deps.Config.Get(mortar.HandlersHealthTimeout); timeout.IsSet() {
				var cancelFunc func()
				ctx, cancelFunc = context.WithTimeout(ctx, timeout.Duration())
				defer cancelFunc()
			}
			conn, handlerErr := wrapper.Dial(ctx, addr, grpc.WithInsecure()) // if you serving TLS you will need to write another health handler
			if handlerErr == nil {
				_, handlerErr = grpc_health_v1.NewHealthClient(conn).Check(ctx, new(grpc_health_v1.HealthCheckRequest))
			}
			if handlerErr == nil {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, "healthy")
			} else {
				serveError(w, handlerErr)
			}
		} else {
			serveError(w, err)
		}
	}
	return partial.HttpHandlerPatternPair{Pattern: "/health", Handler: handler}
}

func (d healthHandlerDeps) getGrpcConnWrapper() (client.GRPCClientConnectionWrapper, string, error) {
	var addr string
	for _, info := range d.WebServer.Ports() {
		if info.Type == server.GRPCServer {
			addr = info.Address
			break
		}
	}
	if len(addr) == 0 {
		return nil, "", fmt.Errorf("no gRPC service to query")
	}
	return d.ClientBuilder.Build(), addr, nil
}
