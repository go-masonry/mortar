package server

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	demopackage "github.com/go-masonry/mortar/http/server/proto"
	serverInt "github.com/go-masonry/mortar/interfaces/http/server"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

type multiListenersSuite struct {
	suite.Suite
	builder serverInt.GRPCWebServiceBuilder
	web     serverInt.WebService
}

func TestSeveralListeners(t *testing.T) {
	suite.Run(t, new(multiListenersSuite))
}

func (ms *multiListenersSuite) TestPingRest() {
	response, err := http.Get("http://localhost:8889/v1/demo/ping")
	ms.Require().NoError(err)
	ms.Require().Equal(http.StatusOK, response.StatusCode)
	defer response.Body.Close()
	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	ms.Require().NoError(err)
	ms.Require().Contains(result, "out")
	ms.Require().Equal("-pong", result["out"])
}

func (ms *multiListenersSuite) TestPingGRPC() {
	conn, err := grpc.Dial(":8888", grpc.WithInsecure())
	ms.Require().NoError(err)
	demoClient := demopackage.NewDemoClient(conn)
	response, err := demoClient.Ping(context.Background(), &demopackage.PingRequest{
		In: "ping",
	})
	ms.Require().NoError(err)
	ms.Require().Equal("ping-pong", response.GetOut())
}

func (ms *multiListenersSuite) TestCustomPath() {
	response, err := http.Get("http://localhost:8890/custom/path")
	ms.NoError(err)
	ms.Equal(http.StatusAccepted, response.StatusCode)
	defer response.Body.Close()
	all, err := ioutil.ReadAll(response.Body)
	ms.NoError(err)
	ms.Equal("it's custom", string(all))
}

func (ms *multiListenersSuite) SetupSuite() {
	ms.builder = Builder().
		SetLogger(func(ctx context.Context, format string, args ...interface{}) {
			ms.T().Logf(format, args...)
		}).
		// GRPC
		ListenOn(":8888").
		RegisterGRPCAPIs(func(srv *grpc.Server) {
			demopackage.RegisterDemoServer(srv, new(demoImpl))
		}).
		// REST 1 with GRPC Gateway
		AddRESTServerConfiguration().
		ListenOn(":8889").
		RegisterGRPCGatewayHandlers(func(mux *runtime.ServeMux, endpoint string) error {
			return demopackage.RegisterDemoHandlerFromEndpoint(context.Background(), mux, endpoint, []grpc.DialOption{grpc.WithInsecure()})
		}).BuildRESTPart().
		// REST 2 without GRPC Gateway
		AddRESTServerConfiguration().
		ListenOn(":8890").
		AddHandlerFunc("/custom/path", customHandler).
		BuildRESTPart()
}

func (ms *multiListenersSuite) BeforeTest(_, _ string) {
	go func() {
		var err error
		ms.web, err = ms.builder.Build()
		if err != nil {
			ms.FailNow("error during setup", "%v", err)
		}
		ms.web.Run(context.Background())
	}()
	time.Sleep(500 * time.Millisecond) // compensate on build machine
}

func (ms *multiListenersSuite) AfterTest(_, _ string) {
	ms.web.Stop(context.Background())
}

type demoImpl struct {
	demopackage.UnimplementedDemoServer
}

func (d demoImpl) Ping(ctx context.Context, request *demopackage.PingRequest) (*demopackage.PongResponse, error) {
	return &demopackage.PongResponse{Out: request.GetIn() + "-pong"}, nil
}

func customHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("it's custom"))
}
