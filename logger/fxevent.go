package logger

import (
	"context"
	"strings"

	"github.com/go-masonry/mortar/interfaces/log"
	"go.uber.org/fx/fxevent"
)

func CreateFxEventLogger(log log.Logger) fxevent.Logger {
	return &logWrapper{log}
}

type logWrapper struct{ log.Logger }

func (zw *logWrapper) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		zw.
			WithField("callee", e.FunctionName).
			WithField("caller", e.CallerName).
			Info(context.TODO(), "OnStart hook executing")
	case *fxevent.OnStartExecuted:
		logger := zw.
			WithError(e.Err).
			WithField("callee", e.FunctionName).
			WithField("caller", e.CallerName)
		if e.Err != nil {
			logger.Info(context.TODO(), "OnStart hook failed")
		} else {
			logger.
				WithField("runtime", e.Runtime.String()).
				Info(context.TODO(), "OnStart hook executed")
		}
	case *fxevent.OnStopExecuting:
		zw.
			WithField("callee", e.FunctionName).
			WithField("caller", e.CallerName).
			Info(context.TODO(), "OnStop hook executing")
	case *fxevent.OnStopExecuted:
		logger := zw.
			WithError(e.Err).
			WithField("callee", e.FunctionName).
			WithField("caller", e.CallerName)
		if e.Err != nil {
			logger.Info(context.TODO(), "OnStop hook failed")
		} else {
			logger.
				WithField("runtime", e.Runtime.String()).
				Info(context.TODO(), "OnStop hook executed")
		}
	case *fxevent.Supplied:
		zw.
			WithField("type", e.TypeName).
			WithError(e.Err).
			Info(context.TODO(), "supplied")
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			zw.
				WithField("constructor", e.ConstructorName).
				WithField("type", rtype).
				Info(context.TODO(), "provided")
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
			logger.Info(context.TODO(), "invoked")
		}
	case *fxevent.Stopping:
		zw.
			WithField("signal", strings.ToUpper(e.Signal.String())).
			Info(context.TODO(), "received signal")
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
			zw.Info(context.TODO(), "started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			zw.
				WithError(e.Err).
				Error(context.TODO(), "custom logger initialization failed")
		} else {
			zw.
				WithField("function", e.ConstructorName).
				Info(context.TODO(), "initialized custom fxevent.Logger")
		}
	}
}
