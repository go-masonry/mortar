package logger

import (
	"bytes"
	"context"
	"fmt"
	logInt "github.com/go-masonry/mortar/interfaces/log"
	"github.com/stretchr/testify/assert"
	"log"
	"strings"
	"testing"
	"time"
)

func TestDefaultLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := Builder().SetWriter(&buf).Build()
	logger.Trace(nil, "trace %s", "line")
	assert.Contains(t, buf.String(), "trace line")
}

func TestNoTimeNoCaller(t *testing.T) {
	var buf bytes.Buffer
	logger := Builder().SetWriter(&buf).ExcludeTime().Build()
	logger.Debug(nil, "no time")
	assert.Equal(t, "no time", strings.TrimSpace(buf.String()))
}

func TestTimeAndCaller(t *testing.T) {
	var buf bytes.Buffer
	logger := Builder().SetWriter(&buf).IncludeCallerAndSkipFrames(0).Build()
	logger.Info(nil, "with caller")
	assert.Regexp(t, `^\d{4}\/\d{2}\/\d{2}.+\d{2}\:\d{2}\:\d{2} .+default_test\.go\:\d{2}\: with caller`, buf.String())
}

func TestLevelHigher(t *testing.T) {
	var buf bytes.Buffer
	logger := Builder().SetWriter(&buf).SetLevel(logInt.ErrorLevel).Build()
	logger.Warn(nil, "nothing printed")
	assert.Empty(t, buf.String())
}

func TestFieldsErrorNoEffect(t *testing.T) {
	var buf bytes.Buffer
	logger := Builder().SetWriter(&buf).ExcludeTime().AddStaticFields(
		map[string]interface{}{
			"one": 1,
		},
	).Build()
	logger.WithError(fmt.Errorf("error")).WithField("field", "absent").Error(nil, "no fields")
	assert.Equal(t, "no fields", strings.TrimSpace(buf.String()))
}

func TestConfiguration(t *testing.T) {
	var buf bytes.Buffer
	logger := Builder().SetWriter(&buf).SetLevel(logInt.InfoLevel).SetCustomTimeFormatter(time.RFC3339).AddContextExtractors(
		func(ctx context.Context) map[string]interface{} {
			panic("no way")
		},
	).IncludeCallerAndSkipFrames(2).Build()
	configuration := logger.Configuration()
	assert.Same(t, &buf, configuration.Writer())
	assert.Equal(t, logInt.InfoLevel, configuration.Level())
	assert.Nil(t, configuration.ContextExtractors())
	assert.IsType(t, &log.Logger{}, configuration.Implementation())
	includeCaller, depth := configuration.CallerConfiguration()
	assert.True(t, includeCaller)
	assert.Equal(t, 2, depth)
	timeFieldConfiguration, pattern := configuration.TimeFieldConfiguration()
	assert.False(t, timeFieldConfiguration)
	assert.Empty(t, pattern)
}
