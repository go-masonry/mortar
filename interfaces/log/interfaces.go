package log

import (
	"context"
	"io"
	"strings"
)

//go:generate mockgen -source=interfaces.go -destination=mock/mock.go

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

// If unable to parse a Trace level will be returned as a default
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

type LoggerConfiguration interface {
	Writer() io.Writer
	Level() Level
	ContextExtractors() []ContextExtractor
	TimeFieldConfiguration() (bool, string)
	CallerConfiguration() (bool, int)
	Implementation() interface{}
}

// Make sure that this function returns fast and is "thread safe"
type ContextExtractor func(ctx context.Context) map[string]interface{}

type Builder interface {
	// TODO Allow adding static values to LOG such as "appname, host, etc"
	// Set output writer [os.Stderr, os.Stdout, bufio.Writer, ...]
	SetWriter(io.Writer) Builder
	// Set system log level
	// This might filter all logs where it's corresponding log level is lower than what is set here
	SetLevel(level Level) Builder
	// Add static fields to each log. Example: Application name, host, git commit
	AddStaticFields(fields map[string]interface{}) Builder
	// Make sure that each extractor function returns fast and is "thread safe"
	AddContextExtractors(hooks ...ContextExtractor) Builder
	// ExcludeTime removes implicit time field with it's value from the log
	ExcludeTime() Builder
	// SetCustomTimeFormatter allows to set custom time formatter, applicable if ExcludeTime() wasn't called
	SetCustomTimeFormatter(format string) Builder
	// IncludeCaller adds path and row number
	// One should only take into account it's own frame stack
	IncludeCallerAndSkipFrames(skip int) Builder
	// Build() returns a Logger implementation, always
	Build() Logger
}

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
	//
	// Note:
	// 	ctx can be nil
	Custom(ctx context.Context, level Level, format string, args ...interface{})
}

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
