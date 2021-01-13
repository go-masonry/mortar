package constructors_test

import (
	"context"
	"testing"

	"github.com/go-masonry/mortar/constructors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestDefaultJWTTokenExtractor(t *testing.T) {
	extractor := constructors.DefaultJWTTokenExtractor()
	ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{
		"authorization": []string{"bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE0MjM5MDIyfQ.0qQhlRLwbuLtNboNltWixpEh8vyP-DP-uUJoKO3m388"},
	})
	token, err := extractor.FromContext(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE0MjM5MDIyfQ.0qQhlRLwbuLtNboNltWixpEh8vyP-DP-uUJoKO3m388", token.Raw())
}

func TestDefaultJWTTokenExtractorWithGRPCGatewayPrefix(t *testing.T) {
	extractor := constructors.DefaultJWTTokenExtractor()
	ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{
		"grpcgateway-authorization": []string{"bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE0MjM5MDIyfQ.0qQhlRLwbuLtNboNltWixpEh8vyP-DP-uUJoKO3m388"},
	})
	token, err := extractor.FromContext(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE0MjM5MDIyfQ.0qQhlRLwbuLtNboNltWixpEh8vyP-DP-uUJoKO3m388", token.Raw())
}
