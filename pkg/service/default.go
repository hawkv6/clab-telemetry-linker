package service

import (
	"sync"

	"github.com/hawkv6/clab-telemetry-linker/pkg/config"
	"github.com/hawkv6/clab-telemetry-linker/pkg/consumer"
	"github.com/hawkv6/clab-telemetry-linker/pkg/logging"
	"github.com/hawkv6/clab-telemetry-linker/pkg/processor"
	"github.com/hawkv6/clab-telemetry-linker/pkg/publisher"
	"github.com/sirupsen/logrus"
)

type DefaultService struct {
	log       *logrus.Entry
	config    config.Config
	consumer  consumer.Consumer
	processor processor.Processor
	publisher publisher.Publisher
	wg        sync.WaitGroup
}

func NewDefaultService(config config.Config, receiver consumer.Consumer, processor processor.Processor, publisher publisher.Publisher) *DefaultService {
	return &DefaultService{
		log:       logging.DefaultLogger.WithField("subsystem", subsystem),
		config:    config,
		consumer:  receiver,
		processor: processor,
		publisher: publisher,
		wg:        sync.WaitGroup{},
	}
}
func (service *DefaultService) Start() {
	service.log.Infoln("Start all services")
	service.wg.Add(3)
	go func() {
		defer service.wg.Done()
		service.consumer.Start()
	}()
	go func() {
		defer service.wg.Done()
		service.processor.Start()
	}()
	go func() {
		defer service.wg.Done()
		service.publisher.Start()
	}()
}
func (service *DefaultService) Stop() {
	service.log.Infoln("Stopping all services")
	if err := service.consumer.Stop(); err != nil {
		service.log.Errorln("Error stopping consumer: ", err)
	}
	service.processor.Stop()
	if err := service.publisher.Stop(); err != nil {
		service.log.Errorln("Error stopping publisher: ", err)
	}
	service.wg.Wait()
}
