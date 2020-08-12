package client

import (
	"container/list"
	"context"

	"github.com/go-masonry/mortar/interfaces/http/client"
	"google.golang.org/grpc"
)

type grpcClientConnOptions struct {
	options []grpc.DialOption
}

type grpcClientConnBuilder struct {
	ll *list.List
}

// GRPCClientConnBuilder creates a fresh gRPC connection for client builder
func GRPCClientConnBuilder() client.GRPCClientConnectionBuilder {
	return &grpcClientConnBuilder{
		ll: list.New(),
	}
}

func (g *grpcClientConnBuilder) AddOptions(opts ...grpc.DialOption) client.GRPCClientConnectionBuilder {
	g.ll.PushBack(func(cfg *grpcClientConnOptions) {
		cfg.options = append(cfg.options, opts...)
	})
	return g
}

func (g *grpcClientConnBuilder) Build() client.GRPCClientConnectionWrapper {
	var cfg = new(grpcClientConnOptions)
	for e := g.ll.Front(); e != nil; e = e.Next() {
		f := e.Value.(func(connOptions *grpcClientConnOptions))
		f(cfg)
	}
	return &grpcClientConnImpl{
		options: cfg,
	}
}

type grpcClientConnImpl struct {
	options *grpcClientConnOptions
}

func (g *grpcClientConnImpl) Dial(ctx context.Context, target string, extraOptions ...grpc.DialOption) (*grpc.ClientConn, error) {
	var allOptions = append(g.options.options, extraOptions...)
	if ctx == nil {
		ctx = context.Background()
	}
	return grpc.DialContext(ctx, target, allOptions...)
}
