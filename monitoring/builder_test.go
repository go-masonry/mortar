package monitoring

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-masonry/mortar/interfaces/monitor"
	mock_monitor "github.com/go-masonry/mortar/interfaces/monitor/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestBuilder(t *testing.T) {
	var onErrorCalled bool
	builder := Builder().
		AddExtractors(func(ctx context.Context) monitor.Tags {
			return monitor.Tags{
				"one": "1",
				"two": "2",
			}
		}).
		DoOnError(func(err error) {
			onErrorCalled = assert.Error(t, err)
		}).
		SetTags(monitor.Tags{
			"always": "there",
		})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBuilder := mock_monitor.NewMockBuilder(ctrl)
	mockReporter := mock_monitor.NewMockBricksReporter(ctrl)
	mockReporter.EXPECT().Metrics().Return(nil)
	mockBuilder.EXPECT().Build().Return(mockReporter)
	reporter := builder.Build(mockBuilder).(*mortarReporter)
	reporter.cfg.onError(fmt.Errorf("some error"))
	// Assertions
	assert.True(t, onErrorCalled)
	assert.Len(t, reporter.cfg.extractors, 1)
	assert.Contains(t, reporter.cfg.tags, "always")
}

func TestBuilderDefaults(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBuilder := mock_monitor.NewMockBuilder(ctrl)
	mockReporter := mock_monitor.NewMockBricksReporter(ctrl)
	mockReporter.EXPECT().Metrics().Return(nil)
	mockBuilder.EXPECT().Build().Return(mockReporter)
	reporter := Builder().Build(mockBuilder).(*mortarReporter)

	assert.NotNil(t, reporter.cfg.onError)
}
