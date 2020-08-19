package log

import (
	"context"
	"strings"
)

//go:generate mockgen -source=interfaces.go -destination=mock/mock.go

// Level is a log level enum
type Level int8

const (
	// TraceLevel defines trace log level.
	TraceLevel Level = iota
	// DebugLevel defines debug log level.
	DebugLevel
	// InfoLevel defines info log level.
	InfoLevel
	// WarnLevel defines warn log level.
	WarnLevel
	// ErrorLevel defines error log level.
	ErrorLevel
)

func (l Level) String() string {
	switch l {
	case ErrorLevel:
		return "error"
	case WarnLevel:
		return "warn"
	case InfoLevel:
		return "info"
	case DebugLevel:
		return "debug"
	default:
		return "trace"
	}
}

// ParseLevel tries to parse level from string, if unable to parse a Trace level will be returned as a default
func ParseLevel(str string) Level {
	switch strings.ToLower(str) {
	case "error":
		return ErrorLevel
	case "warn":
		return WarnLevel
	case "info":
		return InfoLevel
	case "debug":
		return DebugLevel
	default:
		return TraceLevel
	}
}

// LoggerConfiguration get some of the logger configuration options and also the implementation
type LoggerConfiguration interface {
	Level() Level
	Implementation() interface{}
}

// ContextExtractor is an alias for a function that extract values from context
// Make sure that this function returns fast and is "thread safe"
type ContextExtractor func(ctx context.Context) map[string]interface{}

// Builder defines log configuration options
type Builder interface {
	// IncrementSkipFrames **peels** an additional layer(s) to show the actual log line position.
	//
	// **Note**
	//
	// This one is really important, the implementation must increment and not override
	// the value of skip frames.
	//
	// One should take into account only it's own frame stack, mortar will add 2 on top.
	IncrementSkipFrames(addition int) Builder
	// Set system log level
	//
	// Optional:
	//	Logger Implementation should have a default log level. This allows filtering log lines with level below default one.
	SetLevel(level Level) Builder
	// Build() returns a Logger implementation, always
	Build() Logger
}

// Messages part of the Logger interface
type Messages interface {
	// Highly detailed tracing messages. Produces the most voluminous output. Used by developers for developers
	// Some implementations doesn't have that granularity and use Debug level instead
	//
	// Note:
	// 	ctx can be nil
	Trace(ctx context.Context, format string, args ...interface{})
	// Relatively detailed tracing messages. Used mostly by developers to debug the flow
	//
	// Note:
	// 	ctx can be nil
	Debug(ctx context.Context, format string, args ...interface{})
	// Informational messages that might make sense to users unfamiliar with this application
	//
	// Note:
	// 	ctx can be nil
	Info(ctx context.Context, format string, args ...interface{})
	// Potentially harmful situations of interest to users that indicate potential problems.
	//
	// Note:
	// 	ctx can be nil
	Warn(ctx context.Context, format string, args ...interface{})
	// Very severe error events that might cause the application to terminate or misbehave.
	// It's not intended to use to log every 'error', for that use 'WithError(err).<Trace|Debug|Info|...>(...)'
	//
	// Note:
	// 	ctx can be nil
	Error(ctx context.Context, format string, args ...interface{})
	// Custom is a Convenience function that will enable you to set log level dynamically
	//  and perhaps skip additional frames.
	//
	// Note:
	// 	ctx can be nil
	Custom(ctx context.Context, level Level, skipAdditionalFrames int, format string, args ...interface{})
}

// Fields part of the Logger interface
type Fields interface {
	Messages
	// Add an error to the log structure, output depends on the implementation
	WithError(err error) Fields
	// Add an informative field to the log structure, output depends on the implementation
	//
	// Avoid using the following names:
	// 		time, message, error, caller, stack
	WithField(name string, value interface{}) Fields
}

// Logger is a simple interface that defines logging in the system
type Logger interface {
	Fields
	// Implementor returns the actual lib/struct that is responsible for the above logic
	Configuration() LoggerConfiguration
}
