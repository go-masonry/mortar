package monitoring

import (
	"context"
	"fmt"
	"time"

	"github.com/go-masonry/mortar/interfaces/monitor"
	mock_monitor "github.com/go-masonry/mortar/interfaces/monitor/mock"
	"github.com/golang/mock/gomock"
)

func (m *reporterSuite) TestTimer() {
	bricksTimerMock := mock_monitor.NewMockBricksTimer(m.ctrl)
	timerMock := mock_monitor.NewMockTimer(m.ctrl)
	// Expect call to create timer
	m.bricksMetricsMocked.EXPECT().Timer("rate", "rate of something", []string{"one", "three"}).Return(bricksTimerMock, nil)
	// Expect call to update external timer with tags values
	bricksTimerMock.EXPECT().WithTags(gomock.Eq(map[string]string{"one": "11", "three": "33"})).Return(timerMock, nil)
	// Expect call to Record
	timerMock.EXPECT().Record(500 * time.Millisecond)

	// create everything and inc counter
	histogram := m.metrics.Timer("rate", "rate of something").WithTags(monitor.Tags{"one": "11"}).WithContext(context.Background())
	histogram.Record(500 * time.Millisecond)
}

func (m *reporterSuite) TestTimerFailures() {
	var errorCalledCounter = 0
	m.metrics = newMetric(m.registry, &monitorConfig{onError: func(error) {
		errorCalledCounter++
	}})
	// Failed to create Timer
	m.bricksMetricsMocked.EXPECT().Timer(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("no metric for you"))
	timer := m.metrics.Timer("failing", "metric")
	m.Equal(1, errorCalledCounter) // creation
	timer.Record(time.Second)
	m.Equal(2, errorCalledCounter)

	// Creating Histogram succeeded but every call to WithTags fails
	bricksTimerMock := mock_monitor.NewMockBricksTimer(m.ctrl)
	m.bricksMetricsMocked.EXPECT().Timer(gomock.Any(), gomock.Any(), gomock.Any()).Return(bricksTimerMock, nil)
	timer = m.metrics.Timer("success", "but not really")
	bricksTimerMock.EXPECT().WithTags(gomock.Len(0)).Return(nil, fmt.Errorf("different error"))
	timer.Record(time.Millisecond)
	m.Equal(3, errorCalledCounter)
}
