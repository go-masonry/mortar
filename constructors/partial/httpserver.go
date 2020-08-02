package partial

import (
	"fmt"
	"github.com/go-masonry/mortar/http/server"
	"github.com/go-masonry/mortar/interfaces/cfg"
	serverInt "github.com/go-masonry/mortar/interfaces/http/server"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/go-masonry/mortar/mortar"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"net/http"
)

const (
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
	// Order of interceptors is not guaranteed, if it's important then add them manually
	UnaryInterceptors            []grpc.UnaryServerInterceptor `group:"unaryServerInterceptors"`
	InternalHttpHandlers         []HttpHandlerPatternPair      `group:"internalHttpHandlers"`
	InternalHttpHandlerFunctions []HttpHandlerFuncPatternPair  `group:"internalHttpHandlerFunctions"`
}

// PartialHttpServerBuilder true to it's name, it is partially initialized builder, like other builders
// you can influence it's configuration
//
// 	- GRPC Port: we will look for a key mortar.ServerGRPCPort within the configuration map
// 	- UnaryServerInterceptors: Since we are using uber.Fx for DI we can expect any number of unary server interceptors
//		All unary server interceptors must be grouped under a fx.Group named: 'unaryServerInterceptors'
// 	- InternalHttpHandlers: Since we are using uber.Fx for DI we can expect any number of http handlers
//		All Internal HTTP Handlers must be grouped under a fx.Group named: 'internalHttpHandlers'
//	- InternalHttpHandlerFunctions: Since we are using uber.Fx for DI we can expect any number of http handler functions
//		All Internal HTTP Handler Functions must be grouped under a fx.Group named: 'internalHttpHandlerFunctions'
func HttpServerBuilder(deps httpServerDeps) serverInt.GRPCWebServiceBuilder {
	builder := server.Builder().SetLogger(deps.Logger.Info)
	// GRPC
	if grpcPort := deps.Config.Get(mortar.ServerGRPCPort); grpcPort.IsSet() {
		builder = builder.ListenOn(fmt.Sprintf(":%d", grpcPort.Int()))
	}
	if len(deps.UnaryInterceptors) > 0 {
		interceptorsOption := grpc.ChainUnaryInterceptor(deps.UnaryInterceptors...)
		builder = builder.AddGRPCServerOptions(interceptorsOption)
	}
	return deps.buildInternalREST(builder)
}

func (deps httpServerDeps) buildInternalREST(builder serverInt.GRPCWebServiceBuilder) serverInt.GRPCWebServiceBuilder {
	// Internal REST
	internalPort := deps.Config.Get(mortar.ServerRESTInternalPort)
	includeInternalREST := internalPort.IsSet() && (len(deps.InternalHttpHandlerFunctions) > 0 || len(deps.InternalHttpHandlers) > 0)
	if includeInternalREST {
		restBuilder := builder.AddRESTServerConfiguration().
			ListenOn(fmt.Sprintf(":%d", internalPort.Int()))
		for _, handlerPair := range deps.InternalHttpHandlers {
			restBuilder = restBuilder.AddHandler(handlerPair.Pattern, handlerPair.Handler)
		}
		for _, handlerFuncPair := range deps.InternalHttpHandlerFunctions {
			restBuilder = restBuilder.AddHandlerFunc(handlerFuncPair.Pattern, handlerFuncPair.HandlerFunc)
		}
		builder = restBuilder.BuildRESTPart()
	}
	return builder
}
