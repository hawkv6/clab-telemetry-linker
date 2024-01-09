package processor

import (
	"github.com/hawkv6/clab-telemetry-linker/pkg/consumer"
	"github.com/sirupsen/logrus"
)

var subsystem = "processor"

type Processor interface {
	Start()
}

type DefaultProcessor struct {
	log     *logrus.Entry
	msgChan chan consumer.Message
}

func NewDefaultProcessor(msgChan chan consumer.Message) *DefaultProcessor {
	return &DefaultProcessor{
		log:     logrus.WithField("subsystem", subsystem),
		msgChan: msgChan,
	}
}

func (processor *DefaultProcessor) Start() {
	for {
		msg := <-processor.msgChan
		processor.log.Infof("Received message: %v", msg)
	}
}
