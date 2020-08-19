package naive

import (
	"context"
	"fmt"
	"log"

	logInt "github.com/go-masonry/mortar/interfaces/log"
)

type defaultLogger struct {
	cfg    *defaultConfig
	logger *log.Logger
}

func (d *defaultLogger) Level() logInt.Level {
	return d.cfg.level
}

func (d *defaultLogger) Implementation() interface{} {
	return d.logger
}

func (d *defaultLogger) Trace(ctx context.Context, format string, args ...interface{}) {
	d.Custom(ctx, logInt.TraceLevel, format, args...)
}

func (d *defaultLogger) Debug(ctx context.Context, format string, args ...interface{}) {
	d.Custom(ctx, logInt.DebugLevel, format, args...)
}

func (d *defaultLogger) Info(ctx context.Context, format string, args ...interface{}) {
	d.Custom(ctx, logInt.InfoLevel, format, args...)
}

func (d *defaultLogger) Warn(ctx context.Context, format string, args ...interface{}) {
	d.Custom(ctx, logInt.WarnLevel, format, args...)
}

func (d *defaultLogger) Error(ctx context.Context, format string, args ...interface{}) {
	d.Custom(ctx, logInt.ErrorLevel, format, args...)
}

func (d *defaultLogger) Custom(ctx context.Context, level logInt.Level, format string, args ...interface{}) {
	if d.cfg.level <= level {
		d.log(format, args...)
	}
}

func (d *defaultLogger) WithError(err error) logInt.Fields {
	// not supported
	return d
}

func (d *defaultLogger) WithField(name string, value interface{}) logInt.Fields {
	// not supported
	return d
}

func (d *defaultLogger) Configuration() logInt.LoggerConfiguration {
	return d
}

func (d *defaultLogger) log(format string, args ...interface{}) {
	if len(args) > 0 {
		d.logger.Output(d.cfg.depth, fmt.Sprintf(format, args...))
	} else {
		d.logger.Output(d.cfg.depth, format)
	}
}

func newDefaultLogger(cfg *defaultConfig) logInt.Logger {
	flags := log.LstdFlags
	if cfg.excludeTime {
		flags = 0
	}
	if cfg.includeCaller {
		flags |= log.Llongfile
	}
	logger := log.New(cfg.writer, "", flags)
	return &defaultLogger{
		logger: logger,
		cfg:    cfg,
	}
}

var _ logInt.Logger = (*defaultLogger)(nil)
