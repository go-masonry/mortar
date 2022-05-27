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
	"github.com/go-masonry/mortar/interfaces/monitor"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

// PanicHandlerCounter is the metric name to count all recovered panics
const PanicHandlerCounter = "panic_handler_total"

// Group order is not guaranteed, if it's important then add them manually
const (
	// FxGroupBuilderCallbacks defines group name
	FxGroupBuilderCallbacks = "builderCallbacks"
	// FxGroupGRPCServerAPIs defines group name
	FxGroupGRPCServerAPIs = "grpcServerAPIs"
	// FxGroupGRPCGatewayGeneratedHandlers defines group name
	FxGroupGRPCGatewayGeneratedHandlers = "grpcGatewayGeneratedHandlers"
	// FxGroupGRPCGatewayMuxOptions defines group name
	FxGroupGRPCGatewayMuxOptions = "grpcGatewayMuxOptions"
	// FxGroupExternalBuilderCallbacks defines group name
	FxGroupExternalBuilderCallbacks = "externalBuilderCallbacks"
	// FxGroupExternalHTTPHandlers defines group name
	FxGroupExternalHTTPHandlers = "externalHttpHandlers"
	// FxGroupExternalHTTPHandlerFunctions defines group name
	FxGroupExternalHTTPHandlerFunctions = "externalHttpHandlerFunctions"
	// FxGroupExternalHTTPInterceptors defines group name
	FxGroupExternalHTTPInterceptors = "externalHttpInterceptors"
	// FxGroupUnaryServerInterceptors defines group name
	FxGroupUnaryServerInterceptors = "unaryServerInterceptors"
	// FxGroupInternalBuilderCallbacks defines group name
	FxGroupInternalBuilderCallbacks = "internalBuilderCallbacks"
	// FxGroupInternalHTTPHandlers defines group name
	FxGroupInternalHTTPHandlers = "internalHttpHandlers"
	// FxGroupInternalHTTPHandlerFunctions defines group name
	FxGroupInternalHTTPHandlerFunctions = "internalHttpHandlerFunctions"
	// FxGroupInternalHTTPInterceptors defines group name
	FxGroupInternalHTTPInterceptors = "internalHttpInterceptors"
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

// BuilderCallback is the callback to modify the builder or do something else
type BuilderCallback func(serverInt.GRPCWebServiceBuilder) serverInt.GRPCWebServiceBuilder

// RESTBuilderCallback is the callback to modify the builder or do something else
type RESTBuilderCallback func(serverInt.RESTBuilder) serverInt.RESTBuilder

type httpServerDeps struct {
	fx.In

	Config  cfg.Config
	Logger  log.Logger
	Metrics monitor.Metrics `optional:"true"`

	// Builder
	BuilderCallbacks []BuilderCallback `group:"builderCallbacks"`

	// GRPC
	GRPCServerAPIs    []serverInt.GRPCServerAPI     `group:"grpcServerAPIs"`
	UnaryInterceptors []grpc.UnaryServerInterceptor `group:"unaryServerInterceptors"`

	// External REST
	GRPCGatewayGeneratedHandlers []serverInt.GRPCGatewayGeneratedHandlers `group:"grpcGatewayGeneratedHandlers"`
	GRPCGatewayMuxOptions        []runtime.ServeMuxOption                 `group:"grpcGatewayMuxOptions"`
	ExternalBuilderCallbacks     []RESTBuilderCallback                    `group:"externalBuilderCallbacks"`
	ExternalHTTPHandlers         []HTTPHandlerPatternPair                 `group:"externalHttpHandlers"`
	ExternalHTTPHandlerFunctions []HTTPHandlerFuncPatternPair             `group:"externalHttpHandlerFunctions"`
	ExternalHTTPInterceptors     []serverInt.GRPCGatewayInterceptor       `group:"externalHttpInterceptors"`

	// Internal REST
	InternalBuilderCallbacks     []RESTBuilderCallback              `group:"internalBuilderCallbacks"`
	InternalHTTPHandlers         []HTTPHandlerPatternPair           `group:"internalHttpHandlers"`
	InternalHTTPHandlerFunctions []HTTPHandlerFuncPatternPair       `group:"internalHttpHandlerFunctions"`
	InternalHTTPInterceptors     []serverInt.GRPCGatewayInterceptor `group:"internalHttpInterceptors"`
}

// HTTPServerBuilder true to it's name, it is partially initialized builder.
//
// It uses some default assumptions and configurations, which are mostly good.
// However, if you need to customize your configuration it's better to build yours from scratch
//
func HTTPServerBuilder(deps httpServerDeps) serverInt.GRPCWebServiceBuilder {
	builder := server.Builder().SetPanicHandler(deps.panicHandler).SetLogger(deps.Logger.Debug)
	host := deps.Config.Get(confkeys.Host).String()
	// GRPC port
	if grpcPort := deps.Config.Get(confkeys.ExternalGRPCPort); grpcPort.IsSet() {
		builder = builder.ListenOn(fmt.Sprintf("%s:%d", host, grpcPort.Int()))
	}
	// GRPC server interceptors
	if len(deps.UnaryInterceptors) > 0 {
		interceptorsOption := grpc.ChainUnaryInterceptor(deps.UnaryInterceptors...)
		builder = builder.AddGRPCServerOptions(interceptorsOption)
	}
	builder = deps.buildExternalAPI(builder)
	builder = deps.buildInternalAPI(builder)
	for _, callback := range deps.BuilderCallbacks {
		builder = callback(builder)
	}

	return builder
}

func (deps httpServerDeps) buildExternalAPI(builder serverInt.GRPCWebServiceBuilder) serverInt.GRPCWebServiceBuilder {
	if len(deps.GRPCServerAPIs) > 0 {
		builder = builder.RegisterGRPCAPIs(deps.GRPCServerAPIs...) // register grpc APIs
	}
	// add GRPC Gateway on top and expose on external REST Port
	host := deps.Config.Get(confkeys.Host).String()
	externalRESTPort := deps.Config.Get(confkeys.ExternalRESTPort)
	if externalRESTPort.IsSet() && (len(deps.ExternalHTTPHandlerFunctions) > 0 || len(deps.ExternalHTTPHandlers) > 0 || len(deps.GRPCGatewayGeneratedHandlers) > 0) {
		restBuilder := builder.AddRESTServerConfiguration().
			ListenOn(fmt.Sprintf("%s:%d", host, externalRESTPort.Int()))

		for _, callback := range deps.ExternalBuilderCallbacks {
			restBuilder = callback(restBuilder)
		}
		for _, handlerPair := range deps.ExternalHTTPHandlers {
			restBuilder = restBuilder.AddHandler(handlerPair.Pattern, handlerPair.Handler)
		}
		for _, handlerFuncPair := range deps.ExternalHTTPHandlerFunctions {
			restBuilder = restBuilder.AddHandlerFunc(handlerFuncPair.Pattern, handlerFuncPair.HandlerFunc)
		}
		if len(deps.ExternalHTTPInterceptors) > 0 {
			restBuilder = restBuilder.AddGRPCGatewayInterceptors(deps.ExternalHTTPInterceptors...)
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
	host := deps.Config.Get(confkeys.Host).String()
	internalPort := deps.Config.Get(confkeys.InternalRESTPort)
	includeInternalREST := internalPort.IsSet() && (len(deps.InternalHTTPHandlerFunctions) > 0 || len(deps.InternalHTTPHandlers) > 0)
	if includeInternalREST {
		restBuilder := builder.
			AddRESTServerConfiguration().
			ListenOn(fmt.Sprintf("%s:%d", host, internalPort.Int()))

		for _, callback := range deps.InternalBuilderCallbacks {
			restBuilder = callback(restBuilder)
		}
		for _, handlerPair := range deps.InternalHTTPHandlers {
			restBuilder = restBuilder.AddHandler(handlerPair.Pattern, handlerPair.Handler)
		}
		for _, handlerFuncPair := range deps.InternalHTTPHandlerFunctions {
			restBuilder = restBuilder.AddHandlerFunc(handlerFuncPair.Pattern, handlerFuncPair.HandlerFunc)
		}
		if len(deps.InternalHTTPInterceptors) > 0 {
			restBuilder = restBuilder.AddGRPCGatewayInterceptors(deps.InternalHTTPInterceptors...)
		}
		restBuilder = restBuilder.RegisterGRPCGatewayHandlers(health.RegisterInternalGRPCGatewayHandler) // Health
		builder = restBuilder.BuildRESTPart()
	}
	return builder
}

func (deps httpServerDeps) panicHandler(r interface{}) error {
	if deps.Metrics != nil {
		deps.Metrics.Counter(PanicHandlerCounter, "Count gRPC panic recoveries").Inc()
	}
	switch t := r.(type) {
	case string, fmt.Stringer:
		return fmt.Errorf("panic handled, %s", t)
	case error:
		return fmt.Errorf("panic handled, %w", t)
	default:
		return fmt.Errorf("panic handled, %v", t)
	}
}
