package tests

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-masonry/mortar/interfaces/cfg"
	confkeys "github.com/go-masonry/mortar/interfaces/cfg/keys"
	mock_cfg "github.com/go-masonry/mortar/interfaces/cfg/mock"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/go-masonry/mortar/interfaces/monitor"
	mock_monitor "github.com/go-masonry/mortar/interfaces/monitor/mock"
	"github.com/go-masonry/mortar/logger/naive"
	"github.com/go-masonry/mortar/middleware/interceptors/client"
	"github.com/golang/mock/gomock"
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
		}
		return fmt.Errorf("assertions failed")
	}
	err := s.clientInterceptor(ctxWithIncoming, "", nil, nil, nil, invoker)
	s.NoError(err)
}
func (s *middlewareSuite) TestHTTPClientInterceptorHeaderCopier() {
	md := metadata.MD{
		"one-of-a-kind": []string{"1", "2"},
		"another-kind":  []string{"not extracted"},
	}
	ctxWithIncoming := metadata.NewIncomingContext(context.Background(), md)
	fakeHandler := func(req *http.Request) (*http.Response, error) {
		if !(s.NotEmpty(req.Header, "not exist/empty 'one-of-a-kind' header") &&
			s.Contains(req.Header, "One-Of-A-Kind") &&
			s.NotContains(req.Header, "another-kind") &&
			s.ElementsMatch([]string{"1", "2"}, req.Header.Values("one-of-a-kind"))) {
			return nil, fmt.Errorf("assertions failed")
		}
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
	req = req.WithContext(ctxWithIncoming)
	_, err := s.restClientInterceptor(req, fakeHandler)
	s.NoError(err)
}

func (s *middlewareSuite) testClientInterceptorHeaderCopierBeforeTest() fx.Option {
	s.cfgMock.EXPECT().Get(confkeys.ForwardIncomingGRPCMetadataHeadersList).DoAndReturn(func(key string) cfg.Value {
		value := mock_cfg.NewMockValue(s.ctrl)
		value.EXPECT().StringSlice().Return([]string{
			"one", "two",
		})
		return value
	})
	return fx.Options(
		fx.Provide(client.CopyGRPCHeadersClientInterceptor),
		fx.Provide(client.CopyGRPCHeadersHTTPClientInterceptor),
		fx.Populate(&s.clientInterceptor),
		fx.Populate(&s.restClientInterceptor),
	)
}

func (s *middlewareSuite) TestDumpRESTClientInterceptor() {
	fakeHandler := func(req *http.Request) (*http.Response, error) {
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
	res, err := s.restClientInterceptor(req, fakeHandler)
	s.NoError(err)
	defer res.Body.Close()
	s.Contains(s.loggerOutput.String(), "Request:\nGET /path HTTP/1.1\r\nHost: somewhere") // req path
	s.Contains(s.loggerOutput.String(), "foo")                                             // response body
	s.Equal(http.StatusOK, res.StatusCode)
}

func (s *middlewareSuite) testDumpRESTClientInterceptorBeforeTest() fx.Option {
	return fx.Options(
		fx.Provide(client.DumpRESTClientInterceptor),
		fx.Provide(func() log.Logger {
			return naive.Builder().SetWriter(&s.loggerOutput).SetLevel(log.DebugLevel).Build()
		}),
		fx.Populate(&s.restClientInterceptor),
	)
}

func (s *middlewareSuite) TestRESTClientMetrics() {
	mockTimer := mock_monitor.NewMockTagsAwareTimer(s.ctrl)
	mockTimer.EXPECT().Record(gomock.Any()).After(
		mockTimer.EXPECT().WithContext(gomock.Any()).Return(mockTimer),
	)
	s.metricsMock.EXPECT().Timer("client_calls_duration", gomock.Any()).Return(mockTimer).After(
		s.metricsMock.EXPECT().WithTags(gomock.Any()).DoAndReturn(func(tags monitor.Tags) *mock_monitor.MockMetrics {
			s.Equal(monitor.Tags{
				"target": "wonder",
				"path":   "/land",
				"err":    "false",
				"ctype":  "rest",
			}, tags)
			return s.metricsMock
		}),
	)
	fakeHandler := func(req *http.Request) (*http.Response, error) {
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
	req, _ := http.NewRequestWithContext(context.TODO(), http.MethodHead, "http://wonder:/land", nil)
	res, err := s.restClientInterceptor(req, fakeHandler)
	s.NoError(err)
	defer res.Body.Close()
}

func (s *middlewareSuite) TestGRPCClientMetrics() {
	mockTimer := mock_monitor.NewMockTagsAwareTimer(s.ctrl)
	mockTimer.EXPECT().Record(gomock.Any()).After(
		mockTimer.EXPECT().WithContext(gomock.Any()).Return(mockTimer),
	)
	s.metricsMock.EXPECT().Timer("client_calls_duration", gomock.Any()).Return(mockTimer).After(
		s.metricsMock.EXPECT().WithTags(gomock.Any()).DoAndReturn(func(tags monitor.Tags) *mock_monitor.MockMetrics {
			s.Equal(monitor.Tags{
				"target": "",
				"path":   "/wonder.Land/Thing",
				"err":    "true",
				"ctype":  "grpc",
			}, tags)
			return s.metricsMock
		}),
	)
	invoker := func(_ context.Context, _ string, _, _ interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return errors.New("fake error")
	}
	err := s.clientInterceptor(context.TODO(), "/wonder.Land/Thing", nil, nil, &grpc.ClientConn{}, invoker)
	s.Error(err)
}

func (s *middlewareSuite) testClientMetricsBeforeTest() fx.Option {
	return fx.Options(
		fx.Provide(func() monitor.Metrics {
			return s.metricsMock
		}),
		fx.Provide(client.MonitorRESTClientCallsInterceptor),
		fx.Provide(client.MonitorGRPCClientCallsInterceptor),
		fx.Provide(func() log.Logger { return naive.Builder().SetWriter(&s.loggerOutput).Build() }),
		fx.Populate(
			&s.restClientInterceptor,
			&s.clientInterceptor,
		),
	)
}
