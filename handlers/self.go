package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-masonry/mortar/constructors/partial"
	"github.com/go-masonry/mortar/interfaces/cfg"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/go-masonry/mortar/mortar"
	"github.com/go-masonry/mortar/utils"
	"go.uber.org/fx"
	"net/http"
	"os"
	"strings"
)

const (
	ObfuscationEdgeLength = 4
	selfHandlerPrefix     = "/self"
)

type selfHandlerDeps struct {
	fx.In

	Logger log.Logger
	Config cfg.Config
}

type InternalHandlers interface {
	BuildInfo() http.HandlerFunc
	ConfigMap() http.HandlerFunc
}

func SelfHandlersOption() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group:  partial.FxGroupInternalHttpHandlers + ",flatten",
			Target: SelfHandlers,
		})
}

func SelfHandlers(deps selfHandlerDeps) []partial.HttpHandlerPatternPair {
	return []partial.HttpHandlerPatternPair{
		{Pattern: selfHandlerPrefix + "/build", Handler: deps.BuildInfo()},
		{Pattern: selfHandlerPrefix + "/config", Handler: deps.ConfigMap()},
	}
}

func (s *selfHandlerDeps) BuildInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		information := mortar.GetBuildInformation(true)
		if err := json.NewEncoder(w).Encode(information); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s.Logger.WithError(err).Warn(nil, "failed to serve build info")
		}
	}
}

func (s *selfHandlerDeps) ConfigMap() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-type", "application/json; charset=utf-8")
		output := make(map[string]interface{})
		output["config"] = s.getConfigVariables()
		output["environment"] = s.getEnvVariables()
		if err := json.NewEncoder(w).Encode(output); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s.Logger.WithError(err).Warn(nil, "failed to server config map")
		}
	}
}

func (s *selfHandlerDeps) getConfigVariables() map[string]interface{} {
	return s.obfuscateMapWhereNeeded("", s.Config.Map())
}

func (s *selfHandlerDeps) obfuscateMapWhereNeeded(prefix string, confMap map[string]interface{}) map[string]interface{} {
	output := make(map[string]interface{})
	for k, v := range confMap {
		obfuscateKey := fmt.Sprintf("%s.%s", prefix, k)
		if mValue, ok := v.(map[string]interface{}); ok {
			output[k] = s.obfuscateMapWhereNeeded(obfuscateKey, mValue)
		} else {
			output[k] = s.obfuscateIfNeeded(obfuscateKey, v)
		}
	}
	return output
}

func (s *selfHandlerDeps) getEnvVariables() map[string]string {
	output := make(map[string]string)
	for _, keyValue := range os.Environ() {
		if keyValueSlice := strings.Split(keyValue, "="); len(keyValueSlice) > 1 {
			key := keyValueSlice[0]
			output[key] = s.obfuscateIfNeeded(key, keyValueSlice[1])
		}
	}
	return output
}

func (s *selfHandlerDeps) obfuscateIfNeeded(key string, value interface{}) string {
	var valueAsString string
	if value == nil {
		return ""
	}
	switch v := value.(type) {
	case string, fmt.Stringer:
		valueAsString = fmt.Sprintf("%s", v)
	default:
		valueAsString = fmt.Sprintf("%v", v)
	}
	hideKeys := s.Config.Get(mortar.HandlersSelfObfuscateConfigKeys).StringSlice() // if none exist slice will be empty
	for _, hidePart := range hideKeys {
		if strings.Contains(strings.ToLower(key), strings.ToLower(hidePart)) {
			return utils.Obfuscate(valueAsString, ObfuscationEdgeLength)
		}
	}
	return valueAsString
}
