package logger

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

// Builder creates a fresh default Logger builder, this will eventually build a std logger wrapper without structured logging
func Builder() logInt.Builder {
	return &defaultBuilder{
		ll: list.New(),
	}
}

func (d *defaultBuilder) SetWriter(writer io.Writer) logInt.Builder {
	d.ll.PushBack(func(cfg *defaultConfig) {
		cfg.writer = writer
	})
	return d
}

func (d *defaultBuilder) SetLevel(level logInt.Level) logInt.Builder {
	d.ll.PushBack(func(cfg *defaultConfig) {
		cfg.level = level
	})
	return d
}

func (d *defaultBuilder) AddStaticFields(fields map[string]interface{}) logInt.Builder {
	return d
}

func (d *defaultBuilder) AddContextExtractors(hooks ...logInt.ContextExtractor) logInt.Builder {
	return d
}

func (d *defaultBuilder) ExcludeTime() logInt.Builder {
	d.ll.PushBack(func(cfg *defaultConfig) {
		cfg.excludeTime = true
	})
	return d
}

func (d *defaultBuilder) SetCustomTimeFormatter(format string) logInt.Builder {
	return d
}

func (d *defaultBuilder) IncludeCallerAndSkipFrames(skip int) logInt.Builder {
	d.ll.PushBack(func(cfg *defaultConfig) {
		cfg.depth = defaultSkipDepth + skip // 2 is used within the log package
		cfg.includeCaller = true
	})
	return d
}

func (d *defaultBuilder) Build() logInt.Logger {
	cfg := &defaultConfig{
		writer:        os.Stderr,
		level:         logInt.TraceLevel,
		depth:         defaultSkipDepth,
		excludeTime:   false,
		includeCaller: false,
	}
	for e := d.ll.Front(); e != nil; e = e.Next() {
		f := e.Value.(func(config *defaultConfig))
		f(cfg)
	}
	return newDefaultLogger(cfg)
}
