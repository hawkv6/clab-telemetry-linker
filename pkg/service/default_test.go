package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/hawkv6/clab-telemetry-linker/pkg/config"
	"github.com/hawkv6/clab-telemetry-linker/pkg/consumer"
	"github.com/hawkv6/clab-telemetry-linker/pkg/processor"
	"github.com/hawkv6/clab-telemetry-linker/pkg/publisher"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewDefaultService(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Test Creating Default Service",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			config := config.NewMockConfig(ctrl)
			consumer := consumer.NewMockConsumer(ctrl)
			processor := processor.NewMockProcessor(ctrl)
			publisher := publisher.NewMockPublisher(ctrl)
			defaultService := NewDefaultService(config, consumer, processor, publisher)
			assert.NotNil(t, defaultService)
		})
	}
}

func TestDefaultService_Start(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Test Starting Default Service",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			config := config.NewMockConfig(ctrl)
			consumer := consumer.NewMockConsumer(ctrl)
			processor := processor.NewMockProcessor(ctrl)
			publisher := publisher.NewMockPublisher(ctrl)
			consumer.EXPECT().Start().Return()
			processor.EXPECT().Start().Return()
			publisher.EXPECT().Start().Return()
			defaultService := NewDefaultService(config, consumer, processor, publisher)
			assert.NotPanics(t, func() {
				defaultService.Start()
				time.Sleep(3 * time.Second)
			})
		})
	}
}

func TestDefaultService_Stop(t *testing.T) {
	tests := []struct {
		name       string
		wantsError bool
	}{
		{
			name:       "Test Stopping Default Service without error",
			wantsError: false,
		},
		{
			name:       "Test Stopping Default Service with error in consumer",
			wantsError: true,
		},
		{
			name:       "Test Stopping Default Service with error in producer",
			wantsError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			config := config.NewMockConfig(ctrl)
			consumer := consumer.NewMockConsumer(ctrl)
			processor := processor.NewMockProcessor(ctrl)
			publisher := publisher.NewMockPublisher(ctrl)
			defaultService := NewDefaultService(config, consumer, processor, publisher)
			processor.EXPECT().Stop().Return()
			if tt.wantsError {
				consumer.EXPECT().Stop().Return(fmt.Errorf("error stopping consumer"))
				publisher.EXPECT().Stop().Return(fmt.Errorf("error error stopping publisher"))
			} else {
				consumer.EXPECT().Stop().Return(nil)
				publisher.EXPECT().Stop().Return(nil)
			}
			defaultService.Stop()
		})
	}
}
