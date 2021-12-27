package server

import (
	"context"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

//go:generate mockgen -source=interfaces.go -destination=mock/mock.go

// WebServerType string enum
type WebServerType string

const (
	// GRPCServer type
	GRPCServer WebServerType = "GRPC"
	// RESTServer type
	RESTServer WebServerType = "REST"
)

// ListenInfo defines port info
type ListenInfo struct {
	Address string        `json:"address"`
	Port    int           `json:"port"`
	Type    WebServerType `json:"type"`
}

// WebService defines our web service functions
type WebService interface {
	Run(ctx context.Context) error
	Stop(ctx context.Context) error
	Ports() []ListenInfo
}

// GRPCServerAPI alias for gRPC API function registration
type GRPCServerAPI func(server *grpc.Server)

// GRPCWebServiceBuilder defines gRPC web service builder options
type GRPCWebServiceBuilder interface {
	ListenOn(addr string) GRPCWebServiceBuilder
	SetCustomGRPCServer(customServer *grpc.Server) GRPCWebServiceBuilder
	SetCustomListener(listener net.Listener) GRPCWebServiceBuilder
	RegisterGRPCAPIs(register ...GRPCServerAPI) GRPCWebServiceBuilder
	AddGRPCServerOptions(options ...grpc.ServerOption) GRPCWebServiceBuilder
	SetPanicHandler(handler func(interface{}) error) GRPCWebServiceBuilder
	SetLogger(logger func(ctx context.Context, format string, args ...interface{})) GRPCWebServiceBuilder
	AddRESTServerConfiguration() RESTBuilder
	Build() (WebService, error)
}

// GRPCGatewayGeneratedHandlers alias for gRPC-gateway endpoint registrations
type GRPCGatewayGeneratedHandlers func(mux *runtime.ServeMux, endpoint string) error

// GRPCGatewayInterceptor alias for gRPC-gateway interceptor
type GRPCGatewayInterceptor func(handler http.Handler) http.Handler

// RESTBuilder defines REST web service builder options
type RESTBuilder interface {
	ListenOn(addr string) RESTBuilder
	SetCustomServer(customServer *http.Server) RESTBuilder
	SetCustomListener(listener net.Listener) RESTBuilder
	AddHandler(pattern string, handler http.Handler) RESTBuilder
	AddHandlerFunc(pattern string, handlerFunc http.HandlerFunc) RESTBuilder
	SetCustomGRPCGatewayMux(mux *runtime.ServeMux) RESTBuilder
	RegisterGRPCGatewayHandlers(handlers ...GRPCGatewayGeneratedHandlers) RESTBuilder
	AddGRPCGatewayOptions(options ...runtime.ServeMuxOption) RESTBuilder
	AddGRPCGatewayInterceptors(interceptors ...GRPCGatewayInterceptor) RESTBuilder
	BuildRESTPart() GRPCWebServiceBuilder
}
