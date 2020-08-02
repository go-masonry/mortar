package constructors

import (
	"context"
	"fmt"
	"github.com/go-masonry/mortar/health"
	"github.com/go-masonry/mortar/interfaces/http/server"
	"github.com/go-masonry/mortar/interfaces/log"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type webServiceDependencies struct {
	fx.In

	LifeCycle fx.Lifecycle

	Logger            log.Logger
	WebServiceBuilder server.WebService
}

// Service should be invoked by FX, it will build the entire dependencies graph and add lifecycle hooks
func Service(deps webServiceDependencies) {
	deps.LifeCycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go deps.WebServiceBuilder.Run(ctx) // this should exit only when service was shutdown
			return deps.pingService(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return deps.WebServiceBuilder.Stop(ctx)
		},
	})
}

func (deps webServiceDependencies) pingService(ctx context.Context) (err error) {
	err = fmt.Errorf("failed to check internal service health")
	ports := deps.WebServiceBuilder.Ports()
	if grpcAddress := deps.getGRPCAddress(); len(grpcAddress) > 0 {
		var conn *grpc.ClientConn
		if conn, err = grpc.DialContext(ctx, grpcAddress, grpc.WithInsecure()); err == nil {
			defer conn.Close()
			healthClient := health.NewHealthClient(conn)
			_, err = healthClient.Check(ctx, &health.HealthCheckRequest{})
		}
	}
	if err == nil {
		for _, info := range ports {
			deps.Logger.Debug(ctx, "Service is accepting %s calls on %s", info.Type, info.Address)
		}
	}
	return
}
func (deps webServiceDependencies) getGRPCAddress() string {
	for _, info := range deps.WebServiceBuilder.Ports() {
		if info.Type == server.GRPCServer {
			return info.Address
		}
	}
	return ""
}
