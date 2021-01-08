package client

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-masonry/mortar/interfaces/cfg"
	confkeys "github.com/go-masonry/mortar/interfaces/cfg/keys"
	"github.com/go-masonry/mortar/interfaces/http/client"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type copyHeadersDeps struct {
	fx.In

	Config cfg.Config
}

// CopyGRPCHeadersClientInterceptor copies filtered Headers found in the Incoming metadata.MD to the Outgoing one.
//
// This is useful if you want to propagate them to the next service when using grpc Client
//
// For Example: "authorization" header containing user token
func CopyGRPCHeadersClientInterceptor(deps copyHeadersDeps) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			headerPrefixes := deps.Config.Get(confkeys.ForwardIncomingGRPCMetadataHeadersList).StringSlice()
			for _, headerPrefix := range headerPrefixes {
				for k, vs := range md {
					if strings.HasPrefix(strings.ToLower(k), headerPrefix) {
						for _, v := range vs {
							ctx = metadata.AppendToOutgoingContext(ctx, k, v)
						}
					}
				}
			}
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// CopyGRPCHeadersHTTPClientInterceptor copies filtered Headers found in the Incoming metadata.MD to the Outgoing Request Headers.
//
// This is useful if you want to propagate them to the next service when using `http.Client`
//
// For Example: "authorization" header containing user token
func CopyGRPCHeadersHTTPClientInterceptor(deps copyHeadersDeps) client.HTTPClientInterceptor {
	return func(req *http.Request, handler client.HTTPHandler) (resp *http.Response, err error) {
		if md, ok := metadata.FromIncomingContext(req.Context()); ok {
			headerPrefixes := deps.Config.Get(confkeys.ForwardIncomingGRPCMetadataHeadersList).StringSlice()
			for _, headerPrefix := range headerPrefixes {
				for k, vs := range md {
					if strings.HasPrefix(strings.ToLower(k), headerPrefix) {
						// Remember the key will be canonicalized by `http.CanonicalHeaderKey`
						for _, v := range vs {
							req.Header.Add(k, v)
						}
					}
				}
			}
		}
		return handler(req)
	}
}
