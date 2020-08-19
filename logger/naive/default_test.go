package naive

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"testing"

	logInt "github.com/go-masonry/mortar/interfaces/log"
	"github.com/stretchr/testify/assert"
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
	logger := Builder().SetWriter(&buf).IncludeCaller().Build()
	logger.Info(nil, "with caller")
	assert.Regexp(t, `^\d{4}\/\d{2}\/\d{2}.+\d{2}\:\d{2}\:\d{2} .+naive/default_test\.go\:\d{2}\: with caller`, buf.String())
}

func TestLevelHigher(t *testing.T) {
	var buf bytes.Buffer
	logger := Builder().SetWriter(&buf).SetLevel(logInt.ErrorLevel).Build()
	logger.Warn(nil, "nothing printed")
	assert.Empty(t, buf.String())
}

func TestConfiguration(t *testing.T) {
	var buf bytes.Buffer
	logger := Builder().SetWriter(&buf).IncludeCaller().SetLevel(logInt.InfoLevel).IncrementSkipFrames(2).Build()
	configuration := logger.Configuration()
	assert.Equal(t, logInt.InfoLevel, configuration.Level())
	assert.IsType(t, &log.Logger{}, configuration.Implementation())
}

func TestNotSupported(t *testing.T) {
	var buf bytes.Buffer
	logger := Builder().SetWriter(&buf).SetLevel(logInt.ErrorLevel).Build()
	logger.WithError(fmt.Errorf("an error")).WithField("one", "two").Error(nil, "warning line")
	assert.NotContains(t, buf.String(), "an error")
	assert.NotContains(t, buf.String(), "one")
	assert.NotContains(t, buf.String(), "two")
	assert.Contains(t, buf.String(), "warning line")
}
