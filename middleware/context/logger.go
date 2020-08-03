package context

import (
	"context"
	"github.com/go-masonry/mortar/interfaces/cfg"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/go-masonry/mortar/mortar"
	"go.uber.org/fx"
	"google.golang.org/grpc/metadata"
	"sort"
	"strings"
)

type loggerContextExtractorDeps struct {
	fx.In

	Config cfg.Config
}

func LoggerGRPCIncomingContextExtractor(deps loggerContextExtractorDeps) log.ContextExtractor {
	var includedHeaders []string
	if headers := deps.Config.Get(mortar.MiddlewareLoggerHeaders); headers.IsSet() {
		for _, header := range headers.StringSlice() {
			includedHeaders = append(includedHeaders, strings.ToLower(header))
		}
	}
	sort.Slice(includedHeaders, func(i, j int) bool {
		return len(includedHeaders[i]) < len(includedHeaders[j])
	})
	return headerPrefixes(includedHeaders).Extract
}

type headerPrefixes []string // if this slice will be very large it's better to build a trie map

func (h headerPrefixes) Extract(ctx context.Context) map[string]interface{} {
	var output = make(map[string]interface{})
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		for key, value := range md {
			lower := strings.ToLower(key)
			for _, prefix := range h {
				if strings.HasPrefix(lower, prefix) {
					output[lower] = strings.Join(value, ",")
				}
			}
		}
	}
	return output
}
