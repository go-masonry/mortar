package constructors

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-masonry/mortar/auth/jwt"
	jwtInt "github.com/go-masonry/mortar/interfaces/auth/jwt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader                = "authorization"
	grpcGatewayAuthorizationWithPrefix = runtime.MetadataPrefix + "authorization"
)

// DefaultJWTTokenExtractor simple TokenExtractor
func DefaultJWTTokenExtractor() jwtInt.TokenExtractor {
	return jwt.Builder().SetContextExtractor(contextExtractorAuthWithBearer).Build()
}

// Handles use cases where 'authorization' header value is
//
//	bearer <token>
//	basic <token>
func contextExtractorAuthWithBearer(ctx context.Context) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		var headerValue string = strings.Join(md.Get(authorizationHeader), " ")
		if !(len(headerValue) > 0) {
			headerValue = strings.Join(md.Get(grpcGatewayAuthorizationWithPrefix), " ")
		}
		if len(headerValue) > 0 {
			rawTokenWithBearer := strings.Split(headerValue, " ")
			if len(rawTokenWithBearer) == 2 {
				return rawTokenWithBearer[1], nil
			}
			return "", fmt.Errorf("%s/%s header value [%s] is of a wrong format", authorizationHeader, grpcGatewayAuthorizationWithPrefix, headerValue)
		}
		return "", fmt.Errorf("context missing %s/%s header", authorizationHeader, grpcGatewayAuthorizationWithPrefix)
	}
	return "", fmt.Errorf("context missing gRPC incoming key")
}
