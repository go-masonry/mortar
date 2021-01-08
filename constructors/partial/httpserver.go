package partial

import (
	"fmt"
	"net/http"

	"github.com/go-masonry/mortar/http/server"
	"github.com/go-masonry/mortar/http/server/health"
	"github.com/go-masonry/mortar/interfaces/cfg"
	confkeys "github.com/go-masonry/mortar/interfaces/cfg/keys"
	serverInt "github.com/go-masonry/mortar/interfaces/http/server"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

// Group order is not guaranteed, if it's important then add them manually
const (
	// FxGroupGRPCServerAPIs defines group name
	FxGroupGRPCServerAPIs = "grpcServerAPIs"
	// FxGroupGRPCGatewayGeneratedHandlers defines group name
	FxGroupGRPCGatewayGeneratedHandlers = "grpcGatewayGeneratedHandlers"
	// FxGroupGRPCGatewayMuxOptions defines group name
	FxGroupGRPCGatewayMuxOptions = "grpcGatewayMuxOptions"
	// FxGroupExternalHTTPHandlers defines group name
	FxGroupExternalHTTPHandlers = "externalHttpHandlers"
	// FxGroupExternalHTTPHandlerFunctions defines group name
	FxGroupExternalHTTPHandlerFunctions = "externalHttpHandlerFunctions"
	// FxGroupUnaryServerInterceptors defines group name
	FxGroupUnaryServerInterceptors = "unaryServerInterceptors"
	// FxGroupInternalHTTPHandlers defines group name
	FxGroupInternalHTTPHandlers = "internalHttpHandlers"
	// FxGroupInternalHTTPHandlerFunctions defines group name
	FxGroupInternalHTTPHandlerFunctions = "internalHttpHandlerFunctions"
)

// HTTPHandlerPatternPair defines pattern -> handler pair
type HTTPHandlerPatternPair struct {
	Pattern string
	Handler http.Handler
}

// HTTPHandlerFuncPatternPair defines patter -> handler func pair
type HTTPHandlerFuncPatternPair struct {
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
	ExternalHTTPHandlers         []HTTPHandlerPatternPair                 `group:"externalHttpHandlers"`
	ExternalHTTPHandlerFunctions []HTTPHandlerFuncPatternPair             `group:"externalHttpHandlerFunctions"`
	// Internal REST
	InternalHTTPHandlers         []HTTPHandlerPatternPair     `group:"internalHttpHandlers"`
	InternalHTTPHandlerFunctions []HTTPHandlerFuncPatternPair `group:"internalHttpHandlerFunctions"`
}

// HTTPServerBuilder true to it's name, it is partially initialized builder.
//
// It uses some default assumptions and configurations, which are mostly good.
// However, if you need to customize your configuration it's better to build yours from scratch
//
func HTTPServerBuilder(deps httpServerDeps) serverInt.GRPCWebServiceBuilder {
	builder := server.Builder().SetLogger(deps.Logger.Debug)
	// GRPC port
	if grpcPort := deps.Config.Get(confkeys.ExternalGRPCPort); grpcPort.IsSet() {
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
	externalRESTPort := deps.Config.Get(confkeys.ExternalRESTPort)
	if externalRESTPort.IsSet() && (len(deps.ExternalHTTPHandlerFunctions) > 0 || len(deps.ExternalHTTPHandlers) > 0 || len(deps.GRPCGatewayGeneratedHandlers) > 0) {
		restBuilder := builder.AddRESTServerConfiguration().
			ListenOn(fmt.Sprintf(":%d", externalRESTPort.Int()))

		for _, handlerPair := range deps.ExternalHTTPHandlers {
			restBuilder = restBuilder.AddHandler(handlerPair.Pattern, handlerPair.Handler)
		}
		for _, handlerFuncPair := range deps.ExternalHTTPHandlerFunctions {
			restBuilder = restBuilder.AddHandlerFunc(handlerFuncPair.Pattern, handlerFuncPair.HandlerFunc)
		}
		if len(deps.GRPCGatewayGeneratedHandlers) > 0 {
			restBuilder = restBuilder.AddGRPCGatewayOptions(deps.GRPCGatewayMuxOptions...).
				RegisterGRPCGatewayHandlers(deps.GRPCGatewayGeneratedHandlers...)
		}
		builder = restBuilder.BuildRESTPart()

	}
	return builder
}

func (deps httpServerDeps) buildInternalAPI(builder serverInt.GRPCWebServiceBuilder) serverInt.GRPCWebServiceBuilder {
	builder = builder.RegisterGRPCAPIs(health.RegisterInternalHealthService) // add internal GRPC health endpoint
	// Internal
	internalPort := deps.Config.Get(confkeys.InternalRESTPort)
	includeInternalREST := internalPort.IsSet() && (len(deps.InternalHTTPHandlerFunctions) > 0 || len(deps.InternalHTTPHandlers) > 0)
	if includeInternalREST {
		restBuilder := builder.
			AddRESTServerConfiguration().
			ListenOn(fmt.Sprintf(":%d", internalPort.Int()))
		for _, handlerPair := range deps.InternalHTTPHandlers {
			restBuilder = restBuilder.AddHandler(handlerPair.Pattern, handlerPair.Handler)
		}
		for _, handlerFuncPair := range deps.InternalHTTPHandlerFunctions {
			restBuilder = restBuilder.AddHandlerFunc(handlerFuncPair.Pattern, handlerFuncPair.HandlerFunc)
		}
		restBuilder = restBuilder.RegisterGRPCGatewayHandlers(health.RegisterInternalGRPCGatewayHandler) // Health
		builder = restBuilder.BuildRESTPart()
	}
	return builder
}
