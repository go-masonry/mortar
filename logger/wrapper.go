package logger

import (
	"context"

	"github.com/go-masonry/mortar/interfaces/log"
)

const (
	incrementSkipFrames = 2
)

type loggerWrapper struct {
	contextExtractors []log.ContextExtractor
	inner             log.Logger
}

// CreateMortarLogger creates a new mortar logger which is a wrapper to support
// 	- ContextExtractors
//
// **Important**
//	This constructor will call builder.IncrementSkipFrames to peel additional layer of itself.
func CreateMortarLogger(builder log.Builder, contextExtractors ...log.ContextExtractor) log.Logger {
	logger := builder.IncrementSkipFrames(incrementSkipFrames).Build()
	return &loggerWrapper{
		contextExtractors: contextExtractors,
		inner:             logger,
	}
}

func (l *loggerWrapper) Trace(ctx context.Context, format string, args ...interface{}) {
	newEntry(l.contextExtractors, l.inner).Trace(ctx, format, args...)
}

func (l *loggerWrapper) Debug(ctx context.Context, format string, args ...interface{}) {
	newEntry(l.contextExtractors, l.inner).Debug(ctx, format, args...)
}

func (l *loggerWrapper) Info(ctx context.Context, format string, args ...interface{}) {
	newEntry(l.contextExtractors, l.inner).Info(ctx, format, args...)
}

func (l *loggerWrapper) Warn(ctx context.Context, format string, args ...interface{}) {
	newEntry(l.contextExtractors, l.inner).Warn(ctx, format, args...)
}

func (l *loggerWrapper) Error(ctx context.Context, format string, args ...interface{}) {
	newEntry(l.contextExtractors, l.inner).Error(ctx, format, args...)
}

func (l *loggerWrapper) Custom(ctx context.Context, level log.Level, format string, args ...interface{}) {
	newEntry(l.contextExtractors, l.inner).Custom(ctx, level, format, args...)
}

func (l *loggerWrapper) WithError(err error) log.Fields {
	return newEntry(l.contextExtractors, l.inner).WithError(err)
}

func (l *loggerWrapper) WithField(name string, value interface{}) log.Fields {
	return newEntry(l.contextExtractors, l.inner).WithField(name, value)
}

func (l *loggerWrapper) Configuration() log.LoggerConfiguration {
	return l.inner.Configuration()
}
