package logger

import (
	"context"

	"github.com/go-masonry/mortar/interfaces/log"
)

const (
	compensateMortarLoggerWrapper = 1
)

type loggerWrapper struct {
	contextExtractors []log.ContextExtractor
	logger            log.Logger
}

// CreateMortarLogger creates a new mortar logger which is a wrapper to support
// 	- ContextExtractors
//
// **Important**
//	This constructor will call builder.IncrementSkipFrames to peel additional layer of itself.
func CreateMortarLogger(builder log.Builder, contextExtractors ...log.ContextExtractor) log.Logger {
	logger := builder.IncrementSkipFrames(compensateMortarLoggerWrapper).Build() // add 1
	return &loggerWrapper{
		contextExtractors: contextExtractors,
		logger:            logger,
	}
}

func (l *loggerWrapper) Trace(ctx context.Context, format string, args ...interface{}) {
	newEntry(l.contextExtractors, l.logger, false).Trace(ctx, format, args...)
}

func (l *loggerWrapper) Debug(ctx context.Context, format string, args ...interface{}) {
	newEntry(l.contextExtractors, l.logger, false).Debug(ctx, format, args...)
}

func (l *loggerWrapper) Info(ctx context.Context, format string, args ...interface{}) {
	newEntry(l.contextExtractors, l.logger, false).Info(ctx, format, args...)
}

func (l *loggerWrapper) Warn(ctx context.Context, format string, args ...interface{}) {
	newEntry(l.contextExtractors, l.logger, false).Warn(ctx, format, args...)
}

func (l *loggerWrapper) Error(ctx context.Context, format string, args ...interface{}) {
	newEntry(l.contextExtractors, l.logger, false).Error(ctx, format, args...)
}

func (l *loggerWrapper) Custom(ctx context.Context, level log.Level, skipAdditionalFrames int, format string, args ...interface{}) {
	newEntry(l.contextExtractors, l.logger, false).Custom(ctx, level, skipAdditionalFrames, format, args...)
}

func (l *loggerWrapper) WithError(err error) log.Fields {
	return newEntry(l.contextExtractors, l.logger, true).WithError(err)
}

func (l *loggerWrapper) WithField(name string, value interface{}) log.Fields {
	return newEntry(l.contextExtractors, l.logger, true).WithField(name, value)
}

func (l *loggerWrapper) Configuration() log.LoggerConfiguration {
	return l.logger.Configuration()
}
