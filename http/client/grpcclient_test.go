package client

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestGRPCClientConnBuilder(t *testing.T) {
	wrapper := GRPCClientConnBuilder().AddOptions(grpc.WithInsecure()).Build()
	if impl, ok := wrapper.(*grpcClientConnImpl); assert.True(t, ok) {
		assert.Len(t, impl.options.options, 1)
	}
}

func TestGRPCClientConnWrapperNoContext(t *testing.T) {
	wrapper := GRPCClientConnBuilder().AddOptions(grpc.WithInsecure()).Build()
	_, err := wrapper.Dial(nil, ":6666", grpc.WithBlock(), grpc.FailOnNonTempDialError(true))
	assert.Error(t, err)
}
func TestGRPCClientConnWrapperWithContext(t *testing.T) {
	wrapper := GRPCClientConnBuilder().AddOptions(grpc.WithInsecure()).Build()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	_, err := wrapper.Dial(ctx, ":6666", grpc.WithBlock())
	assert.Error(t, err)
}
