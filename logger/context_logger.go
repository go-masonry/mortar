package logger

import (
	"context"

	"github.com/go-masonry/mortar/interfaces/log"
)

type contextAwareLogEntry struct {
	contextExtractors []log.ContextExtractor
	innerLogger       log.Fields
	fields            map[string]interface{}
	err               error
}

func newEntry(contextExtractors []log.ContextExtractor, logger log.Fields) log.Fields {
	return &contextAwareLogEntry{
		contextExtractors: contextExtractors,
		innerLogger:       logger,
		fields:            make(map[string]interface{}),
		err:               nil,
	}
}

func (c *contextAwareLogEntry) Trace(ctx context.Context, format string, args ...interface{}) {
	c.log(ctx, log.TraceLevel, format, args...)
}

func (c *contextAwareLogEntry) Debug(ctx context.Context, format string, args ...interface{}) {
	c.log(ctx, log.DebugLevel, format, args...)
}

func (c *contextAwareLogEntry) Info(ctx context.Context, format string, args ...interface{}) {
	c.log(ctx, log.InfoLevel, format, args...)
}

func (c *contextAwareLogEntry) Warn(ctx context.Context, format string, args ...interface{}) {
	c.log(ctx, log.WarnLevel, format, args...)
}

func (c *contextAwareLogEntry) Error(ctx context.Context, format string, args ...interface{}) {
	c.log(ctx, log.ErrorLevel, format, args...)
}

func (c *contextAwareLogEntry) Custom(ctx context.Context, level log.Level, format string, args ...interface{}) {
	c.log(ctx, level, format, args...)
}

func (c *contextAwareLogEntry) WithError(err error) log.Fields {
	c.err = err
	return c
}

func (c *contextAwareLogEntry) WithField(name string, value interface{}) log.Fields {
	c.fields[name] = value
	return c
}

func (c *contextAwareLogEntry) log(ctx context.Context, level log.Level, format string, args ...interface{}) {
	if ctx == nil {
		ctx = context.Background()
	}
	logger := c.enrich(ctx)
	for k, v := range c.fields {
		logger = logger.WithField(k, v)
	}
	if c.err != nil {
		logger = logger.WithError(c.err)
	}
	logger.Custom(ctx, level, format, args...)
}

func (c *contextAwareLogEntry) enrich(ctx context.Context) (logger log.Fields) {
	defer func() {
		if r := recover(); r != nil {
			c.innerLogger.WithField("__panic__", r).Error(ctx, "one of the context extractors panicked")
			logger = c.innerLogger
		}
	}()
	logger = c.innerLogger
	for _, extractor := range c.contextExtractors {
		for k, v := range extractor(ctx) {
			logger = logger.WithField(k, v)
		}
	}
	return
}
