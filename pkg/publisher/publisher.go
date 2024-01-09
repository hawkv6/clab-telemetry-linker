package publisher

import (
	"github.com/hawkv6/clab-telemetry-linker/pkg/consumer"
	"github.com/sirupsen/logrus"
)

var subsystem = "publisher"

type Publisher interface {
	Start() error
}

type DefaultPublisher struct {
	log              *logrus.Entry
	processedMsgChan chan consumer.Message
}

func NewDefaultPublisher(msgChan chan consumer.Message) *DefaultPublisher {
	return &DefaultPublisher{
		log:              logrus.WithField("subsystem", subsystem),
		processedMsgChan: msgChan,
	}
}

func (publisher *DefaultPublisher) Start() error {
	for {
		msg := <-publisher.processedMsgChan
		publisher.log.Infof("Received message: %v", msg)
	}
	// publisher.log.Infoln("Start Publisher")
	// return nil
}
