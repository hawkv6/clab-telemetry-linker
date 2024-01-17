package processor

import (
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"github.com/hawkv6/clab-telemetry-linker/pkg/config"
	"github.com/hawkv6/clab-telemetry-linker/pkg/consumer"
	"github.com/hawkv6/clab-telemetry-linker/pkg/helpers"
	"github.com/hawkv6/clab-telemetry-linker/pkg/logging"
	"github.com/sirupsen/logrus"
)

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

func (processor *DefaultProcessor) shortenInterfaceName(name string) (error, string) {
	re := regexp.MustCompile(`GigabitEthernet(\d+)/(\d+)/(\d+)/(\d+)`)
	if !re.MatchString(name) {
		return fmt.Errorf("interface name %s does not match expected pattern", name), ""
	}
	return nil, re.ReplaceAllString(name, "Gi$1-$2-$3-$4")
}

func (processor *DefaultProcessor) getDelayValues(impairmentsPrefix string) (error, float64, float64) {
	delay := processor.config.GetValue(impairmentsPrefix + "delay")
	delayValueUsec := 0.0
	jitterValueUsec := 0.0
	if delay != "" {
		delayValue, err := strconv.ParseFloat(delay, 64)
		delayValueUsec = delayValue * 1000
		if err != nil {
			return fmt.Errorf("Failed to convert delay to float64: %v", err), 0, 0
		}
		jitter := processor.config.GetValue(impairmentsPrefix + "jitter")
		if jitter != "" {
			jitterValue, err := strconv.ParseFloat(jitter, 64)
			if err != nil {
				return fmt.Errorf("Failed to convert jitter to float64: %v", err), 0, 0
			}
			jitterValueUsec = jitterValue * 1000
		}
	}
	return nil, delayValueUsec, jitterValueUsec
}

func (processor *DefaultProcessor) setDelayValues(msg *consumer.DelayMessage, delay float64, jitter float64, randomFactor float64) {
	if delay == 0.0 {
		msg.Average = msg.Average + msg.Average*randomFactor
	} else {
		msg.Average = msg.Average + (delay + delay*randomFactor)
	}
	if jitter == 0.0 {
		msg.Maximum = msg.Average + msg.Average*randomFactor
		msg.Minimum = msg.Average - msg.Average*randomFactor
		msg.Variance = msg.Maximum - msg.Minimum
	} else {
		msg.Maximum = msg.Average + 0.5*jitter
		msg.Minimum = msg.Average - 0.5*jitter
		msg.Variance = jitter
	}
}

func (processor *DefaultProcessor) processDelayMessage(msg *consumer.DelayMessage) {
	processor.log.Debugf("Process delay of node %s of interface %s", msg.Tags.Source, msg.Tags.InterfaceName)
	err, shortInterfaceName := processor.shortenInterfaceName(msg.Tags.InterfaceName)
	if err != nil {
		processor.log.Debugf("Failed to shorten interface name: %v", err)
		return
	}
	impairmentsPrefix := processor.helper.GetDefaultImpairmentsPrefix(msg.Tags.Source, shortInterfaceName)
	err, delay, jitter := processor.getDelayValues(impairmentsPrefix)
	if err != nil {
		processor.log.Errorf("Failed to get delay values: %v", err)
		return
	}

	// Add normalized random factor to the delay
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomFactor := (math.Log10(delay+1)*0.2 - 0.1) * 0.05 * (r.Float64()*2 - 1)
	processor.setDelayValues(msg, delay, jitter, randomFactor)
	processor.log.Debugf("Adjusted delay of node %s of interface %s to: %f", msg.Tags.Source, msg.Tags.InterfaceName, delay)
	processor.processedMsgChan <- msg
}

func (processor *DefaultProcessor) getLossValue(impairmentsPrefix string) (error, float64) {
	loss := processor.config.GetValue(impairmentsPrefix + "loss")
	if loss != "" {
		if lossValue, err := strconv.ParseFloat(loss, 64); err != nil {
			return fmt.Errorf("Failed to convert loss to float64: %v", err), 0
		} else {
			return nil, lossValue
		}
	}
	return nil, 0.0
}

func (processor *DefaultProcessor) setLossValue(msg *consumer.LossMessage, loss float64, randomFactor float64) {
	msg.LossPercentage = loss + loss*randomFactor
}
func (processor *DefaultProcessor) processLossMessage(msg *consumer.LossMessage) {
	processor.log.Debugf("Process loss of node %s of interface %s", msg.Tags.Source, msg.Tags.InterfaceName)
	err, shortInterfaceName := processor.shortenInterfaceName(msg.Tags.InterfaceName)
	if err != nil {
		processor.log.Debugf("Failed to shorten interface name: %v", err)
		return
	}
	err, loss := processor.getLossValue(processor.helper.GetDefaultImpairmentsPrefix(msg.Tags.Source, shortInterfaceName))
	if err != nil {
		processor.log.Errorf("Failed to get loss value: %v", err)
		return
	}

	// Add normalized random factor to the loss
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomFactor := (math.Log10(loss+1)*0.2 - 0.1) * 0.1 * (r.Float64()*2 - 1)
	processor.setLossValue(msg, loss, randomFactor)

	processor.log.Debugf("Adjusted loss of node %s of interface %s to: %f", msg.Tags.Source, msg.Tags.InterfaceName, loss)
	processor.processedMsgChan <- msg
}

func (proessor *DefaultProcessor) getBandwidthValue(impairmentsPrefix string) (error, float64) {
	bandwidth := proessor.config.GetValue(impairmentsPrefix + "rate")
	if bandwidth != "" {
		if bandwidthValue, err := strconv.ParseFloat(bandwidth, 64); err != nil {
			return fmt.Errorf("Failed to convert bandwidth to float64: %v", err), 0
		} else {
			return nil, bandwidthValue
		}
	}
	return nil, 1000000 // 1Gbps
}
func (processor *DefaultProcessor) processBandwidthMessage(msg *consumer.BandwidthMessage) {
	processor.log.Debugf("Process bandwidth of node %s of interface %s", msg.Tags.Source, msg.Tags.InterfaceName)
	err, shortInterfaceName := processor.shortenInterfaceName(msg.Tags.InterfaceName)
	if err != nil {
		processor.log.Debugf("Failed to shorten interface name: %v", err)
		return
	}
	err, bandwidth := processor.getBandwidthValue(processor.helper.GetDefaultImpairmentsPrefix(msg.Tags.Source, shortInterfaceName))
	if err != nil {
		processor.log.Errorf("Failed to get bandwidth value: %v", err)
		return
	}
	msg.Bandwidth = bandwidth
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
			processor.log.Infoln("Stopping processor")
			return
		}
	}
}

func (processor *DefaultProcessor) Stop() {
	processor.quitChan <- true
}
