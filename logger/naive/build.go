package naive

import (
	"container/list"
	"io"
	"os"

	logInt "github.com/go-masonry/mortar/interfaces/log"
)

const defaultSkipDepth = 4

type defaultConfig struct {
	writer        io.Writer
	level         logInt.Level
	depth         int
	excludeTime   bool
	includeCaller bool
}

type defaultBuilder struct {
	ll *list.List
}

// NativeLogBuilder is a helper interface to configure native log.Logger instance.
type NativeLogBuilder interface {
	logInt.Builder
	// SetWriter set where output should be printed
	SetWriter(writer io.Writer) NativeLogBuilder
	// ExcludeTime configures standard Logger to exclude any time field
	ExcludeTime() NativeLogBuilder
	// IncludeCaller adds caller:line to the output
	IncludeCaller() NativeLogBuilder
}

// Builder creates a fresh default Logger builder, this will eventually build a std logger wrapper without structured logging
func Builder() NativeLogBuilder {
	return &defaultBuilder{
		ll: list.New(),
	}
}

func (d *defaultBuilder) SetWriter(writer io.Writer) NativeLogBuilder {
	d.ll.PushBack(func(cfg *defaultConfig) {
		cfg.writer = writer
	})
	return d
}

func (d *defaultBuilder) ExcludeTime() NativeLogBuilder {
	d.ll.PushBack(func(cfg *defaultConfig) {
		cfg.excludeTime = true
	})
	return d
}

func (d *defaultBuilder) IncludeCaller() NativeLogBuilder {
	d.ll.PushBack(func(cfg *defaultConfig) {
		cfg.includeCaller = true
	})
	return d
}

func (d *defaultBuilder) IncrementSkipFrames(inc int) logInt.Builder {
	d.ll.PushBack(func(cfg *defaultConfig) {
		cfg.depth += inc
	})
	return d
}

func (d *defaultBuilder) SetLevel(level logInt.Level) logInt.Builder {
	d.ll.PushBack(func(cfg *defaultConfig) {
		cfg.level = level
	})
	return d
}

func (d *defaultBuilder) Build() logInt.Logger {
	cfg := &defaultConfig{
		writer:        os.Stderr,
		level:         logInt.TraceLevel,
		depth:         defaultSkipDepth, // 2 is used within the log package
		excludeTime:   false,
		includeCaller: false,
	}
	for e := d.ll.Front(); e != nil; e = e.Next() {
		f := e.Value.(func(config *defaultConfig))
		f(cfg)
	}
	return newDefaultLogger(cfg)
}

var _ logInt.Builder = (*defaultBuilder)(nil)
