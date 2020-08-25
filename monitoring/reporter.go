package monitoring

import (
	"context"

	"github.com/go-masonry/mortar/interfaces/monitor"
)

type mortarReporter struct {
	externalReporter monitor.BricksReporter
	cfg              *monitorConfig
}

// NewMortarReporter creates a new mortar monitoring reporter which is a wrapper to support
// 	- ContextExtractors
// 	- Default Tags, for example: {"version":"v1.0.1", "git":"22a85d8"}
//
// Meaning, it is possible to also extract tag values from the context, this is useful when the value is set per request/call within the context.Context:
// 	- Canary release https://martinfowler.com/bliki/CanaryRelease.html identifier
// 	- Authentication Token values, but avoid using high cardinality values such as UserID
//
func newMortarReporter(builder monitor.Builder, cfg *monitorConfig) monitor.Reporter {
	panic("implement me")
}

func (r *mortarReporter) Connect(ctx context.Context) error {
	return r.externalReporter.Connect(ctx)
}

func (r *mortarReporter) Close(ctx context.Context) error {
	return r.externalReporter.Close(ctx)
}

func (r *mortarReporter) Metrics() monitor.Metrics {
	panic("not implemented") // TODO: Implement
}
