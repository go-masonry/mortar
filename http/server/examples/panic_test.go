package examples

import (
	"context"
	"fmt"
	"github.com/go-masonry/mortar/http/server"
	demopackage "github.com/go-masonry/mortar/http/server/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"testing"
)

func TestPanicHandler(t *testing.T) {
	panicHandler := func(i interface{}) error {
		require.Equal(t, "ohh my god", i)
		return fmt.Errorf("panic handled")
	}
	service, err := server.Builder().
		// GRPC
		ListenOn(":8888").
		RegisterGRPCAPIs(func(srv *grpc.Server) {
			demopackage.RegisterDemoServer(srv, new(panicImpl))
		}).SetPanicHandler(panicHandler).Build()
	require.NoError(t, err)
	defer service.Stop(context.Background())
	go service.Run(context.Background()) // run service
	conn, err := grpc.Dial(":8888", grpc.WithInsecure())
	require.NoError(t, err)
	demoClient := demopackage.NewDemoClient(conn)
	_, err = demoClient.Ping(context.Background(), &demopackage.PingRequest{
		In: "ping",
	})
	assert.EqualError(t, err, "rpc error: code = Unknown desc = panic handled")
}

type panicImpl struct {
	demopackage.UnimplementedDemoServer
}

func (panicImpl) Ping(context.Context, *demopackage.PingRequest) (*demopackage.PongResponse, error) {
	panic("ohh my god")
}
