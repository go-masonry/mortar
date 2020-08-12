package handlers

import (
	"net/http/pprof"

	"github.com/go-masonry/mortar/constructors/partial"
)

const (
	profilePrefix = internalPatternPrefix + "/debug/pprof"
)

// InternalProfileHandlerFunctions profile handlers
func InternalProfileHandlerFunctions() []partial.HTTPHandlerFuncPatternPair {
	return []partial.HTTPHandlerFuncPatternPair{
		{Pattern: profilePrefix, HandlerFunc: pprof.Index},
		{Pattern: profilePrefix + "/cmdline", HandlerFunc: pprof.Cmdline},
		{Pattern: profilePrefix + "/profile", HandlerFunc: pprof.Profile},
		{Pattern: profilePrefix + "/symbol", HandlerFunc: pprof.Symbol},
		{Pattern: profilePrefix + "/trace", HandlerFunc: pprof.Trace},
	}
}
