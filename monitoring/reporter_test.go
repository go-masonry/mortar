package monitoring

import (
	"context"
	"testing"

	"github.com/go-masonry/mortar/interfaces/monitor"
	mock_monitor "github.com/go-masonry/mortar/interfaces/monitor/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type reporterSuite struct {
	suite.Suite

	ctrl                *gomock.Controller
	bricksMetricsMocked *mock_monitor.MockBricksMetrics
	mockedReporter      *mock_monitor.MockBricksReporter

	reporter monitor.Reporter
	metrics  monitor.Metrics
}

func TestMetrics(t *testing.T) {
	suite.Run(t, new(reporterSuite))
}

func (m *reporterSuite) SetupTest() {
	m.ctrl = gomock.NewController(m.T())
	m.bricksMetricsMocked = mock_monitor.NewMockBricksMetrics(m.ctrl)
	m.mockedReporter = mock_monitor.NewMockBricksReporter(m.ctrl)
	m.mockedReporter.EXPECT().Metrics().Return(m.bricksMetricsMocked)
	m.reporter = newMortarReporter(&monitorConfig{
		reporter: m.mockedReporter,
		tags:     monitor.Tags{"one": "1", "three": "3"},
		extractors: []monitor.ContextExtractor{
			func(context.Context) monitor.Tags {
				return monitor.Tags{
					"three": "33",
				}
			},
		},
	})
	m.metrics = m.reporter.Metrics()
}

func (m *reporterSuite) TearDownTest() {
	m.ctrl.Finish()
}

func (m *reporterSuite) TestReporterConnectClose() {
	m.mockedReporter.EXPECT().Connect(gomock.Not(nil)).Return(nil)
	m.mockedReporter.EXPECT().Close(gomock.Not(nil)).Return(nil)
	m.reporter.Connect(context.Background())
	m.reporter.Close(context.Background())
}

func (m *reporterSuite) TestStaticTags() {
	bricksCounterMock := mock_monitor.NewMockBricksCounter(m.ctrl)
	counterMock := mock_monitor.NewMockCounter(m.ctrl)
	m.bricksMetricsMocked.EXPECT().Counter("rate", "rate of something", []string{"one", "three"}).Return(bricksCounterMock, nil)
	// Expect a call with static tag values
	bricksCounterMock.EXPECT().WithTags(gomock.Eq(map[string]string{"one": "1", "three": "3"})).Return(counterMock, nil)
	counterMock.EXPECT().Inc()
	// make sure counter is not created again
	m.metrics.Counter("rate", "rate of something").Inc()
}

func (m *reporterSuite) TestWithCustomTags() {
	bricksCounterMock := mock_monitor.NewMockBricksCounter(m.ctrl)
	counterMock := mock_monitor.NewMockCounter(m.ctrl)
	// Expect a call with static + custom tag values
	m.bricksMetricsMocked.EXPECT().Counter("additional", "tags", []string{"one", "three", "two"}).Return(bricksCounterMock, nil)
	counter := m.metrics.WithTags(monitor.Tags{"two": "2"}).Counter("additional", "tags")
	// expect one tag value was changed
	bricksCounterMock.EXPECT().WithTags(gomock.Eq(map[string]string{"one": "10", "two": "2", "three": "3"})).Return(counterMock, nil)
	counter = counter.WithTags(monitor.Tags{"one": "10"}) // overwrite tag value of "one"
	// Inc counter without tag changes
	counterMock.EXPECT().Inc()
	counter.Inc()
	// Again validate that tag changes are still there
	bricksCounterMock.EXPECT().WithTags(gomock.Eq(map[string]string{"one": "10", "two": "2", "three": "3"})).Return(counterMock, nil)
	counterMock.EXPECT().Inc()
	counter.Inc()
}
