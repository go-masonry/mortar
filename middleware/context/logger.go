package context

import (
	"context"
	"sort"
	"strings"

	"github.com/go-masonry/mortar/interfaces/cfg"
	confkeys "github.com/go-masonry/mortar/interfaces/cfg/keys"
	"github.com/go-masonry/mortar/interfaces/log"
	"go.uber.org/fx"
	"google.golang.org/grpc/metadata"
)

type loggerContextExtractorDeps struct {
	fx.In

	Config cfg.Config
}

// LoggerGRPCIncomingContextExtractor creates a context extractor for logger
// This is useful if you want to add different fields from gRPC incoming metadata.MD to a log entry
func LoggerGRPCIncomingContextExtractor(deps loggerContextExtractorDeps) log.ContextExtractor {
	var includedHeaders []string
	if headers := deps.Config.Get(confkeys.LoggerIncomingGRPCMetadataHeadersExtractor); headers.IsSet() {
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
