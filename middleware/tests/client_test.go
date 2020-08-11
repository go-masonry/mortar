package tests

import (
	"context"
	"fmt"
	"github.com/go-masonry/mortar/interfaces/cfg"
	mock_cfg "github.com/go-masonry/mortar/interfaces/cfg/mock"
	"github.com/go-masonry/mortar/middleware/interceptors/client"
	"github.com/go-masonry/mortar/mortar"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (s *middlewareSuite) TestClientInterceptorHeaderCopier() {
	md := metadata.MD{
		"one-of-a-kind": []string{"1"},
		"another-kind":  []string{"not extracted"},
	}
	ctxWithIncoming := metadata.NewIncomingContext(context.Background(), md)
	invoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		outgoingContext, b := metadata.FromOutgoingContext(ctx)
		result := s.True(b, "no outgoing md found") &&
			s.Contains(outgoingContext, "one-of-a-kind") &&
			s.NotContains(outgoingContext, "another-kind")
		if result {
			return nil
		} else {
			return fmt.Errorf("assertions failed")
		}
	}
	err := s.clientInterceptor(ctxWithIncoming, "", nil, nil, nil, invoker)
	s.NoError(err)
}

func (s *middlewareSuite) testClientInterceptorHeaderCopierBeforeTest() fx.Option {
	s.cfgMock.EXPECT().Get(mortar.MiddlewareServerGRPCCopyHeadersPrefixes).DoAndReturn(func(key string) cfg.Value {
		value := mock_cfg.NewMockValue(s.ctrl)
		value.EXPECT().StringSlice().Return([]string{
			"one", "two",
		})
		return value
	})
	return fx.Options(
		fx.Provide(client.CopyGRPCHeadersClientInterceptor),
		fx.Populate(&s.clientInterceptor),
	)
}
