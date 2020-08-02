package constructors

import (
	"github.com/go-masonry/mortar/interfaces/cfg"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/go-masonry/mortar/mortar"
	"go.uber.org/fx"
)

const FxGroupLoggerContextExtractors = "loggerContextExtractors"

type LoggerDeps struct {
	fx.In
	Config            cfg.Config
	LoggerBuilder     log.Builder
	ContextExtractors []log.ContextExtractor `group:"loggerContextExtractors"`
}

// DefaultLogger is a constructor that will create a logger with some default values, however
// you can still customize it.
//
// 	- Level: Given config.Config map, we will look for self.LoggerLevelKey value or use Builder default
// 	- ContextExtractors: Since we are using uber.Fx for DI we can expect any number of context extractors
//		All context extractors must be grouped under a fx.Group named: 'loggerExtractors'
func DefaultLogger(deps LoggerDeps) log.Logger {
	var logLevel = log.InfoLevel
	if levelValue := deps.Config.Get(mortar.LoggerLevelKey); levelValue.IsSet() {
		logLevel = log.ParseLevel(levelValue.String())
	}
	return deps.LoggerBuilder.
		SetLevel(logLevel).
		AddContextExtractors(deps.ContextExtractors...).
		Build()
}
