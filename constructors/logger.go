package constructors

import (
	"github.com/go-masonry/mortar/interfaces/cfg"
	logInt "github.com/go-masonry/mortar/interfaces/log"
	defaultLogger "github.com/go-masonry/mortar/logger"
	"github.com/go-masonry/mortar/mortar"
	"go.uber.org/fx"
	"log"
)

const FxGroupLoggerContextExtractors = "loggerContextExtractors"
const (
	application = "app"
	hostname    = "host"
	gitCommit   = "git"

	callerSkipDepth = 0
)

type LoggerDeps struct {
	fx.In

	Config            cfg.Config
	LoggerBuilder     logInt.Builder            `optional:"true"`
	ContextExtractors []logInt.ContextExtractor `group:"loggerContextExtractors"`
}

// DefaultLogger is a constructor that will create a logger with some default values on top of provided ones
func DefaultLogger(deps LoggerDeps) logInt.Logger {
	var logLevel = logInt.InfoLevel
	if levelValue := deps.Config.Get(mortar.LoggerLevelKey); levelValue.IsSet() {
		logLevel = logInt.ParseLevel(levelValue.String())
	}
	return deps.getLogBuilder().
		SetLevel(logLevel).
		AddStaticFields(deps.selfStaticFields()).
		AddContextExtractors(deps.ContextExtractors...).
		IncludeCallerAndSkipFrames(callerSkipDepth).
		Build()
}

func (d LoggerDeps) selfStaticFields() map[string]interface{} {
	output := make(map[string]interface{})
	info := mortar.GetBuildInformation()
	appName := d.Config.Get(mortar.Name).String()
	if len(appName) > 0 && d.Config.Get(mortar.LoggerStaticName).Bool() {
		output[application] = appName
	}
	if len(info.Hostname) > 0 && d.Config.Get(mortar.LoggerStaticHost).Bool() {
		output[hostname] = info.Hostname
	}
	if len(info.GitCommit) > 0 && d.Config.Get(mortar.LoggerStaticGit).Bool() {
		output[gitCommit] = info.GitCommit
	}
	return output
}

func (d LoggerDeps) getLogBuilder() logInt.Builder {
	if d.LoggerBuilder != nil {
		return d.LoggerBuilder
	}
	log.Printf("[Mortar] WARNING \tNo logger builder provided, using default logger. Some features will not be supported")
	return defaultLogger.Builder()
}
