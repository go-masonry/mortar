package grpcgateway

import (
	"net/textproto"
	"strings"

	"github.com/go-masonry/mortar/interfaces/cfg"
	confkeys "github.com/go-masonry/mortar/interfaces/cfg/keys"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
)

type grpcGatewayHeadersDeps struct {
	fx.In

	Config cfg.Config
}

// MapHTTPHeadersToClientMetadataMuxOption maps incoming HTTP Headers to gRPC Context by checking if they match a list of prefixes
func MapHTTPHeadersToClientMetadataMuxOption(deps grpcGatewayHeadersDeps) runtime.ServeMuxOption {
	return runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
		for _, prefix := range deps.Config.Get(confkeys.MapHTTPRequestHeadersToGRPCMetadata).StringSlice() {
			// `key` is already canonicalized by Grpc-Gateway before calling this code
			if strings.HasPrefix(key, textproto.CanonicalMIMEHeaderKey(prefix)) {
				return key, true
			}
		}
		return runtime.DefaultHeaderMatcher(key)
	})
}
