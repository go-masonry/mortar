package tests

import (
	"context"

	"github.com/go-masonry/mortar/interfaces/cfg"
	mock_cfg "github.com/go-masonry/mortar/interfaces/cfg/mock"
	contextMiddleware "github.com/go-masonry/mortar/middleware/context"
	"github.com/go-masonry/mortar/mortar"
	"go.uber.org/fx"
	"google.golang.org/grpc/metadata"
)

func (s *middlewareSuite) TestLoggerGRPCIncomingContextExtractor() {
	md := metadata.MD{
		"one-of-a-kind": []string{"1"},
		"another-kind":  []string{"not extracted"},
	}
	ctxWithValues := metadata.NewIncomingContext(context.Background(), md)
	extracted := s.logExtractor(ctxWithValues)
	s.Contains(extracted, "one-of-a-kind")
	s.NotContains(extracted, "another-kind")
}

func (s *middlewareSuite) testLoggerGRPCIncomingContextExtractorBeforeTest() fx.Option {
	s.cfgMock.EXPECT().Get(mortar.MiddlewareLoggerHeaders).DoAndReturn(func(key string) cfg.Value {
		value := mock_cfg.NewMockValue(s.ctrl)
		value.EXPECT().IsSet().Return(true)
		value.EXPECT().StringSlice().Return([]string{
			"one", "two",
		})
		return value
	})
	return fx.Options(
		fx.Provide(contextMiddleware.LoggerGRPCIncomingContextExtractor),
		fx.Populate(&s.logExtractor),
	)
}
