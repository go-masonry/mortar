package health

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type healthService struct {
	UnimplementedHealthServer
}

// RegisterInternalGRPCGatewayHandler grpc-gateway health handler
func RegisterInternalGRPCGatewayHandler(mux *runtime.ServeMux, endpoint string) error {
	return RegisterHealthHandlerFromEndpoint(context.Background(), mux, endpoint, []grpc.DialOption{grpc.WithInsecure()})
}

// RegisterInternalHealthService grpc server health api registration
func RegisterInternalHealthService(srv *grpc.Server) {
	RegisterHealthServer(srv, ImplementedHealthService())
}

// ImplementedHealthService internal health service
func ImplementedHealthService() HealthServer {
	return &healthService{}
}

func (*healthService) Check(context.Context, *HealthCheckRequest) (*HealthCheckResponse, error) {
	return new(HealthCheckResponse), nil
}
