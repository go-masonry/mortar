package monitoring

import (
	"context"
	"fmt"

	"github.com/go-masonry/mortar/interfaces/monitor"
	mock_monitor "github.com/go-masonry/mortar/interfaces/monitor/mock"
	"github.com/golang/mock/gomock"
)

func (m *reporterSuite) TestHistogram() {
	bricksHistogramMock := mock_monitor.NewMockBricksHistogram(m.ctrl)
	histogramMock := mock_monitor.NewMockHistogram(m.ctrl)
	// Expect call to create histogram
	m.bricksMetricsMocked.EXPECT().Histogram("rate", "rate of something", []float64{0.1, 0.2}, []string{"one", "three"}).Return(bricksHistogramMock, nil)
	// Expect call to update external histogram with tags values
	bricksHistogramMock.EXPECT().WithTags(gomock.Eq(map[string]string{"one": "11", "three": "33"})).Return(histogramMock, nil)
	// Expect call to Record
	histogramMock.EXPECT().Record(0.5)

	// create everything and inc counter
	histogram := m.metrics.Histogram("rate", "rate of something", []float64{0.1, 0.2}).WithTags(monitor.Tags{"one": "11"}).WithContext(context.Background())
	histogram.Record(0.5)
}

func (m *reporterSuite) TestHistogramFailures() {
	var errorCalledCounter = 0
	m.metrics = newMetric(m.bricksMetricsMocked, &monitorConfig{onError: func(error) {
		errorCalledCounter++
	}})
	// Failed to create Histogram
	m.bricksMetricsMocked.EXPECT().Histogram(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("no metric for you"))
	histogram := m.metrics.Histogram("failing", "metric", []float64{0.1})
	m.Equal(1, errorCalledCounter) // creation
	histogram.Record(0.5)
	m.Equal(2, errorCalledCounter)

	// Creating Histogram succeeded but every call to WithTags fails
	bricksHistogramMock := mock_monitor.NewMockBricksHistogram(m.ctrl)
	m.bricksMetricsMocked.EXPECT().Histogram(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(bricksHistogramMock, nil)
	histogram = m.metrics.Histogram("success", "but not really", []float64{0.1})
	bricksHistogramMock.EXPECT().WithTags(gomock.Len(0)).Return(nil, fmt.Errorf("different error"))
	histogram.Record(0.6)
	m.Equal(3, errorCalledCounter)
}
