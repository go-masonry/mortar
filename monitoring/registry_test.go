package monitoring

import (
	"testing"

	mock_monitor "github.com/go-masonry/mortar/interfaces/monitor/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRegistryCounterCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	externalMock := mock_monitor.NewMockBricksMetrics(ctrl)
	brickCounterMock := mock_monitor.NewMockBricksCounter(ctrl)
	registry := newRegistry(externalMock)
	externalMock.EXPECT().Counter(gomock.Any(), gomock.Any(), gomock.Any()).Return(brickCounterMock, nil)
	firstCounter, _ := registry.loadOrStoreCounter("name", "desc", "one", "two")
	secondCounter, _ := registry.loadOrStoreCounter("name", "desc", "one", "two")
	assert.Equal(t, firstCounter, secondCounter)
}

func TestRegistryGaugeCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	externalMock := mock_monitor.NewMockBricksMetrics(ctrl)
	brickGaugeMock := mock_monitor.NewMockBricksGauge(ctrl)
	registry := newRegistry(externalMock)
	externalMock.EXPECT().Gauge(gomock.Any(), gomock.Any(), gomock.Any()).Return(brickGaugeMock, nil)
	firstGauge, _ := registry.loadOrStoreGauge("name", "desc", "one", "two")
	secondGauge, _ := registry.loadOrStoreGauge("name", "desc", "one", "two")
	assert.Equal(t, firstGauge, secondGauge)
}

func TestRegistryTimerCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	externalMock := mock_monitor.NewMockBricksMetrics(ctrl)
	brickTimerMock := mock_monitor.NewMockBricksTimer(ctrl)
	registry := newRegistry(externalMock)
	externalMock.EXPECT().Timer(gomock.Any(), gomock.Any(), gomock.Any()).Return(brickTimerMock, nil)
	firstTimer, _ := registry.loadOrStoreTimer("name", "desc", "one", "two")
	secondTimer, _ := registry.loadOrStoreTimer("name", "desc", "one", "two")
	assert.Equal(t, firstTimer, secondTimer)
}

func TestRegistryHistogramCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	externalMock := mock_monitor.NewMockBricksMetrics(ctrl)
	brickHistogramMock := mock_monitor.NewMockBricksHistogram(ctrl)
	registry := newRegistry(externalMock)
	externalMock.EXPECT().Histogram(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(brickHistogramMock, nil)
	firstHistogram, _ := registry.loadOrStoreHistogram("name", "desc", nil, "one", "two")
	secondHistogram, _ := registry.loadOrStoreHistogram("name", "desc", nil, "one", "two")
	assert.Equal(t, firstHistogram, secondHistogram)
}

func TestIDGen(t *testing.T) {
	ID := calcID("name", "first", "second")
	assert.Equal(t, "name_first_second", ID)
}
