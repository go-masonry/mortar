package monitoring

import (
	"fmt"
	"strconv"
	"sync"
	"testing"

	"github.com/go-masonry/mortar/interfaces/monitor"
	mock_monitor "github.com/go-masonry/mortar/interfaces/monitor/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	sameID := calcID("name", "second", "first")
	assert.Equal(t, ID, sameID)
}

func TestRegistryReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	externalMock := mock_monitor.NewMockBricksMetrics(ctrl)
	registry := newRegistry(externalMock)
	externalMock.EXPECT().Histogram(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("constant error")).Times(2)
	_, err1 := registry.loadOrStoreHistogram("name", "desc", nil, "one", "two")
	_, err2 := registry.loadOrStoreHistogram("name", "desc", nil, "one", "two")
	assert.EqualError(t, err1, "constant error")
	assert.EqualError(t, err2, "constant error")
	externalMock.EXPECT().Timer(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("constant error")).Times(2)
	_, err1 = registry.loadOrStoreTimer("name", "desc", "one", "two")
	_, err2 = registry.loadOrStoreTimer("name", "desc", "one", "two")
	assert.EqualError(t, err1, "constant error")
	assert.EqualError(t, err2, "constant error")
	externalMock.EXPECT().Counter(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("constant error")).Times(2)
	_, err1 = registry.loadOrStoreCounter("name", "desc", "one", "two")
	_, err2 = registry.loadOrStoreCounter("name", "desc", "one", "two")
	assert.EqualError(t, err1, "constant error")
	assert.EqualError(t, err2, "constant error")
	externalMock.EXPECT().Gauge(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("constant error")).Times(2)
	_, err1 = registry.loadOrStoreGauge("name", "desc", "one", "two")
	_, err2 = registry.loadOrStoreGauge("name", "desc", "one", "two")
	assert.EqualError(t, err1, "constant error")
	assert.EqualError(t, err2, "constant error")
}

func TestConcurrentGaugeCreation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	externalMock := mock_monitor.NewMockBricksMetrics(ctrl)
	brickGaugeMock := mock_monitor.NewMockBricksGauge(ctrl)
	registry := newRegistry(externalMock)
	for i := 0; i < 10; i++ {
		var created bool
		var mu sync.Mutex
		name := strconv.Itoa(i)
		var wg sync.WaitGroup
		var err1, err2 error
		externalMock.EXPECT().Gauge(name, gomock.Any(), gomock.Any()).DoAndReturn(func(name string, _ string, _ ...string) (monitor.BricksGauge, error) {
			mu.Lock()
			defer mu.Unlock()
			if created {
				return nil, fmt.Errorf("already exists")
			}
			created = true
			return brickGaugeMock, nil
		}).MaxTimes(2).MinTimes(1)
		wg.Add(2)
		// concurrent run 1
		go func(n string) {
			defer wg.Done()
			_, err1 = registry.loadOrStoreGauge(n, "")
		}(name)

		// concurrent run 2
		go func(n string) {
			defer wg.Done()
			_, err2 = registry.loadOrStoreGauge(n, "")
		}(name)
		wg.Wait()
		require.NoError(t, err1, "[%d] first error", i)
		require.NoError(t, err2, "[%d] second error", i)
	}
}

func TestConcurrentTimerCreation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	externalMock := mock_monitor.NewMockBricksMetrics(ctrl)
	brickTimerMock := mock_monitor.NewMockBricksTimer(ctrl)
	registry := newRegistry(externalMock)
	for i := 0; i < 10; i++ {
		var created bool
		var mu sync.Mutex
		name := strconv.Itoa(i)
		var wg sync.WaitGroup
		var err1, err2 error
		externalMock.EXPECT().Timer(name, gomock.Any(), gomock.Any()).DoAndReturn(func(name string, _ string, _ ...string) (monitor.BricksTimer, error) {
			mu.Lock()
			defer mu.Unlock()
			if created {
				return nil, fmt.Errorf("already exists: %s", name)
			}
			created = true
			return brickTimerMock, nil
		}).MaxTimes(2).MinTimes(1)
		wg.Add(2)
		// concurrent run 1
		go func(n string) {
			defer wg.Done()
			_, err1 = registry.loadOrStoreTimer(n, "")
		}(name)

		// concurrent run 2
		go func(n string) {
			defer wg.Done()
			_, err2 = registry.loadOrStoreTimer(n, "")
		}(name)
		wg.Wait()
		require.NoError(t, err1, "[%d] first error", i)
		require.NoError(t, err2, "[%d] second error", i)
	}
}

func TestConcurrentCounterCreation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	externalMock := mock_monitor.NewMockBricksMetrics(ctrl)
	brickCounterMock := mock_monitor.NewMockBricksCounter(ctrl)
	registry := newRegistry(externalMock)
	for i := 0; i < 10; i++ {
		var created bool
		var mu sync.Mutex
		name := strconv.Itoa(i)
		var wg sync.WaitGroup
		var err1, err2 error
		externalMock.EXPECT().Counter(name, gomock.Any(), gomock.Any()).DoAndReturn(func(name string, _ string, _ ...string) (monitor.BricksCounter, error) {
			mu.Lock()
			defer mu.Unlock()
			if created {
				return nil, fmt.Errorf("already exists")
			}
			created = true
			return brickCounterMock, nil
		}).MaxTimes(2).MinTimes(1)
		wg.Add(2)
		// concurrent run 1
		go func(n string) {
			defer wg.Done()
			_, err1 = registry.loadOrStoreCounter(n, "")
		}(name)

		// concurrent run 2
		go func(n string) {
			defer wg.Done()
			_, err2 = registry.loadOrStoreCounter(n, "")
		}(name)
		wg.Wait()
		require.NoError(t, err1, "[%d] first error", i)
		require.NoError(t, err2, "[%d] second error", i)
	}
}

func TestConcurrentHistogramCreation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	externalMock := mock_monitor.NewMockBricksMetrics(ctrl)
	brickHistogramMock := mock_monitor.NewMockBricksHistogram(ctrl)
	registry := newRegistry(externalMock)
	for i := 0; i < 10; i++ {
		var created bool
		var mu sync.Mutex
		name := strconv.Itoa(i)
		var wg sync.WaitGroup
		var err1, err2 error
		externalMock.EXPECT().Histogram(name, gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(name string, _ string, _ []float64, _ ...string) (monitor.BricksHistogram, error) {
			mu.Lock()
			defer mu.Unlock()
			if created {
				return nil, fmt.Errorf("already exists")
			}
			created = true
			return brickHistogramMock, nil
		}).MaxTimes(2).MinTimes(1)
		wg.Add(2)
		// concurrent run 1
		go func(n string) {
			defer wg.Done()
			_, err1 = registry.loadOrStoreHistogram(n, "", nil)
		}(name)

		// concurrent run 2
		go func(n string) {
			defer wg.Done()
			_, err2 = registry.loadOrStoreHistogram(n, "", nil)
		}(name)
		wg.Wait()
		require.NoError(t, err1, "[%d] first error", i)
		require.NoError(t, err2, "[%d] second error", i)
	}
}
