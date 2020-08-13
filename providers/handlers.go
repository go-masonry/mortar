package providers

import (
	"github.com/go-masonry/mortar/handlers"
	"github.com/go-masonry/mortar/providers/groups"
	"go.uber.org/fx"
)

// InternalDebugHandlersFxOption adds Internal Debug Handlers to the graph
func InternalDebugHandlersFxOption() fx.Option {
	return fx.Provide(fx.Annotated{
		Group:  groups.InternalHTTPHandlers + ",flatten",
		Target: handlers.InternalDebugHandlers,
	})
}

// InternalDebugHandlers is a constructor that creates Internal Debug HTTP Handlers
//	*Note* normally this dependency is part of a group. If you want to create it as a standalone
// 	dependency, remember that there can be only one of this kind in the graph.
//
// Consider using InternalDebugHandlersFxOption if you only want to provide it.
var InternalDebugHandlers = handlers.InternalDebugHandlers

// InternalProfileHandlerFunctionsFxOption adds Internal Profiler Handler to the graph
func InternalProfileHandlerFunctionsFxOption() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group:  groups.InternalHTTPHandlerFunctions + ",flatten",
			Target: handlers.InternalProfileHandlerFunctions,
		})
}

// InternalProfileHandlerFunctions is a constructor that creates Internal Profile HTTP Handlers
//	*Note* normally this dependency is part of a group. If you want to create it as a standalone
// 	dependency, remember that there can be only one of this kind in the graph.
//
// Consider using InternalProfileHandlerFunctionsFxOption if you only want to provide it.
var InternalProfileHandlerFunctions = handlers.InternalProfileHandlerFunctions

// InternalSelfHandlersFxOption adds Internal Self Build and Config information HTTP Handlers to the graph
//
// Adds these endpoint on Internal web service
//	- /self/build
//	- /self/config
func InternalSelfHandlersFxOption() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group:  groups.InternalHTTPHandlers + ",flatten",
			Target: handlers.SelfHandlers,
		})
}

// SelfHandlers is a constructor that creates Internal Self HTTP Handlers
//	*Note* normally this dependency is part of a group. If you want to create it as a standalone
// 	dependency, remember that there can be only one of this kind in the graph.
//
// Adds these endpoint on Internal web service
//	- /self/build
//	- /self/config
//
// Consider using InternalSelfHandlersFxOption if you only want to provide it.
var SelfHandlers = handlers.SelfHandlers
