package monitoring

import (
	"fmt"
	"time"

	"github.com/go-masonry/mortar/interfaces/monitor"
)

type noop struct {
	name, desc string
	err        error
	onError    func(error)
}

type noopCounter struct {
	*noop
}

func (n *noopCounter) WithTags(tags map[string]string) (monitor.Counter, error) {
	return n, nil
}

func newNoopCounter(err error, onError func(error)) monitor.BricksCounter {
	return &noopCounter{&noop{
		err:     err,
		onError: onError,
	}}
}

type noopGauge struct {
	*noop
}

func (n *noopGauge) WithTags(tags map[string]string) (monitor.Gauge, error) {
	return n, nil
}

func newNoopGauge(err error, onError func(error)) monitor.BricksGauge {
	return &noopGauge{&noop{
		err:     err,
		onError: onError,
	}}
}

type noopHistogram struct {
	*noop
}

func (n *noopHistogram) WithTags(tags map[string]string) (monitor.Histogram, error) {
	return n, nil
}

func newNoopHistogram(err error, onError func(error)) monitor.BricksHistogram {
	return &noopHistogram{&noop{
		err:     err,
		onError: onError,
	}}
}

type noopTimer struct {
	*noop
}

func (n *noopTimer) WithTags(tags map[string]string) (monitor.Timer, error) {
	return n, nil
}

func (n *noopTimer) Record(d time.Duration) {
	n.noop.Record(d.Seconds())
}

func newNoopTimer(err error, onError func(error)) monitor.BricksTimer {
	return &noopTimer{&noop{
		err:     err,
		onError: onError,
	}}
}

// Inc increments the counter by 1
func (n *noop) Inc() {
	n.do()
}

// Add adds the given value to the counter, negative values are not advised
func (n *noop) Add(v float64) {
	n.do()
}

// Record value
func (n *noop) Record(v float64) {
	n.do()
}

// Set sets Gauge value
func (n *noop) Set(v float64) {
	n.do()
}

// Dec adds -1
func (n *noop) Dec() {
	n.do()
}

func (n *noop) do() {
	n.onError(
		fmt.Errorf("still trying to use failed metric %s:%s, %w", n.name, n.desc, n.err),
	)
}
