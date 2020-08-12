package trace

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

// If you need a mocked Tracer use one provided by the opentracing library
//	"github.com/opentracing/opentracing-go/mocktracer"

//go:generate mockgen -source=interfaces.go -destination=mock/mock.go

// OpenTracer defines service tracer with Connect/Close options on demand and not on creation
type OpenTracer interface {
	Connect(ctx context.Context) error
	Tracer() opentracing.Tracer
	Close(ctx context.Context) error
}
