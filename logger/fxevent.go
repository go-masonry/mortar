package logger

import (
	"context"
	"strings"

	"github.com/go-masonry/mortar/interfaces/cfg"
	confkeys "github.com/go-masonry/mortar/interfaces/cfg/keys"
	"github.com/go-masonry/mortar/interfaces/log"
	"go.uber.org/fx/fxevent"
)

// CreateFxEventLogger is a constructor to create fxevent.Logger
// This one is used by fx itself to output its events
func CreateFxEventLogger(logger log.Logger, cfg cfg.Config) fxevent.Logger {
	level := log.InfoLevel
	if cfg.Get(confkeys.LogStartStop).IsSet() {
		level = log.ParseLevel(cfg.Get(confkeys.LogStartStop).String())
	}
	return &logWrapper{Logger: logger, startStopLogLevel: level}
}

type logWrapper struct {
	log.Logger
	startStopLogLevel log.Level
}

func (zw *logWrapper) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		zw.
			WithField("callee", e.FunctionName).
			WithField("caller", e.CallerName).
			Debug(context.TODO(), "OnStart hook executing")
	case *fxevent.OnStartExecuted:
		logger := zw.
			WithError(e.Err).
			WithField("callee", e.FunctionName).
			WithField("caller", e.CallerName)
		if e.Err != nil {
			logger.Error(context.TODO(), "OnStart hook failed")
		} else {
			logger.
				WithField("runtime", e.Runtime.String()).
				Debug(context.TODO(), "OnStart hook executed")
		}
	case *fxevent.OnStopExecuting:
		zw.
			WithField("callee", e.FunctionName).
			WithField("caller", e.CallerName).
			Debug(context.TODO(), "OnStop hook executing")
	case *fxevent.OnStopExecuted:
		logger := zw.
			WithError(e.Err).
			WithField("callee", e.FunctionName).
			WithField("caller", e.CallerName)
		if e.Err != nil {
			logger.Error(context.TODO(), "OnStop hook failed")
		} else {
			logger.
				WithField("runtime", e.Runtime.String()).
				Debug(context.TODO(), "OnStop hook executed")
		}
	case *fxevent.Supplied:
		zw.
			WithField("type", e.TypeName).
			WithError(e.Err).
			Debug(context.TODO(), "supplied")
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			zw.
				WithField("constructor", e.ConstructorName).
				WithField("type", rtype).
				Debug(context.TODO(), "provided")
		}
		if e.Err != nil {
			zw.
				WithError(e.Err).
				Error(context.TODO(), "error encountered while applying options")
		}
	case *fxevent.Invoking:
		// Do nothing. Will log on Invoked.

	case *fxevent.Invoked:
		logger := zw.
			WithError(e.Err).
			WithField("stack", e.Trace).
			WithField("function", e.FunctionName)
		if e.Err != nil {
			logger.Error(context.TODO(), "invoke failed")
		} else {
			logger.Debug(context.TODO(), "invoked")
		}
	case *fxevent.Stopping:
		zw.
			WithField("signal", strings.ToUpper(e.Signal.String())).
			Custom(context.TODO(), zw.startStopLogLevel, 0, "received termination signal")
	case *fxevent.Stopped:
		if e.Err != nil {
			zw.WithError(e.Err).Error(context.TODO(), "stop failed")
		}
	case *fxevent.RollingBack:
		zw.
			WithError(e.StartErr).
			Error(context.TODO(), "start failed, rolling back")
	case *fxevent.RolledBack:
		if e.Err != nil {
			zw.
				WithError(e.Err).
				Error(context.TODO(), "start failed, rolling back")
		}
	case *fxevent.Started:
		if e.Err != nil {
			zw.
				WithError(e.Err).
				Error(context.TODO(), "start failed")
		} else {
			zw.Custom(context.TODO(), zw.startStopLogLevel, 0, "service started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			zw.
				WithError(e.Err).
				Error(context.TODO(), "custom logger initialization failed")
		} else {
			zw.
				WithField("function", e.ConstructorName).
				Debug(context.TODO(), "initialized custom fxevent.Logger")
		}
	}
}
