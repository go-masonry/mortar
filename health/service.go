package health

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

type healthService struct {
	UnimplementedHealthServer
}

func RegisterInternalGRPCGatewayHandler(mux *runtime.ServeMux, endpoint string) error{
	return RegisterHealthHandlerFromEndpoint(context.Background(), mux, endpoint, []grpc.DialOption{grpc.WithInsecure()})
}

func RegisterInternalHealthService(srv *grpc.Server) {
	RegisterHealthServer(srv, ImplementedHealthService())
}

func ImplementedHealthService() HealthServer {
	return &healthService{}
}

func (*healthService) Check(context.Context, *HealthCheckRequest) (*HealthCheckResponse, error) {
	return new(HealthCheckResponse), nil
}
