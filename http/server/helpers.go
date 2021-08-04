package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"

	"google.golang.org/grpc"
)

// Interfaces
type mux interface {
	Serve(listener net.Listener) error
}

type grpcServerStopper interface {
	GracefulStop()
	Stop()
}

type restServerShutdown interface {
	Shutdown(ctx context.Context) error
	Close() error
}

func defaultPanicHandler(r interface{}) error {
	switch t := r.(type) {
	case string, fmt.Stringer:
		return fmt.Errorf("panic handled, %s", t)
	case error:
		return fmt.Errorf("panic handled, %w", t)
	default:
		return fmt.Errorf("panic handled, %v", t)
	}
}

func panicHandlerUnaryInterceptor(panicHandler func(interface{}) error) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			r := recover()
			if r != nil {
				debug.PrintStack()
				err = panicHandler(r)
			}
		}()
		resp, err = handler(ctx, req)
		return
	}
}

func panicHandlerStreamInterceptor(panicHandler func(interface{}) error) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			r := recover()
			if r != nil {
				debug.PrintStack()
				err = panicHandler(r)
			}
		}()
		err = handler(srv, ss)
		return
	}
}

func createListener(network, addr string) (net.Listener, error) {
	if len(addr) == 0 {
		addr = "localhost:0"
	}
	return net.Listen(network, addr)
}

type muxHandler interface {
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

func extractPort(addr string) int {
	if index := strings.LastIndex(addr, ":"); index >= 0 && index+1 < len(addr) {
		if port, err := strconv.Atoi(addr[index+1:]); err == nil {
			return port
		}
	}
	return 0
}
