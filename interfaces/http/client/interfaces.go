package client

import (
	"context"
	"google.golang.org/grpc"
	"net/http"
)

//go:generate mockgen -source=interfaces.go -destination=mock/mock.go

//********************************************************************************
// http.Client
//********************************************************************************

// HttpHandler is just an alias to http.RoundTriper.RoundTrip function
type HttpHandler func(*http.Request) (*http.Response, error)

// HttpClientInterceptor is a user defined function that can alter a request before it's sent
// and/or alter a response before it's returned to the caller
type HttpClientInterceptor func(*http.Request, HttpHandler) (*http.Response, error)

// HTTPClientBuilder is a builder interface to build http.Client with interceptors
type HttpClientBuilder interface {
	AddInterceptors(...HttpClientInterceptor) HttpClientBuilder
	WithPreconfiguredClient(*http.Client) HttpClientBuilder
	Build() *http.Client
}

//********************************************************************************
// grpc.Client
//********************************************************************************

// GRPCClientConnectionWrapper is a convenience wrapper to support predefined dial Options
// provided by GRPCClientConnectionBuilder
type GRPCClientConnectionWrapper interface {
	// Context can be nil
	Dial(ctx context.Context, target string, extraOptions ...grpc.DialOption) (*grpc.ClientConn, error)
}

// GRPCClientConnectionBuilder is a convenience builder to gather []grpc.DialOption
type GRPCClientConnectionBuilder interface {
	AddOptions(opts ...grpc.DialOption) GRPCClientConnectionBuilder
	Build() GRPCClientConnectionWrapper
}
