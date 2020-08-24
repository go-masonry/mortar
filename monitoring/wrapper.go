package monitoring

import (
	"github.com/go-masonry/mortar/interfaces/monitor"
)

// NewMortarReporter creates a new mortar monitoring reporter which is a wrapper to support
// 	- ContextExtractors
// 	- Default Tags, for example: {"version":"v1.0.1", "git":"22a85d8"}
//
// Meaning, it is possible to also extract tag values from the context, this is useful when the value is set per request/call within the context.Context:
// 	- Canary release https://martinfowler.com/bliki/CanaryRelease.html identifier
// 	- Authentication Token values, but avoid using high cardinality values such as UserID
//
func NewMortarReporter(builder monitor.Builder, tags monitor.Tags, contextExtractors ...monitor.ContextExtractor) monitor.Reporter {
	panic("implement me")
}
