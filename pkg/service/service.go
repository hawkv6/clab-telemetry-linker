package service

import (
	"github.com/hawkv6/clab-telemetry-linker/pkg/config"
	"github.com/hawkv6/clab-telemetry-linker/pkg/consumer"
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
		log:       logrus.WithField("subsystem", subsystem),
		config:    config,
		consumer:  receiver,
		processor: processor,
		publisher: publisher,
	}
}
func (service *DefaultService) Start() error {
	go service.consumer.Start()
	go service.processor.Start()
	service.log.Infoln("Start Service")
	return nil
}
func (service *DefaultService) Stop() error {
	service.consumer.Stop()
	service.log.Infoln("Stop Service")
	return nil
}
