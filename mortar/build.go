package mortar

import (
	"fmt"
	"os"
	"time"
)

// Some of these variables should be populated during build using LDFLAGS
var (
	// should be injected with LDFLAG
	gitCommit string
	// should be injected with LDFLAG
	version string
	// should be injected with LDFLAG
	buildTimestamp string // "2006-01-02T15:04:05Z07:00" defined in RFC3339
	// should be injected with LDFLAG
	buildTag string

	// initialized During init()
	initTime time.Time
	// initialized During init()
	hostname string
)

func init() {
	initTime = time.Now()
	if host, err := os.Hostname(); err == nil {
		hostname = host
	} else {
		hostname = err.Error()
	}
}

// Information is a struct that will hold all the statically injected information during build
type Information struct {
	GitCommit string       `json:"git_commit,omitempty"`
	Version   string       `json:"version,omitempty"`
	BuildTag  string       `json:"build_tag,omitempty"`
	BuildTime time.Time    `json:"build_time,omitempty"`
	InitTime  time.Time    `json:"init_time,omitempty"`
	UpTime    JSONDuration `json:"up_time,omitempty"`
	Hostname  string       `json:"hostname,omitempty"`
}

// GetBuildInformation returns this service build information
func GetBuildInformation(includeExplanations ...bool) (info Information) {
	info.GitCommit = gitCommit
	info.Version = version
	info.BuildTag = buildTag
	info.InitTime = initTime
	info.UpTime = JSONDuration(time.Since(initTime))
	info.BuildTime = time.Time{} // Zero
	info.Hostname = hostname
	if len(buildTimestamp) != 0 {
		// try to parse
		if t, err := time.Parse(time.RFC3339, buildTimestamp); err == nil {
			info.BuildTime = t
		}
	}

	if len(includeExplanations) > 0 && includeExplanations[0] {
		if len(gitCommit) == 0 {
			info.GitCommit = "wasn't provided during build"
		}
		if len(version) == 0 {
			info.Version = "wasn't provided during build"
		}
		if len(buildTag) == 0 {
			info.BuildTag = "wasn't provided during build"
		}
	}
	return
}

// JSONDuration is an alias to time.Duration for Json marshaling
type JSONDuration time.Duration

// MarshalJSON for JsonDuration is a helper function to better marshal time.Duration
func (jd JSONDuration) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, time.Duration(jd))), nil
}
