package providers

import (
	"github.com/go-masonry/mortar/constructors"
	"github.com/go-masonry/mortar/constructors/partial"
	"github.com/go-masonry/mortar/middleware/context"
	"github.com/go-masonry/mortar/middleware/interceptors/server"
	"go.uber.org/fx"
)

// LoggerFxOption adds Default Logger to the graph
func LoggerFxOption() fx.Option {
	return fx.Provide(constructors.DefaultLogger)
}

// LoggerGRPCIncomingContextExtractorFxOption adds Logger Context Extractor using values within incoming grpc metadata.MD
//
// This one will be included during Logger build
func LoggerGRPCIncomingContextExtractorFxOption() fx.Option {
	return fx.Provide(fx.Annotated{
		Group:  constructors.FxGroupLoggerContextExtractors,
		Target: context.LoggerGRPCIncomingContextExtractor,
	})
}

// LoggerGRPCInterceptorFxOption adds Server Unary Interceptor that will log Request and Response if needed
func LoggerGRPCInterceptorFxOption() fx.Option {
	return fx.Provide(fx.Annotated{
		Group:  partial.FxGroupUnaryServerInterceptors,
		Target: server.LoggerGRPCInterceptor,
	})
}
