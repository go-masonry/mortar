package server

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"testing"
	"time"

	demopackage "github.com/go-masonry/mortar/http/server/proto"
	serverInt "github.com/go-masonry/mortar/interfaces/http/server"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/soheilhy/cmux"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

type onePortSuite struct {
	suite.Suite
	builder serverInt.GRPCWebServiceBuilder
	web     serverInt.WebService
	cMux    cmux.CMux
}

func TestOnePortListener(t *testing.T) {
	suite.Run(t, new(onePortSuite))
}

func (os *onePortSuite) TestPingRestOnSamePort() {
	response, err := http.Get("http://localhost:8888/v1/demo/ping")
	os.NoError(err)
	os.Equal(http.StatusOK, response.StatusCode)
	defer response.Body.Close()
	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	os.NoError(err)
	os.Contains(result, "out")
	os.Equal("-pong", result["out"])
}

func (os *onePortSuite) TestPingGRPCOnSamePort() {
	conn, err := grpc.Dial(":8888", grpc.WithInsecure())
	os.NoError(err)
	demoClient := demopackage.NewDemoClient(conn)
	response, err := demoClient.Ping(context.Background(), &demopackage.PingRequest{
		In: "ping",
	})
	os.NoError(err)
	os.Equal("ping-pong", response.GetOut())
}

func (os *onePortSuite) makeBuilderAndCmux() {
	listener, err := net.Listen("tcp", "localhost:8888")
	os.Require().NoError(err)
	os.cMux = cmux.New(listener)
	restL := os.cMux.Match(cmux.HTTP1())
	grpcL := os.cMux.Match(cmux.Any())

	os.builder = Builder().
		SetLogger(func(ctx context.Context, format string, args ...interface{}) {
			os.T().Logf(format, args...)
		}).
		// GRPC
		SetCustomListener(grpcL).
		RegisterGRPCAPIs(func(srv *grpc.Server) {
			demopackage.RegisterDemoServer(srv, new(demoImpl))
		}).
		// REST 1 with GRPC Gateway
		AddRESTServerConfiguration().
		SetCustomListener(restL).
		RegisterGRPCGatewayHandlers(func(mux *runtime.ServeMux, endpoint string) error {
			return demopackage.RegisterDemoHandlerFromEndpoint(context.Background(), mux, endpoint, []grpc.DialOption{grpc.WithInsecure()})
		}).
		BuildRESTPart()
}

func (os *onePortSuite) BeforeTest(_, _ string) {
	os.makeBuilderAndCmux()
	go func() {
		var err error
		os.web, err = os.builder.Build()
		if err != nil {
			os.FailNow("error during setup", "%v", err)
		}
		go os.cMux.Serve()               // cmux must serve on it's own since it's blocking
		os.web.Run(context.Background()) // this one is blocking
	}()
	time.Sleep(500 * time.Millisecond) // compensate on build machine
}

func (os *onePortSuite) AfterTest(_, _ string) {
	os.web.Stop(context.Background())
	// cmux listener will close since we are calling it's "child" to close :)
}
