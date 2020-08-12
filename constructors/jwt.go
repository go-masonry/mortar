package constructors

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-masonry/mortar/auth/jwt"
	jwtInt "github.com/go-masonry/mortar/interfaces/auth/jwt"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
)

// DefaultJWTTokenExtractor simple TokenExtractor
func DefaultJWTTokenExtractor() jwtInt.TokenExtractor {
	return jwt.Builder().SetContextExtractor(contextExtractorAuthWithBearer).Build()
}

// Handles use cases where 'authorization' header value is
// 		bearer <token>
//		basic <token>
func contextExtractorAuthWithBearer(ctx context.Context) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if headerValue := strings.Join(md.Get(authorizationHeader), " "); len(headerValue) > 0 {
			rawTokenWithBearer := strings.Split(headerValue, " ")
			if len(rawTokenWithBearer) == 2 {
				return rawTokenWithBearer[1], nil
			}
			return "", fmt.Errorf("%s header value [%s] is of a wrong format", authorizationHeader, headerValue)
		}
		return "", fmt.Errorf("context missing %s header", authorizationHeader)
	}
	return "", fmt.Errorf("context missing gRPC incomming key")
}
