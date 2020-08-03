package providers

import (
	"github.com/go-masonry/mortar/constructors"
	"go.uber.org/fx"
)

// JWTExtractorFxOption adds default JWT extractor from context.Context to the graph
func JWTExtractorFxOption() fx.Option {
	return fx.Provide(constructors.DefaultJWTTokenExtractor)
}
