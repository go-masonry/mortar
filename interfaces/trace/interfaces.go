package trace

import (
	"context"
	"github.com/opentracing/opentracing-go"
)

type OpenTracer interface {
	Connect(ctx context.Context) error
	Tracer() opentracing.Tracer
	Close(ctx context.Context) error
}
