package processor

import (
	"regexp"
	"strconv"

	"github.com/hawkv6/clab-telemetry-linker/pkg/config"
	"github.com/hawkv6/clab-telemetry-linker/pkg/consumer"
	"github.com/hawkv6/clab-telemetry-linker/pkg/helpers"
	"github.com/hawkv6/clab-telemetry-linker/pkg/logging"
	"github.com/sirupsen/logrus"
)

var subsystem = "processor"

type Processor interface {
	Start()
	Stop()
}

type DefaultProcessor struct {
	log                *logrus.Entry
	config             config.Config
	unprocessedMsgChan chan consumer.Message
	processedMsgChan   chan consumer.Message
	quitChan           chan bool
	helper             helpers.Helper
}

func NewDefaultProcessor(config config.Config, unprocessedMsgChan chan consumer.Message, processedMsgChan chan consumer.Message, helper helpers.Helper) *DefaultProcessor {
	return &DefaultProcessor{
		log:                logging.DefaultLogger.WithField("subsystem", subsystem),
		config:             config,
		unprocessedMsgChan: unprocessedMsgChan,
		processedMsgChan:   processedMsgChan,
		quitChan:           make(chan bool),
		helper:             helper,
	}
}

func shortenInterfaceName(name string) string {
	re := regexp.MustCompile(`GigabitEthernet(\d+)/(\d+)/(\d+)/(\d+)`)
	return re.ReplaceAllString(name, "Gi$1-$2-$3-$4")
}

func (processor *DefaultProcessor) processDelayMessage(msg *consumer.DelayMessage) {
	processor.log.Debugf("Process delay of node %s of interface %s", msg.Tags.Source, msg.Tags.InterfaceName)
	shortInterfaceName := shortenInterfaceName(msg.Tags.InterfaceName)
	impairmentsPrefix := processor.helper.GetDefaultImpairmentsPrefix(msg.Tags.Source, shortInterfaceName)
	delay := processor.config.GetValue(impairmentsPrefix + "delay")
	if delay != "" {
		delayValue, err := strconv.ParseFloat(delay, 64)
		delayValueUsec := delayValue * 1000
		if err != nil {
			processor.log.Errorf("Failed to convert delay to float64: %v", err)
			return
		}
		msg.Average = delayValueUsec
		msg.Maximum = delayValueUsec
		msg.Minimum = delayValueUsec
		processor.log.Debugf("Adjusted delay of node %s of interface %s to: %f", msg.Tags.Source, msg.Tags.InterfaceName, delayValueUsec)
	}
	processor.processedMsgChan <- msg
}

func (processor *DefaultProcessor) processLossMessage(msg *consumer.LossMessage) {
	processor.log.Debugf("Process loss of node %s of interface %s", msg.Tags.Source, msg.Tags.InterfaceName)
	shortInterfaceName := shortenInterfaceName(msg.Tags.InterfaceName)
	impairmentsPrefix := processor.helper.GetDefaultImpairmentsPrefix(msg.Tags.Source, shortInterfaceName)
	loss := processor.config.GetValue(impairmentsPrefix + "loss")
	if loss != "" {
		lossValue, err := strconv.ParseFloat(loss, 64)
		if err != nil {
			processor.log.Errorf("Failed to convert loss to float64: %v", err)
			return
		}
		msg.LossPercentage = lossValue
		processor.log.Debugf("Adjusted loss of node %s of interface %s to: %f", msg.Tags.Source, msg.Tags.InterfaceName, lossValue)
	}
	processor.processedMsgChan <- msg
}

func (processor *DefaultProcessor) processBandwidthMessage(msg *consumer.BandwidthMessage) {
	processor.log.Debugf("Process bandwidth of node %s of interface %s", msg.Tags.Source, msg.Tags.InterfaceName)
	shortInterfaceName := shortenInterfaceName(msg.Tags.InterfaceName)
	impairmentsPrefix := processor.helper.GetDefaultImpairmentsPrefix(msg.Tags.Source, shortInterfaceName)
	bandwidth := processor.config.GetValue(impairmentsPrefix + "rate")
	if bandwidth != "" {
		bandwidthValue, err := strconv.ParseFloat(bandwidth, 64)
		if err != nil {
			processor.log.Errorf("Failed to convert bandwidth to float64: %v", err)
			return
		}
		msg.Bandwidth = bandwidthValue
		processor.log.Debugf("Adjusted bandwidth of node %s of interface %s to: %f", msg.Tags.Source, msg.Tags.InterfaceName, bandwidthValue)
	}
	processor.processedMsgChan <- msg
}

func (processor *DefaultProcessor) processMessage(msg consumer.Message) {
	switch msg := msg.(type) {
	case *consumer.DelayMessage:
		processor.processDelayMessage(msg)
	case *consumer.LossMessage:
		processor.processLossMessage(msg)
	case *consumer.BandwidthMessage:
		processor.processBandwidthMessage(msg)
	default:
		processor.log.Errorf("Skipping unknown message type: %v", msg)
	}
}

func (processor *DefaultProcessor) Start() {
	processor.log.Infoln("Starting processing messages")
	for {
		select {
		case msg := <-processor.unprocessedMsgChan:
			processor.processMessage(msg)
		case <-processor.quitChan:
			processor.log.Debug("Stopping processor")
			return
		}
	}
}

func (processor *DefaultProcessor) Stop() {
	processor.quitChan <- true
}
