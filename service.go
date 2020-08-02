package mortar

import (
	"context"
	"fmt"
	"github.com/go-masonry/mortar/interfaces/http/server"
	"github.com/go-masonry/mortar/interfaces/log"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
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
			healthClient := grpc_health_v1.NewHealthClient(conn)
			_, err = healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: ""})
		}
	}
	if err == nil {
		for port, srvType := range ports {
			deps.Logger.Debug(ctx, "Service is accepting %s calls on %d", srvType, port)
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
