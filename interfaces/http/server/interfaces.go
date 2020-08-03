package server

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

type WebServerType string

const (
	GRPCServer WebServerType = "GRPC"
	RESTServer WebServerType = "REST"
)

type ListenInfo struct {
	Address string        `json:"address"`
	Port    int           `json:"port"`
	Type    WebServerType `json:"type"`
}
type WebService interface {
	Run(ctx context.Context) error
	Stop(ctx context.Context) error
	Ports() []ListenInfo
}

type GRPCServerAPI func(server *grpc.Server)

type GRPCWebServiceBuilder interface {
	ListenOn(addr string) GRPCWebServiceBuilder
	SetCustomGRPCServer(server *grpc.Server) GRPCWebServiceBuilder
	SetCustomListener(listener net.Listener) GRPCWebServiceBuilder
	RegisterGRPCAPIs(register ...GRPCServerAPI) GRPCWebServiceBuilder
	AddGRPCServerOptions(options ...grpc.ServerOption) GRPCWebServiceBuilder
	SetPanicHandler(handler func(interface{}) error) GRPCWebServiceBuilder
	SetLogger(logger func(ctx context.Context, format string, args ...interface{})) GRPCWebServiceBuilder
	AddRESTServerConfiguration() RESTBuilder
	Build() (WebService, error)
}

type GRPCGatewayGeneratedHandlers func(mux *runtime.ServeMux, endpoint string) error

type RESTBuilder interface {
	ListenOn(addr string) RESTBuilder
	SetCustomServer(server *http.Server) RESTBuilder
	SetCustomListener(listener net.Listener) RESTBuilder
	AddHandler(pattern string, handler http.Handler) RESTBuilder
	AddHandlerFunc(pattern string, handlerFunc http.HandlerFunc) RESTBuilder
	SetCustomGRPCGatewayMux(mux *runtime.ServeMux) RESTBuilder
	RegisterGRPCGatewayHandlers(handlers ...GRPCGatewayGeneratedHandlers) RESTBuilder
	AddGRPCGatewayOptions(options ...runtime.ServeMuxOption) RESTBuilder
	BuildRESTPart() GRPCWebServiceBuilder
}
