package mortar

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildInfo(t *testing.T) {
	t.Run("without defaults", testGetBuildInformationWithoutDefaults)
	t.Run("without defaults and explanations", testGetBuildInformationWithoutDefaultsWithExplanations)
	injectValues()
	t.Run("with defaults", testGetBuildInformationWithDefaults)
}

func testGetBuildInformationWithDefaults(t *testing.T) {
	info := GetBuildInformation(true)
	assert.Equal(t, "1234", info.GitCommit)
	assert.Equal(t, "abc", info.BuildTag)
	assert.Equal(t, "v0.0.1", info.Version)
	expectedTime, _ := time.Parse(time.RFC3339, "2020-08-12T17:11:51Z")
	assert.Equal(t, expectedTime, info.BuildTime)
	hostname, _ := os.Hostname()
	assert.Equal(t, hostname, info.Hostname)
	assert.WithinDuration(t, time.Now(), info.InitTime, time.Second)
	assert.WithinDuration(t, time.Now(), time.Now().Add(time.Duration(info.UpTime)), time.Second)
}

func testGetBuildInformationWithoutDefaultsWithExplanations(t *testing.T) {
	info := GetBuildInformation(true)
	assert.Equal(t, "wasn't provided during build", info.GitCommit)
	assert.Equal(t, "wasn't provided during build", info.BuildTag)
	assert.Equal(t, "wasn't provided during build", info.Version)
	assert.Equal(t, time.Time{}, info.BuildTime)
	hostname, _ := os.Hostname()
	assert.Equal(t, hostname, info.Hostname)
	assert.WithinDuration(t, time.Now(), info.InitTime, time.Second)
	assert.WithinDuration(t, time.Now(), time.Now().Add(time.Duration(info.UpTime)), time.Second)
}

func testGetBuildInformationWithoutDefaults(t *testing.T) {
	info := GetBuildInformation()
	assert.Equal(t, "", info.GitCommit)
	assert.Equal(t, "", info.BuildTag)
	assert.Equal(t, "", info.Version)
	assert.Equal(t, time.Time{}, info.BuildTime)
	hostname, _ := os.Hostname()
	assert.Equal(t, hostname, info.Hostname)
	assert.WithinDuration(t, time.Now(), info.InitTime, time.Second)
	assert.WithinDuration(t, time.Now(), time.Now().Add(time.Duration(info.UpTime)), time.Second)
}

func injectValues() {
	gitCommit = "1234"
	buildTag = "abc"
	buildTimestamp = "2020-08-12T17:11:51Z"
	version = "v0.0.1"
}

func TestDurationMarshaler(t *testing.T) {
	info := GetBuildInformation()
	jsonBytes, err := json.Marshal(info)
	require.NoError(t, err)
	var infoMap map[string]interface{}
	err = json.Unmarshal(jsonBytes, &infoMap)
	require.NoError(t, err)
	assert.Regexp(t, ".+s$", infoMap["up_time"]) // ends with either `µs` or `ms`
}
