package partial

import (
	"net/http"
	"testing"

	"github.com/go-masonry/mortar/interfaces/cfg"
	confkeys "github.com/go-masonry/mortar/interfaces/cfg/keys"
	mock_cfg "github.com/go-masonry/mortar/interfaces/cfg/mock"
	serverInt "github.com/go-masonry/mortar/interfaces/http/server"
	"github.com/go-masonry/mortar/interfaces/log"
	mock_log "github.com/go-masonry/mortar/interfaces/log/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type partialSuite struct {
	suite.Suite

	ctrl    *gomock.Controller
	cfgMock *mock_cfg.MockConfig
	logMock *mock_log.MockLogger
}

func TestPartialServer(t *testing.T) {
	suite.Run(t, new(partialSuite))
}

func (s *partialSuite) TestExternalHTTPGroups() {
	var serverBuilder serverInt.GRPCWebServiceBuilder
	testApp := fxtest.New(s.T(),
		fx.Provide(HTTPServerBuilder),
		fx.Provide(
			func() cfg.Config {
				return s.cfgMock
			},
			func() log.Logger {
				return s.logMock
			},
		),
		s.setupGroups(),
		fx.Populate(&serverBuilder),
	)
	testApp.RequireStart()
	defer testApp.RequireStop()

	// Test will make sure that we are trying to register same path twice using FX Groups
	s.PanicsWithValue("http: multiple registrations for /notfound", func() {
		serverBuilder.Build()
	})
}

func (s *partialSuite) SetupTest() {
	// This one runs before `BeforeTest`
	s.ctrl = gomock.NewController(s.T())
	s.cfgMock = mock_cfg.NewMockConfig(s.ctrl)
	s.logMock = mock_log.NewMockLogger(s.ctrl)
	// host
	s.cfgMock.EXPECT().Get(confkeys.Host).DoAndReturn(func(key string) cfg.Value {
		value := mock_cfg.NewMockValue(s.ctrl)
		value.EXPECT().String().Return("localhost")
		return value
	})
	// grpc port
	s.cfgMock.EXPECT().Get(confkeys.ExternalGRPCPort).DoAndReturn(func(key string) cfg.Value {
		value := mock_cfg.NewMockValue(s.ctrl)
		value.EXPECT().IsSet().Return(true)
		value.EXPECT().Int().Return(1234)
		return value
	})
	// host
	s.cfgMock.EXPECT().Get(confkeys.Host).DoAndReturn(func(key string) cfg.Value {
		value := mock_cfg.NewMockValue(s.ctrl)
		value.EXPECT().String().Return("localhost")
		return value
	})
	// external rest port
	s.cfgMock.EXPECT().Get(confkeys.ExternalRESTPort).DoAndReturn(func(key string) cfg.Value {
		value := mock_cfg.NewMockValue(s.ctrl)
		value.EXPECT().IsSet().Return(true)
		value.EXPECT().Int().Return(1235)
		return value
	})
	// host
	s.cfgMock.EXPECT().Get(confkeys.Host).DoAndReturn(func(key string) cfg.Value {
		value := mock_cfg.NewMockValue(s.ctrl)
		value.EXPECT().String().Return("localhost")
		return value
	})
	// internal rest port
	s.cfgMock.EXPECT().Get(confkeys.InternalRESTPort).DoAndReturn(func(key string) cfg.Value {
		value := mock_cfg.NewMockValue(s.ctrl)
		value.EXPECT().IsSet().Return(true)
		return value
	})
}

func (s *partialSuite) setupGroups() fx.Option {
	return fx.Provide(
		// grpc
		fx.Annotated{
			Group: FxGroupGRPCServerAPIs,
			Target: func() serverInt.GRPCServerAPI {
				return func(srv *grpc.Server) {
					grpc_health_v1.RegisterHealthServer(srv, health.NewServer())
				}
			},
		},
		fx.Annotated{
			Group: FxGroupExternalHTTPHandlers + ",flatten",
			Target: func() []HTTPHandlerPatternPair {
				return []HTTPHandlerPatternPair{
					{Pattern: "/notfound", Handler: http.NotFoundHandler()},
				}
			},
		},
		fx.Annotated{
			Group: FxGroupExternalHTTPHandlerFunctions,
			Target: func() HTTPHandlerFuncPatternPair {
				return HTTPHandlerFuncPatternPair{
					Pattern: "/notfound", HandlerFunc: http.NotFound, // the same as above should return error
				}
			},
		},
	)
}
