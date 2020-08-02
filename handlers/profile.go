package handlers

import (
	"github.com/go-masonry/mortar/constructors/partial"
	"go.uber.org/fx"
	"net/http/pprof"
)

const (
	profilePrefix = internalPatternPrefix + "/debug/pprof"
)

func InternalProfileHandlerFunctionsOption() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group:  partial.FxGroupInternalHttpHandlerFunctions + ",flatten",
			Target: InternalProfileHandlerFunctions,
		})
}

func InternalProfileHandlerFunctions() []partial.HttpHandlerFuncPatternPair {
	return []partial.HttpHandlerFuncPatternPair{
		{Pattern: profilePrefix, HandlerFunc: pprof.Index},
		{Pattern: profilePrefix + "/cmdline", HandlerFunc: pprof.Cmdline},
		{Pattern: profilePrefix + "/profile", HandlerFunc: pprof.Profile},
		{Pattern: profilePrefix + "/symbol", HandlerFunc: pprof.Symbol},
		{Pattern: profilePrefix + "/trace", HandlerFunc: pprof.Trace},
	}
}
