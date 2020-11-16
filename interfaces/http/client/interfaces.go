package client

import (
	"context"
	"net/http"

	"google.golang.org/grpc"
)

//go:generate mockgen -source=interfaces.go -destination=mock/mock.go

//********************************************************************************
// http.Client
//********************************************************************************

// HTTPHandler is just an alias to http.RoundTriper.RoundTrip function
type HTTPHandler func(*http.Request) (*http.Response, error)

// HTTPClientInterceptor is a user defined function that can alter a request before it's sent
// and/or alter a response before it's returned to the caller
type HTTPClientInterceptor func(*http.Request, HTTPHandler) (*http.Response, error)

// HTTPClientBuilder is a builder interface to build http.Client with interceptors
type HTTPClientBuilder interface {
	AddInterceptors(...HTTPClientInterceptor) HTTPClientBuilder
	WithPreconfiguredClient(*http.Client) HTTPClientBuilder
	Build() *http.Client
}

// NewHTTPClientBuilder REST HTTP builder
//
// Useful when you want to create several *http.Client with different options
type NewHTTPClientBuilder func() HTTPClientBuilder

//********************************************************************************
// grpc.Client
//********************************************************************************

// GRPCClientConnectionWrapper is a convenience wrapper to support predefined dial Options
// provided by GRPCClientConnectionBuilder
type GRPCClientConnectionWrapper interface {
	// Context can be nil
	Dial(ctx context.Context, target string, extraOptions ...grpc.DialOption) (grpc.ClientConnInterface, error)
}

// GRPCClientConnectionBuilder is a convenience builder to gather []grpc.DialOption
type GRPCClientConnectionBuilder interface {
	AddOptions(opts ...grpc.DialOption) GRPCClientConnectionBuilder
	Build() GRPCClientConnectionWrapper
}
