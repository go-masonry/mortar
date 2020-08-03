package constructors

import (
	"github.com/go-masonry/mortar/interfaces/cfg"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/go-masonry/mortar/mortar"
	"go.uber.org/fx"
)

const FxGroupLoggerContextExtractors = "loggerContextExtractors"
const (
	application = "app"
	hostname    = "host"
	gitCommit   = "git"
)

type LoggerDeps struct {
	fx.In

	Config            cfg.Config
	LoggerBuilder     log.Builder
	ContextExtractors []log.ContextExtractor `group:"loggerContextExtractors"`
}

// DefaultLogger is a constructor that will create a logger with some default values on top of provided ones
func DefaultLogger(deps LoggerDeps) log.Logger {
	var logLevel = log.InfoLevel
	if levelValue := deps.Config.Get(mortar.LoggerLevelKey); levelValue.IsSet() {
		logLevel = log.ParseLevel(levelValue.String())
	}
	appName := deps.Config.Get(mortar.Name).String() // empty string is just fine
	return deps.LoggerBuilder.
		SetLevel(logLevel).
		AddStaticFields(selfStaticFields(appName)).
		AddContextExtractors(deps.ContextExtractors...).
		Build()
}

func selfStaticFields(name string) map[string]interface{} {
	output := make(map[string]interface{})
	info := mortar.GetBuildInformation()
	if len(name) > 0 {
		output[application] = name
	}
	if len(info.Hostname) > 0 {
		output[hostname] = info.Hostname
	}
	if len(info.GitCommit) > 0 {
		output[gitCommit] = info.GitCommit
	}
	return output
}
