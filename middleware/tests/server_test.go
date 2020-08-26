package tests

import (
	"context"
	"fmt"
	"time"

	"github.com/go-masonry/mortar/logger/naive"
	"github.com/golang/mock/gomock"

	"github.com/go-masonry/mortar/interfaces/cfg"
	mock_cfg "github.com/go-masonry/mortar/interfaces/cfg/mock"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/go-masonry/mortar/interfaces/monitor"
	mock_monitor "github.com/go-masonry/mortar/interfaces/monitor/mock"
	"github.com/go-masonry/mortar/middleware/interceptors/server"
	"github.com/go-masonry/mortar/mortar"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func (s *middlewareSuite) TestLoggerGRPCInterceptor() {
	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "response", nil
	}
	ctxWithDeadline, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := s.serverInterceptor(ctxWithDeadline, nil, &grpc.UnaryServerInfo{FullMethod: "fake method"}, unaryHandler)
	s.NoError(err)
	s.Contains(s.loggerOutput.String(), "fake method finished")
}

func (s *middlewareSuite) testLoggerGRPCInterceptorBeforeTest() fx.Option {
	s.cfgMock.EXPECT().Get(mortar.ServerGRPCLogLevel).DoAndReturn(func(key string) cfg.Value {
		value := mock_cfg.NewMockValue(s.ctrl)
		value.EXPECT().IsSet().Return(true)
		value.EXPECT().String().Return("debug")
		return value
	})

	s.cfgMock.EXPECT().Get(mortar.MiddlewareServerGRPCLogIncludeRequest).DoAndReturn(func(key string) cfg.Value {
		value := mock_cfg.NewMockValue(s.ctrl)
		value.EXPECT().Bool().Return(true)
		return value
	})

	s.cfgMock.EXPECT().Get(mortar.MiddlewareServerGRPCLogIncludeResponse).DoAndReturn(func(key string) cfg.Value {
		value := mock_cfg.NewMockValue(s.ctrl)
		value.EXPECT().Bool().Return(true)
		return value
	})

	return fx.Options(
		fx.Provide(server.LoggerGRPCInterceptor),
		fx.Provide(func() log.Logger {
			return naive.Builder().SetWriter(&s.loggerOutput).SetLevel(log.DebugLevel).Build()
		}),
		fx.Populate(&s.serverInterceptor),
	)
}

func (s *middlewareSuite) TestMonitorGRPCInterceptor() {
	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "response", fmt.Errorf("some error")
	}
	_, err := s.serverInterceptor(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "package.service/method"}, unaryHandler)
	s.Error(err)
}

func (s *middlewareSuite) testMonitorGRPCInterceptorBeforeTest() fx.Option {
	return fx.Options(
		fx.Provide(server.MonitorGRPCInterceptor),
		fx.Provide(func() log.Logger {
			return naive.Builder().SetWriter(&s.loggerOutput).SetLevel(log.DebugLevel).Build()
		}),
		fx.Provide(func() monitor.Metrics {
			mockedTimer := mock_monitor.NewMockTagsAwareTimer(s.ctrl)
			mockedTimer.EXPECT().Record(gomock.AssignableToTypeOf(time.Second)) // assignable to duration
			mockMetrics := mock_monitor.NewMockMetrics(s.ctrl)
			mockMetrics.EXPECT().WithTags(monitor.Tags{"code": "2"}).Return(mockMetrics)
			mockMetrics.EXPECT().Timer("grpc_method", gomock.Any()).Return(mockedTimer) // method is from the above unary info
			return mockMetrics
		}),
		fx.Populate(&s.serverInterceptor),
	)
}
