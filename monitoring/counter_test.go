package monitoring

import (
	"context"
	"fmt"

	"github.com/go-masonry/mortar/interfaces/monitor"
	mock_monitor "github.com/go-masonry/mortar/interfaces/monitor/mock"
	"github.com/golang/mock/gomock"
)

func (m *reporterSuite) TestCounter() {
	bricksCounterMock := mock_monitor.NewMockBricksCounter(m.ctrl)
	counterMock := mock_monitor.NewMockCounter(m.ctrl)
	// Expect call to create counter
	m.bricksMetricsMocked.EXPECT().Counter("rate", "rate of something", []string{"one", "three"}).Return(bricksCounterMock, nil)
	// Expect call to update external counter with tags values
	bricksCounterMock.EXPECT().WithTags(gomock.Eq(map[string]string{"one": "11", "three": "33"})).Return(counterMock, nil).Times(2)
	// Expect call to Inc()
	counterMock.EXPECT().Inc()
	counterMock.EXPECT().Add(1.1)

	// create everything and inc counter
	counter := m.metrics.Counter("rate", "rate of something").WithTags(monitor.Tags{"one": "11"}).WithContext(context.Background())
	counter.Inc()
	counter.Add(1.1)
}

func (m *reporterSuite) TestCounterFailures() {
	var errorCalledCounter = 0
	m.metrics = newMetric(m.bricksMetricsMocked, &monitorConfig{onError: func(error) {
		errorCalledCounter++
	}})
	// Failed to create Counter
	m.bricksMetricsMocked.EXPECT().Counter(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("no metric for you"))
	counter := m.metrics.Counter("failing", "metric")
	m.Equal(1, errorCalledCounter)
	counter.Inc()
	m.Equal(2, errorCalledCounter)
	counter.Add(1.1)
	m.Equal(3, errorCalledCounter)
	counter.WithTags(monitor.Tags{"irrelevant": "tag"}).Inc()
	m.Equal(4, errorCalledCounter)
	counter.WithContext(context.Background()).Inc()
	m.Equal(5, errorCalledCounter)
	// Creating counter succeeded but every call to WithTags fails
	bricksCounterMock := mock_monitor.NewMockBricksCounter(m.ctrl)
	m.bricksMetricsMocked.EXPECT().Counter(gomock.Any(), gomock.Any(), gomock.Any()).Return(bricksCounterMock, nil)
	counter = m.metrics.Counter("success", "but not really")
	// Inc
	bricksCounterMock.EXPECT().WithTags(gomock.Len(0)).Return(nil, fmt.Errorf("different error"))
	counter.Inc()
	m.Equal(6, errorCalledCounter)
	// Add
	bricksCounterMock.EXPECT().WithTags(gomock.Len(0)).Return(nil, fmt.Errorf("different error"))
	counter.Add(2.2)
	m.Equal(7, errorCalledCounter)
}
