package constructors

import (
	"context"
	"fmt"
	"github.com/go-masonry/mortar/interfaces/auth/jwt"
	"go.uber.org/fx"
	"google.golang.org/grpc/metadata"
	"strings"
)

const (
	authorizationHeader = "authorization"
)

type jwtTokenExtractorDeps struct {
	fx.In

	Builder jwt.ExtractorBuilder
}

func DefaultJWTTokenExtractor(deps jwtTokenExtractorDeps) jwt.TokenExtractor {
	return deps.Builder.SetContextExtractor(contextExtractorAuthWithBearer).Build()
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
			} else {
				return "", fmt.Errorf("%s header value [%s] is of a wrong format", authorizationHeader, headerValue)
			}
		} else {
			return "", fmt.Errorf("context missing %s header", authorizationHeader)
		}
	}
	return "", fmt.Errorf("context missing gRPC incomming key")
}
