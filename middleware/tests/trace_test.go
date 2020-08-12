package tests

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-masonry/mortar/interfaces/cfg"
	mock_cfg "github.com/go-masonry/mortar/interfaces/cfg/mock"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/go-masonry/mortar/logger"
	"github.com/go-masonry/mortar/middleware/interceptors/trace"
	"github.com/go-masonry/mortar/mortar"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func (s *middlewareSuite) TestTracerGRPCClientInterceptor() {
	tracerMock := s.tracer.(*mocktracer.MockTracer)
	invoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return nil
	}
	err := s.clientInterceptor(context.Background(), "package.service/method", "request", nil, nil, invoker)
	s.NoError(err)
	spans := tracerMock.FinishedSpans()
	s.Require().Len(spans, 1)
	clientSpan := spans[0]
	s.Equal("gRPC", clientSpan.Tag("component"))
	s.EqualValues("client", clientSpan.Tag("span.kind"))
	s.Equal("package.service/method", clientSpan.OperationName)
	logs := clientSpan.Logs()
	s.Len(logs, 2, "request or response is missing")
	// No errors during injecting
	s.Empty(s.loggerOutput.String())
}

func (s *middlewareSuite) testTracerGRPCClientInterceptorBeforeTest() fx.Option {
	s.cfgMock.EXPECT().Get(mortar.MiddlewareClientGRPCTraceIncludeRequest).DoAndReturn(func(key string) cfg.Value {
		value := mock_cfg.NewMockValue(s.ctrl)
		value.EXPECT().Bool().Return(true)
		return value
	})
	s.cfgMock.EXPECT().Get(mortar.MiddlewareClientGRPCTraceIncludeResponse).DoAndReturn(func(key string) cfg.Value {
		value := mock_cfg.NewMockValue(s.ctrl)
		value.EXPECT().Bool().Return(true)
		return value
	})
	return fx.Options(
		s.unifiedOptionsForTraceInterceptors(),
		fx.Provide(trace.TracerGRPCClientInterceptor),
		fx.Populate(&s.clientInterceptor),
	)
}

func (s *middlewareSuite) TestTracerRESTClientInterceptor() {
	handler := func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			Status:        "200 OK",
			StatusCode:    200,
			Proto:         "HTTP/1.1",
			ProtoMajor:    1,
			ProtoMinor:    1,
			ContentLength: 3,
			Body:          ioutil.NopCloser(strings.NewReader("foo")),
		}, nil
	}
	req, _ := http.NewRequest(http.MethodGet, "http://somewhere/path", nil)
	_, err := s.restClientInterceptor(req, handler)
	s.NoError(err)
	tracerMock := s.tracer.(*mocktracer.MockTracer)
	spans := tracerMock.FinishedSpans()
	s.Require().Len(spans, 1)
	clientSpan := spans[0]
	s.Equal("REST", clientSpan.Tag("component"))
	s.EqualValues("client", clientSpan.Tag("span.kind"))
	s.Equal("/path", clientSpan.OperationName)
	logs := clientSpan.Logs()
	s.Len(logs, 2, "request or response is missing")
	// No errors during injecting
	s.Empty(s.loggerOutput.String())
}

func (s *middlewareSuite) testTracerRESTClientInterceptorBeforeTest() fx.Option {
	s.cfgMock.EXPECT().Get(mortar.MiddlewareClientRESTTraceIncludeResponse).DoAndReturn(func(key string) cfg.Value {
		value := mock_cfg.NewMockValue(s.ctrl)
		value.EXPECT().Bool().Return(true)
		return value
	})
	s.cfgMock.EXPECT().Get(mortar.MiddlewareClientRESTTraceIncludeRequest).DoAndReturn(func(key string) cfg.Value {
		value := mock_cfg.NewMockValue(s.ctrl)
		value.EXPECT().Bool().Return(true)
		return value
	})
	return fx.Options(
		s.unifiedOptionsForTraceInterceptors(),
		fx.Provide(trace.TracerRESTClientInterceptor),
		fx.Populate(&s.restClientInterceptor),
	)
}

func (s *middlewareSuite) TestGRPCTracingUnaryServerInterceptor() {
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "response", nil
	}
	response, err := s.serverInterceptor(context.Background(), "request", &grpc.UnaryServerInfo{FullMethod: "package.service/method"}, handler)
	s.Equal("response", response)
	s.NoError(err)
	tracerMock := s.tracer.(*mocktracer.MockTracer)
	spans := tracerMock.FinishedSpans()
	s.Require().Len(spans, 1)
	clientSpan := spans[0]
	s.Equal("gRPC", clientSpan.Tag("component"))
	s.EqualValues("server", clientSpan.Tag("span.kind"))
	s.Equal("package.service/method", clientSpan.OperationName)
	logs := clientSpan.Logs()
	s.Len(logs, 2, "request or response is missing")
	// No errors during injecting
	s.Empty(s.loggerOutput.String())
}

func (s *middlewareSuite) testGRPCTracingUnaryServerInterceptorBeforeTest() fx.Option {
	s.cfgMock.EXPECT().Get(mortar.MiddlewareServerGRPCTraceIncludeRequest).DoAndReturn(func(key string) cfg.Value {
		value := mock_cfg.NewMockValue(s.ctrl)
		value.EXPECT().Bool().Return(true)
		return value
	})
	s.cfgMock.EXPECT().Get(mortar.MiddlewareServerGRPCTraceIncludeResponse).DoAndReturn(func(key string) cfg.Value {
		value := mock_cfg.NewMockValue(s.ctrl)
		value.EXPECT().Bool().Return(true)
		return value
	})
	return fx.Options(
		s.unifiedOptionsForTraceInterceptors(),
		fx.Provide(trace.GRPCTracingUnaryServerInterceptor),
		fx.Populate(&s.serverInterceptor),
	)
}

func (s *middlewareSuite) unifiedOptionsForTraceInterceptors() fx.Option {
	return fx.Options(
		fx.Provide(func() log.Logger {
			return logger.Builder().SetWriter(&s.loggerOutput).SetLevel(log.TraceLevel).Build()
		}),
		fx.Provide(func() opentracing.Tracer {
			return mocktracer.New()
		}),
		fx.Populate(&s.tracer),
	)
}
