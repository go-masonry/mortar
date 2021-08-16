package providers

import (
	"github.com/go-masonry/mortar/constructors"
	"github.com/go-masonry/mortar/logger"
	"github.com/go-masonry/mortar/middleware/context"
	"github.com/go-masonry/mortar/middleware/interceptors/server"
	"github.com/go-masonry/mortar/providers/groups"
	"go.uber.org/fx"
)

// LoggerFxOption adds Default Logger to the graph
func LoggerFxOption() fx.Option {
	return fx.Provide(constructors.DefaultLogger)
}

// DefaultLogger is a constructor that creates a default log.Logger based on provided log.Builder
//
// Consider using LoggerFxOption if you only want to invoke it.
var DefaultLogger = constructors.DefaultLogger

// FxEventLoggerOption add new Fx Event option to output fx events using structured logger
func FxEventLoggerOption() fx.Option {
	return fx.WithLogger(logger.CreateFxEventLogger)
}

// CreateFxEventLogger is a constuctor that creates a custom fxevent.Logger
//
// Consider using FxEventLoggerOption if you only want to invoke it.
var CreateFxEventLogger = logger.CreateFxEventLogger

// LoggerGRPCIncomingContextExtractorFxOption adds Logger Context Extractor using values within incoming grpc metadata.MD
//
// This one will be included during Logger build
func LoggerGRPCIncomingContextExtractorFxOption() fx.Option {
	return fx.Provide(fx.Annotated{
		Group:  groups.LoggerContextExtractors,
		Target: context.LoggerGRPCIncomingContextExtractor,
	})
}

// LoggerGRPCIncomingContextExtractor is a constructor that creates log.ContextExtractor.
// This Extractor will then extract selected key:value pairs from the context when writing log
//
//	*Note* normally this dependency is part of a group. If you want to create it as a standalone
// 	dependency, remember that there can be only one of this kind in the graph.
//
// Consider using LoggerGRPCIncomingContextExtractorFxOption if you only want to provide it.
var LoggerGRPCIncomingContextExtractor = context.LoggerGRPCIncomingContextExtractor

// LoggerGRPCInterceptorFxOption adds Unary Server Interceptor that will log Request and Response if needed
func LoggerGRPCInterceptorFxOption() fx.Option {
	return fx.Provide(fx.Annotated{
		Group:  groups.UnaryServerInterceptors,
		Target: server.LoggerGRPCInterceptor,
	})
}

// LoggerGRPCInterceptor is a constructor that creates gRPC Unary Server Interceptor.
// This Interceptor will log gRPC calls with request and response if enabled.
//
//	*Note* normally this dependency is part of a group. If you want to create it as a standalone
// 	dependency, remember that there can be only one of this kind in the graph.
//
// Consider using LoggerGRPCInterceptorFxOption if you only want to provide it.
var LoggerGRPCInterceptor = server.LoggerGRPCInterceptor
