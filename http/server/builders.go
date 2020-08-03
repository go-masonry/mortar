package server

import (
	"container/list"
	"context"
	"github.com/go-masonry/mortar/interfaces/http/server"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

// ******************************************************************************************************************************************************
// ***************************************************************************REST BUILDER***************************************************************
// ******************************************************************************************************************************************************

type restConfig struct {
	addr                string
	server              *http.Server
	listener            net.Listener
	handlers            map[string]http.Handler
	handlerFuncs        map[string]http.HandlerFunc
	grpcGatewayMux      *runtime.ServeMux
	grpcGatewayHandlers []server.GRPCGatewayGeneratedHandlers
	grpcGatewayOptions  []runtime.ServeMuxOption
}

type restBuilder struct {
	parent server.GRPCWebServiceBuilder
	cfg    *restConfig
	ll     *list.List
}

func newRESTBuilder(cfg *restConfig, parent server.GRPCWebServiceBuilder) server.RESTBuilder {
	return &restBuilder{
		parent: parent,
		cfg:    cfg,
		ll:     list.New(),
	}
}

func (r *restBuilder) ListenOn(addr string) server.RESTBuilder {
	r.ll.PushBack(func(cfg *restConfig) {
		cfg.addr = addr
	})
	return r
}

func (r *restBuilder) SetCustomServer(server *http.Server) server.RESTBuilder {
	r.ll.PushBack(func(cfg *restConfig) {
		cfg.server = server
	})
	return r
}

func (r *restBuilder) SetCustomListener(listener net.Listener) server.RESTBuilder {
	r.ll.PushBack(func(cfg *restConfig) {
		cfg.listener = listener
	})
	return r
}

func (r *restBuilder) AddHandler(pattern string, handler http.Handler) server.RESTBuilder {
	r.ll.PushBack(func(cfg *restConfig) {
		if cfg.handlers == nil {
			cfg.handlers = make(map[string]http.Handler)
		}
		cfg.handlers[pattern] = handler
	})
	return r
}

func (r *restBuilder) AddHandlerFunc(pattern string, handlerFunc http.HandlerFunc) server.RESTBuilder {
	r.ll.PushBack(func(cfg *restConfig) {
		if cfg.handlerFuncs == nil {
			cfg.handlerFuncs = make(map[string]http.HandlerFunc)
		}
		cfg.handlerFuncs[pattern] = handlerFunc
	})
	return r
}

func (r *restBuilder) SetCustomGRPCGatewayMux(mux *runtime.ServeMux) server.RESTBuilder {
	r.ll.PushBack(func(cfg *restConfig) {
		cfg.grpcGatewayMux = mux
	})
	return r
}

func (r *restBuilder) RegisterGRPCGatewayHandlers(handlers ...server.GRPCGatewayGeneratedHandlers) server.RESTBuilder {
	r.ll.PushBack(func(cfg *restConfig) {
		cfg.grpcGatewayHandlers = append(cfg.grpcGatewayHandlers, handlers...)
	})
	return r
}

func (r *restBuilder) AddGRPCGatewayOptions(options ...runtime.ServeMuxOption) server.RESTBuilder {
	r.ll.PushBack(func(cfg *restConfig) {
		cfg.grpcGatewayOptions = append(cfg.grpcGatewayOptions, options...)
	})
	return r
}

func (r *restBuilder) BuildRESTPart() server.GRPCWebServiceBuilder {
	for e := r.ll.Front(); e != nil; e = e.Next() {
		f := e.Value.(func(cfg *restConfig))
		f(r.cfg)
	}
	return r.parent
}

// ******************************************************************************************************************************************************
// ***************************************************************************GRPC BUILDER***************************************************************
// ******************************************************************************************************************************************************

type grpcConfig struct {
	addr         string
	server       *grpc.Server
	listener     net.Listener
	registerApi  []server.GRPCServerAPI
	options      []grpc.ServerOption
	panicHandler func(interface{}) error
}

type webServiceConfig struct {
	grpc   *grpcConfig
	rest   []*restConfig
	logger func(ctx context.Context, format string, args ...interface{})
}

type serviceBuilder struct {
	ll *list.List
}

func Builder() server.GRPCWebServiceBuilder {
	return &serviceBuilder{ll: list.New()}
}

func (s *serviceBuilder) ListenOn(addr string) server.GRPCWebServiceBuilder {
	s.ll.PushBack(func(cfg *webServiceConfig) {
		cfg.grpc.addr = addr
	})
	return s
}

func (s *serviceBuilder) SetCustomGRPCServer(server *grpc.Server) server.GRPCWebServiceBuilder {
	s.ll.PushBack(func(cfg *webServiceConfig) {
		cfg.grpc.server = server
	})
	return s
}

func (s *serviceBuilder) SetCustomListener(listener net.Listener) server.GRPCWebServiceBuilder {
	s.ll.PushBack(func(cfg *webServiceConfig) {
		cfg.grpc.listener = listener
	})
	return s
}

func (s *serviceBuilder) RegisterGRPCAPIs(apis ...server.GRPCServerAPI) server.GRPCWebServiceBuilder {
	s.ll.PushBack(func(cfg *webServiceConfig) {
		cfg.grpc.registerApi = append(cfg.grpc.registerApi, apis...)
	})
	return s
}

func (s *serviceBuilder) AddGRPCServerOptions(options ...grpc.ServerOption) server.GRPCWebServiceBuilder {
	s.ll.PushBack(func(cfg *webServiceConfig) {
		cfg.grpc.options = append(cfg.grpc.options, options...)
	})
	return s
}

func (s *serviceBuilder) SetPanicHandler(handler func(interface{}) error) server.GRPCWebServiceBuilder {
	s.ll.PushBack(func(cfg *webServiceConfig) {
		cfg.grpc.panicHandler = handler
	})
	return s
}

func (s *serviceBuilder) SetLogger(logger func(ctx context.Context, format string, args ...interface{})) server.GRPCWebServiceBuilder {
	s.ll.PushBack(func(cfg *webServiceConfig) {
		cfg.logger = logger
	})
	return s
}

func (s *serviceBuilder) AddRESTServerConfiguration() server.RESTBuilder {
	emptyRESTConfig := new(restConfig)
	s.ll.PushBack(func(cfg *webServiceConfig) {
		cfg.rest = append(cfg.rest, emptyRESTConfig)
	})
	return newRESTBuilder(emptyRESTConfig, s)
}

func (s *serviceBuilder) Build() (server.WebService, error) {
	cfg := &webServiceConfig{
		grpc: new(grpcConfig),
	}
	for e := s.ll.Front(); e != nil; e = e.Next() {
		f := e.Value.(func(cfg *webServiceConfig))
		f(cfg)
	}
	if cfg.logger == nil {
		cfg.logger = func(context.Context, string, ...interface{}) {} // no log
	}
	if cfg.grpc.panicHandler == nil {
		cfg.grpc.panicHandler = defaultPanicHandler
	}
	cfg.grpc.options = append([]grpc.ServerOption{ // make sure they are outer most
		grpc.ChainUnaryInterceptor(panicHandlerUnaryInterceptor(cfg.grpc.panicHandler)),
		grpc.ChainStreamInterceptor(panicHandlerStreamInterceptor(cfg.grpc.panicHandler)),
	}, cfg.grpc.options...)
	return newWebService(cfg)
}

// Sanity
var _ server.GRPCWebServiceBuilder = (*serviceBuilder)(nil)
var _ server.RESTBuilder = (*restBuilder)(nil)
