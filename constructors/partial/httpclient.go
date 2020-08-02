package partial

import (
	"github.com/go-masonry/mortar/http/client"
	clientInt "github.com/go-masonry/mortar/interfaces/http/client"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

const (
	FxGroupRESTClientInterceptors      = "restClientInterceptors"
	FxGroupGRPCUnaryClientInterceptors = "grpcUnaryClientInterceptors"
)

// REST HTTP
type HttpClientPartialBuilder func() clientInt.HttpClientBuilder

type httpClientBuilderDeps struct {
	fx.In

	Interceptors []clientInt.HttpClientInterceptor `group:"restClientInterceptors"`
}

// HttpClientBuilder creates an injectable http.Client builder that can be predefined with Interceptors
//
// This function returns a closure that will always create a new builder. That way every usage can add additional
// interceptors without influencing others
func HttpClientBuilder(deps httpClientBuilderDeps) HttpClientPartialBuilder {
	return func() clientInt.HttpClientBuilder {
		return client.HTTPClientBuilder().AddInterceptors(deps.Interceptors...)
	}
}

// GRPC
type grpcClientConnectionBuilderDeps struct {
	fx.In

	Interceptors []grpc.UnaryClientInterceptor `group:"grpcUnaryClientInterceptors"`
}

// GRPCClientConnectionBuilder creates an injectable grpc.ClientConn that can be predefined with Interceptors
// or/and additional options later
func GRPCClientConnectionBuilder(deps grpcClientConnectionBuilderDeps) clientInt.GRPCClientConnectionBuilder {
	interceptors := grpc.WithChainUnaryInterceptor(deps.Interceptors...)
	return client.GRPCClientConnBuilder().AddOptions(interceptors)
}
