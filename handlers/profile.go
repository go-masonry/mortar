package handlers

import (
	"net/http/pprof"

	"github.com/go-masonry/mortar/constructors/partial"
)

const (
	// the path of pprof must start with /debug/pprof because of https://github.com/golang/go/issues/14286
	profilePrefix = internalPatternPrefix + "/pprof"
)

// InternalProfileHandlerFunctions profile handlers
func InternalProfileHandlerFunctions() []partial.HTTPHandlerFuncPatternPair {
	return []partial.HTTPHandlerFuncPatternPair{
		{Pattern: profilePrefix + "/", HandlerFunc: pprof.Index},
		{Pattern: profilePrefix + "/cmdline", HandlerFunc: pprof.Cmdline},
		{Pattern: profilePrefix + "/profile", HandlerFunc: pprof.Profile},
		{Pattern: profilePrefix + "/symbol", HandlerFunc: pprof.Symbol},
		{Pattern: profilePrefix + "/trace", HandlerFunc: pprof.Trace},
	}
}
