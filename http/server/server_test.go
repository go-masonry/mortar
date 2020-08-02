package server

import (
	"context"
	demopackage "github.com/go-masonry/mortar/http/server/proto"
	"github.com/go-masonry/mortar/interfaces/http/server"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"net/http"
	"testing"
)

func TestDefaultListeners(t *testing.T) {
	service, err := Builder().
		RegisterGRPCAPIs(registerGrpcAPI).
		AddRESTServerConfiguration().
		AddGRPCGatewayHandlers(registerGatewayHandler).
		BuildRESTPart().
		Build()
	require.NoError(t, err)
	var ports = make([]int, 2)
	for _, serviceType := range service.Ports() {
		if serviceType.Type == server.GRPCServer {
			ports[0]++
		} else {
			ports[1]++
		}
	}
	assert.ElementsMatch(t, []int{1, 1}, ports) // one for grpc and one for rest
}

func TestListenOnAddresses(t *testing.T) {
	service, err := Builder().
		ListenOn(":8888").
		RegisterGRPCAPIs(registerGrpcAPI).
		AddRESTServerConfiguration().
		ListenOn(":8887").
		AddGRPCGatewayHandlers(registerGatewayHandler).
		BuildRESTPart().
		AddRESTServerConfiguration().
		ListenOn(":8889").
		AddGRPCGatewayHandlers(registerGatewayHandler).
		BuildRESTPart().
		Build()
	require.NoError(t, err)
	defer service.Stop(context.Background()) // make sure listeners are closed
	for _, info := range service.Ports() {
		if info.Type == server.GRPCServer {
			assert.Equal(t, 8888, info.Port, "GRPC port is wrong")
		} else {
			assert.True(t, info.Port == 8887 || info.Port == 8889, "REST ports are wrong")
		}
	}
	assert.Len(t, service.Ports(), 3, "amount of open ports is wrong")
}

func TestCustomListeners(t *testing.T) {
	grpcL, _ := net.Listen("tcp", ":8888")
	defer grpcL.Close()
	restL, _ := net.Listen("tcp", ":8889")
	defer restL.Close()
	service, err := Builder().
		SetCustomListener(grpcL).
		RegisterGRPCAPIs(registerGrpcAPI).
		AddRESTServerConfiguration().
		SetCustomListener(restL).
		AddGRPCGatewayHandlers(registerGatewayHandler).
		BuildRESTPart().
		Build()
	require.NoError(t, err)
	for _, info := range service.Ports() {
		if info.Type == server.GRPCServer {
			assert.Equal(t, 8888, info.Port, "GRPC port is wrong")
		} else {
			assert.Equal(t, 8889, info.Port, "REST port is wrong")
		}
	}
}

func TestSettingCustomGrpcServer(t *testing.T) {
	grpcServer := grpc.NewServer()
	service, err := Builder().SetCustomGRPCServer(grpcServer).RegisterGRPCAPIs(registerGrpcAPI).Build()
	require.NoError(t, err)
	defer service.Stop(context.Background())
	instance, ok := service.(*webService)
	require.True(t, ok, "wrong web service implementation")
	assert.Equal(t, grpcServer, instance.grpcServer)
}

func TestSettingCustomRESTServer(t *testing.T) {
	restServer := &http.Server{}
	service, err := Builder().RegisterGRPCAPIs(registerGrpcAPI).
		AddRESTServerConfiguration().SetCustomServer(restServer).AddGRPCGatewayHandlers(registerGatewayHandler).BuildRESTPart().
		Build()
	require.NoError(t, err)
	defer service.Stop(context.Background())
	instance, ok := service.(*webService)
	require.True(t, ok, "wrong web service implementation")
	assert.Equal(t, restServer, instance.muxAndListeners[1].m) // little hacky I know
}

func TestInternalHealthService(t *testing.T) {
	service, err := Builder().ListenOn(":8888").RegisterGRPCAPIs(registerGrpcAPI).Build()
	require.NoError(t, err)
	defer service.Stop(context.Background())
	go service.Run(context.Background())
	conn, err := grpc.Dial(":8888", grpc.WithInsecure())
	defer conn.Close()
	require.NoError(t, err)
	healthClient := grpc_health_v1.NewHealthClient(conn)
	response, err := healthClient.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{Service: ""})
	require.NoError(t, err)
	assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, response.Status)
}

func registerGrpcAPI(srv *grpc.Server) {
	demopackage.RegisterDemoServer(srv, new(demopackage.UnimplementedDemoServer))
}

func registerGatewayHandler(mux *runtime.ServeMux, endpoint string) error {
	return demopackage.RegisterDemoHandlerFromEndpoint(context.Background(), mux, endpoint, []grpc.DialOption{grpc.WithInsecure()})
}
