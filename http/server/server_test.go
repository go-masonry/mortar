package server

import (
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/go-masonry/mortar/http/server/health"
	demopackage "github.com/go-masonry/mortar/http/server/proto"
	"github.com/go-masonry/mortar/interfaces/http/server"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDefaultListeners(t *testing.T) {
	service, err := Builder().
		RegisterGRPCAPIs(registerGrpcAPI).
		AddRESTServerConfiguration().
		RegisterGRPCGatewayHandlers(registerGatewayHandler).
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
		RegisterGRPCGatewayHandlers(registerGatewayHandler).
		BuildRESTPart().
		AddRESTServerConfiguration().
		ListenOn(":8889").
		RegisterGRPCGatewayHandlers(registerGatewayHandler).
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
		RegisterGRPCGatewayHandlers(registerGatewayHandler).
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
		AddRESTServerConfiguration().SetCustomServer(restServer).RegisterGRPCGatewayHandlers(registerGatewayHandler).BuildRESTPart().
		Build()
	require.NoError(t, err)
	defer service.Stop(context.Background())
	instance, ok := service.(*webService)
	require.True(t, ok, "wrong web service implementation")
	assert.Equal(t, restServer, instance.muxAndListeners[1].m) // little hacky I know
}

func TestInternalHealthService(t *testing.T) {
	service, err := Builder().ListenOn(":8888").RegisterGRPCAPIs(registerGrpcAPI).RegisterGRPCAPIs(health.RegisterInternalHealthService).Build()
	require.NoError(t, err)
	defer service.Stop(context.Background())
	go service.Run(context.Background())
	conn, err := grpc.Dial(":8888", grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()
	healthClient := health.NewHealthClient(conn)
	_, err = healthClient.Check(context.Background(), &health.HealthCheckRequest{})
	require.NoError(t, err)
}

func TestCustomGrpcServerOptions(t *testing.T) {
	service, err := Builder().
		ListenOn(":8888").
		AddGRPCServerOptions(grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			return nil, status.Error(codes.OutOfRange, "out of range")
		})).
		RegisterGRPCAPIs(registerGrpcAPI).Build()
	require.NoError(t, err)
	defer service.Stop(context.Background())
	go service.Run(context.Background())
	conn, err := grpc.Dial(":8888", grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()
	demoClient := demopackage.NewDemoClient(conn)
	_, err = demoClient.Ping(context.Background(), &demopackage.PingRequest{In: "in"})
	assert.EqualError(t, err, "rpc error: code = OutOfRange desc = out of range")
}

func TestAddCustomHandler(t *testing.T) {
	service, err := Builder().
		ListenOn(":8888").
		RegisterGRPCAPIs(registerGrpcAPI).
		AddRESTServerConfiguration().
		ListenOn(":8889").
		AddHandler("/notfound", http.NotFoundHandler()).
		BuildRESTPart().
		Build()
	require.NoError(t, err)
	defer service.Stop(context.Background())
	go service.Run(context.Background())
	resp, err := http.Get("http://localhost:8889/notfound")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestCustomGrpcGatewayOptions(t *testing.T) {
	service, err := Builder().
		ListenOn(":8888").
		RegisterGRPCAPIs(registerGrpcAPI).
		AddRESTServerConfiguration().
		ListenOn(":8889").
		RegisterGRPCGatewayHandlers(registerGatewayHandler).
		AddGRPCGatewayOptions(runtime.WithErrorHandler(func(_ context.Context, _ *runtime.ServeMux, _ runtime.Marshaler, writer http.ResponseWriter, _ *http.Request, _ error) {
			http.Error(writer, "bad one", http.StatusTeapot)
		})).
		BuildRESTPart().
		Build()
	require.NoError(t, err)
	defer service.Stop(context.Background()) // clean
	go service.Run(context.Background())
	resp, err := http.Get("http://localhost:8889/v1/demo/ping")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusTeapot, resp.StatusCode)
	bytes, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "bad one", strings.TrimSpace(string(bytes)))
}

func TestCustomGrpcGatewayMux(t *testing.T) {
	service, err := Builder().
		ListenOn(":8888").
		RegisterGRPCAPIs(registerGrpcAPI).
		AddRESTServerConfiguration().
		ListenOn(":8889").
		RegisterGRPCGatewayHandlers(registerGatewayHandler).
		SetCustomGRPCGatewayMux(runtime.NewServeMux(runtime.WithErrorHandler(func(_ context.Context, _ *runtime.ServeMux, _ runtime.Marshaler, writer http.ResponseWriter, _ *http.Request, _ error) {
			http.Error(writer, "bad one", http.StatusTeapot)
		}))).
		BuildRESTPart().
		Build()
	require.NoError(t, err)
	defer service.Stop(context.Background()) // clean
	go service.Run(context.Background())
	resp, err := http.Get("http://localhost:8889/v1/demo/ping")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusTeapot, resp.StatusCode)
	bytes, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "bad one", strings.TrimSpace(string(bytes)))
}

func registerGrpcAPI(srv *grpc.Server) {
	demopackage.RegisterDemoServer(srv, new(demopackage.UnimplementedDemoServer))
}

func registerGatewayHandler(mux *runtime.ServeMux, endpoint string) error {
	return demopackage.RegisterDemoHandlerFromEndpoint(context.Background(), mux, endpoint, []grpc.DialOption{grpc.WithInsecure()})
}
