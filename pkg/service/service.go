package service

import (
	"github.com/hawkv6/clab-telemetry-linker/pkg/config"
	"github.com/hawkv6/clab-telemetry-linker/pkg/consumer"
	"github.com/hawkv6/clab-telemetry-linker/pkg/logging"
	"github.com/hawkv6/clab-telemetry-linker/pkg/processor"
	"github.com/hawkv6/clab-telemetry-linker/pkg/publisher"
	"github.com/sirupsen/logrus"
)

// Read Config and watch for changes
// Create Kafka Broker, Consumer, Topic
// Start listening for messages (own go routine)
// Start processing messages (own go routine)
// Start writing processed messages back to own topic (own go routine)

var subsystem = "service"

type Service interface {
	Start() error
	Stop() error
}

type DefaultService struct {
	log       *logrus.Entry
	config    config.Config
	consumer  consumer.Consumer
	processor processor.Processor
	publisher publisher.Publisher
}

func NewDefaultService(config config.Config, receiver consumer.Consumer, processor processor.Processor, publisher publisher.Publisher) *DefaultService {
	return &DefaultService{
		log:       logging.DefaultLogger.WithField("subsystem", subsystem),
		config:    config,
		consumer:  receiver,
		processor: processor,
		publisher: publisher,
	}
}
func (service *DefaultService) Start() error {
	go service.consumer.Start()
	go service.processor.Start()
	go service.publisher.Start()
	service.log.Infoln("Start all services")
	return nil
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
}
