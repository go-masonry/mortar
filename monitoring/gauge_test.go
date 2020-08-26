package monitoring

import (
	"context"
	"fmt"

	"github.com/go-masonry/mortar/interfaces/monitor"
	mock_monitor "github.com/go-masonry/mortar/interfaces/monitor/mock"
	"github.com/golang/mock/gomock"
)

func (m *reporterSuite) TestGauge() {
	mockedBricksGauge := mock_monitor.NewMockBricksGauge(m.ctrl)
	gaugeMock := mock_monitor.NewMockGauge(m.ctrl)
	m.bricksMetricsMocked.EXPECT().Gauge("gauge", "gauge desc", []string{"one", "three"}).Return(mockedBricksGauge, nil)

	mockedBricksGauge.EXPECT().WithTags(gomock.Eq(map[string]string{"one": "11", "three": "33"})).
		Return(gaugeMock, nil).
		Times(4)
	gauge := m.metrics.Gauge("gauge", "gauge desc").WithTags(monitor.Tags{"one": "11"}).WithContext(context.Background())
	gaugeMock.EXPECT().Add(1.1)
	gauge.Add(1.1)
	gaugeMock.EXPECT().Inc()
	gauge.Inc()
	gaugeMock.EXPECT().Dec()
	gauge.Dec()
	gaugeMock.EXPECT().Set(5.0)
	gauge.Set(5.0)
}

func (m *reporterSuite) TestGaugeFailures() {
	var errorCalledCounter = 0
	m.metrics = newMetric(m.registry, &monitorConfig{onError: func(error) {
		errorCalledCounter++
	}})
	// Failed to create Counter
	m.bricksMetricsMocked.EXPECT().Gauge(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("no metric for you"))
	gauge := m.metrics.Gauge("failing", "metric")
	m.Equal(1, errorCalledCounter)
	gauge.Inc()
	m.Equal(2, errorCalledCounter)
	gauge.Add(1.1)
	m.Equal(3, errorCalledCounter)
	gauge.WithTags(monitor.Tags{"irrelevant": "tag"}).Inc()
	m.Equal(4, errorCalledCounter)
	gauge.WithContext(context.Background()).Inc()
	m.Equal(5, errorCalledCounter)
	gauge.Dec()
	m.Equal(6, errorCalledCounter)
	gauge.Set(2.0)
	m.Equal(7, errorCalledCounter)
	// Creating gauge succeeded but every call to WithTags/WithContext fails
	bricksGaugeMock := mock_monitor.NewMockBricksGauge(m.ctrl)
	m.bricksMetricsMocked.EXPECT().Gauge(gomock.Any(), gomock.Any(), gomock.Any()).Return(bricksGaugeMock, nil)
	gauge = m.metrics.Gauge("success", "but not really")
	// Inc
	bricksGaugeMock.EXPECT().WithTags(gomock.Len(0)).
		Return(nil, fmt.Errorf("different error")).
		Times(4)
	gauge.Inc()
	m.Equal(8, errorCalledCounter)
	gauge.Add(2.2)
	m.Equal(9, errorCalledCounter)
	gauge.Dec()
	m.Equal(10, errorCalledCounter)
	gauge.Set(2.2)
	m.Equal(11, errorCalledCounter)
}
