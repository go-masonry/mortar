package partial

import (
	"fmt"
	"github.com/go-masonry/mortar/health"
	"github.com/go-masonry/mortar/http/server"
	"github.com/go-masonry/mortar/interfaces/cfg"
	serverInt "github.com/go-masonry/mortar/interfaces/http/server"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/go-masonry/mortar/mortar"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"net/http"
)

// Group order is not guaranteed, if it's important then add them manually
const (
	FxGroupGRPCServerAPIs               = "grpcServerAPIs"
	FxGroupGRPCGatewayGeneratedHandlers = "grpcGatewayGeneratedHandlers"
	FxGroupGRPCGatewayMuxOptions        = "grpcGatewayMuxOptions"
	FxGroupUnaryServerInterceptors      = "unaryServerInterceptors"
	FxGroupInternalHttpHandlers         = "internalHttpHandlers"
	FxGroupInternalHttpHandlerFunctions = "internalHttpHandlerFunctions"
)

type HttpHandlerPatternPair struct {
	Pattern string
	Handler http.Handler
}

type HttpHandlerFuncPatternPair struct {
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type httpServerDeps struct {
	fx.In

	Config cfg.Config
	Logger log.Logger
	// GRPC
	GRPCServerAPIs    []serverInt.GRPCServerAPI     `group:"grpcServerAPIs"`
	UnaryInterceptors []grpc.UnaryServerInterceptor `group:"unaryServerInterceptors"`
	// External REST
	GRPCGatewayGeneratedHandlers []serverInt.GRPCGatewayGeneratedHandlers `group:"grpcGatewayGeneratedHandlers"`
	GRPCGatewayMuxOptions        []runtime.ServeMuxOption                 `group:"grpcGatewayMuxOptions"`
	// Internal REST
	InternalHttpHandlers         []HttpHandlerPatternPair     `group:"internalHttpHandlers"`
	InternalHttpHandlerFunctions []HttpHandlerFuncPatternPair `group:"internalHttpHandlerFunctions"`
}

// HttpServerBuilder true to it's name, it is partially initialized builder.
//
// It uses some default assumptions and configurations, which are mostly good.
// However, if you need to customize your configuration it's better to build yours from scratch
//
func HttpServerBuilder(deps httpServerDeps) serverInt.GRPCWebServiceBuilder {
	builder := server.Builder().SetLogger(deps.Logger.Debug)
	// GRPC port
	if grpcPort := deps.Config.Get(mortar.ServerGRPCPort); grpcPort.IsSet() {
		builder = builder.ListenOn(fmt.Sprintf(":%d", grpcPort.Int()))
	}
	// GRPC server interceptors
	if len(deps.UnaryInterceptors) > 0 {
		interceptorsOption := grpc.ChainUnaryInterceptor(deps.UnaryInterceptors...)
		builder = builder.AddGRPCServerOptions(interceptorsOption)
	}
	builder = deps.buildExternalAPI(builder)
	return deps.buildInternalAPI(builder)
}

func (deps httpServerDeps) buildExternalAPI(builder serverInt.GRPCWebServiceBuilder) serverInt.GRPCWebServiceBuilder {
	if len(deps.GRPCServerAPIs) > 0 {
		builder = builder.RegisterGRPCAPIs(deps.GRPCServerAPIs...) // register grpc APIs
	}
	// add GRPC Gateway on top and expose on external REST Port
	externalRESTPort := deps.Config.Get(mortar.ServerRESTExternalPort)
	if len(deps.GRPCGatewayGeneratedHandlers) > 0 && externalRESTPort.IsSet() {
		builder = builder.AddRESTServerConfiguration().
			ListenOn(fmt.Sprintf(":%d", externalRESTPort.Int())).
			AddGRPCGatewayOptions(deps.GRPCGatewayMuxOptions...).
			RegisterGRPCGatewayHandlers(deps.GRPCGatewayGeneratedHandlers...).
			BuildRESTPart()
	}
	return builder
}

func (deps httpServerDeps) buildInternalAPI(builder serverInt.GRPCWebServiceBuilder) serverInt.GRPCWebServiceBuilder {
	// Internal
	internalPort := deps.Config.Get(mortar.ServerRESTInternalPort)
	includeInternalREST := internalPort.IsSet() && (len(deps.InternalHttpHandlerFunctions) > 0 || len(deps.InternalHttpHandlers) > 0)
	if includeInternalREST {
		restBuilder := builder.
			RegisterGRPCAPIs(health.RegisterInternalHealthService). // add internal GRPC health endpoint
			AddRESTServerConfiguration().
			ListenOn(fmt.Sprintf(":%d", internalPort.Int()))
		for _, handlerPair := range deps.InternalHttpHandlers {
			restBuilder = restBuilder.AddHandler(handlerPair.Pattern, handlerPair.Handler)
		}
		for _, handlerFuncPair := range deps.InternalHttpHandlerFunctions {
			restBuilder = restBuilder.AddHandlerFunc(handlerFuncPair.Pattern, handlerFuncPair.HandlerFunc)
		}
		restBuilder = restBuilder.RegisterGRPCGatewayHandlers(health.RegisterInternalGRPCGatewayHandler) // Health
		builder = restBuilder.BuildRESTPart()
	}
	return builder
}
