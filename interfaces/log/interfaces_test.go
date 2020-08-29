package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLevelString(t *testing.T) {
	levels := map[Level]string{
		TraceLevel: "trace",
		DebugLevel: "debug",
		InfoLevel:  "info",
		WarnLevel:  "warn",
		ErrorLevel: "error",
	}

	for lvl, str := range levels {
		assert.Equal(t, str, lvl.String())
	}
}

func TestParseLevel(t *testing.T) {
	levels := map[string]Level{
		"Trace":    TraceLevel,
		"dEBug":    DebugLevel,
		"INFO":     InfoLevel,
		"warn":     WarnLevel,
		"whatever": TraceLevel,
	}
	for str, lvl := range levels {
		assert.Equal(t, lvl, ParseLevel(str))
	}
}
