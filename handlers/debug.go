package handlers

import (
	"encoding/json"
	"expvar"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/go-masonry/mortar/constructors/partial"
	"github.com/go-masonry/mortar/interfaces/log"
	"go.uber.org/fx"
)

const (
	internalPatternPrefix = "/debug"
)

// StatsInfo some statistics information
type StatsInfo struct {
	Memory          *runtime.MemStats `json:"memory"`
	NumOfCPU        int               `json:"num_of_cpu"`
	NumOfGoRoutines int               `json:"num_of_go_routines"`
}

// DebugHandlers different debug handlers
type DebugHandlers interface {
	DebugVars() http.Handler
	Stats() http.HandlerFunc
	DumpFunc() http.HandlerFunc
}

type debugHandlersDeps struct {
	fx.In

	Logger log.Logger
}

// InternalDebugHandlers defines internal debug handlers
//	- dump heap
//	- expvar
//	- running stats
func InternalDebugHandlers(deps debugHandlersDeps) []partial.HTTPHandlerPatternPair {
	return []partial.HTTPHandlerPatternPair{
		{Pattern: internalPatternPrefix + "/vars", Handler: deps.DebugVars()},
		{Pattern: internalPatternPrefix + "/dump", Handler: deps.DumpFunc()},
		{Pattern: internalPatternPrefix + "/stats", Handler: deps.Stats()},
	}
}

func (d *debugHandlersDeps) DebugVars() http.Handler {
	return expvar.Handler()
}

func (d *debugHandlersDeps) DumpFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		file, err := ioutil.TempFile("", "heapdump")
		if err != nil {
			d.Logger.WithError(err).Info(nil, "failed to create temp file to dump heap into")
			http.Error(w, "internal error, failed to serve heap dump", http.StatusInternalServerError)
			return
		}
		defer func(logger log.Logger, tempFile *os.File) {
			if err := os.Remove(tempFile.Name()); err != nil {
				logger.WithError(err).WithField("tempfile", tempFile.Name()).Warn(nil, "failed to remove temp file")
			}
		}(d.Logger, file) // remove garbage
		debug.WriteHeapDump(file.Fd())
		http.ServeFile(w, req, file.Name())
		if err = file.Close(); err != nil {
			d.Logger.WithError(err).WithField("tempfile", file.Name()).Warn(nil, "temp file wasn't closed")
		}
	}
}

func (d *debugHandlersDeps) Stats() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-type", "application/json; charset=utf-8")
		output := &StatsInfo{
			Memory:          new(runtime.MemStats),
			NumOfCPU:        runtime.NumCPU(),
			NumOfGoRoutines: runtime.NumGoroutine(),
		}
		runtime.ReadMemStats(output.Memory)
		if err := json.NewEncoder(w).Encode(output); err != nil {
			d.Logger.WithError(err).Debug(nil, "failed to serve stats")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
