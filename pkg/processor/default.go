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

func (processor *DefaultProcessor) shortenInterfaceName(name string) (string, error) {
	re := regexp.MustCompile(`GigabitEthernet(\d+)/(\d+)/(\d+)/(\d+)`)
	if !re.MatchString(name) {
		return "", fmt.Errorf("interface name %s does not match expected pattern", name)
	}
	return re.ReplaceAllString(name, "Gi$1-$2-$3-$4"), nil
}

func (processor *DefaultProcessor) getDelayValues(impairmentsPrefix string) (uint32, uint32, error) {
	delay := processor.config.GetValue(impairmentsPrefix + "delay")
	var delayMicroSec uint32 = 0
	var jitterMicroSec uint32 = 0
	if delay != "" {
		delayValue64, err := strconv.ParseUint(delay, 10, 32)
		if err != nil {
			return 0, 0, fmt.Errorf("Failed to convert delay to uint64: %v", err)
		}
		delayValue := uint32(delayValue64)
		delayMicroSec = delayValue * 1000

		jitter := processor.config.GetValue(impairmentsPrefix + "jitter")
		if jitter != "" {
			jitterValue64, err := strconv.ParseUint(jitter, 10, 32)
			if err != nil {
				return 0, 0, fmt.Errorf("Failed to convert jitter to float64: %v", err)
			}
			jitterValue := uint32(jitterValue64)
			jitterMicroSec = jitterValue * 1000
		}
	}
	return delayMicroSec, jitterMicroSec, nil
}

func (processor *DefaultProcessor) setDelayValues(msg *consumer.DelayMessage, delay uint32, jitter uint32, randomFactor float64) {
	if delay == 0 {
		msg.Average = msg.Average + uint32(float64(msg.Average)*randomFactor)
	} else {
		msg.Average = msg.Average + delay + uint32(float64(delay)*randomFactor)
	}
	absVal := uint32(math.Abs(float64(msg.Average) * randomFactor))
	if jitter == 0 {
		msg.Maximum = msg.Average + absVal
		if absVal > msg.Average {
			msg.Minimum = 0
		} else {
			msg.Minimum = msg.Average - absVal
		}
		msg.Variance = msg.Maximum - msg.Minimum
	} else {
		halfJitter := uint32(0.5 * float64(jitter))
		msg.Maximum = msg.Average + halfJitter
		if halfJitter > msg.Average {
			msg.Minimum = 0
		} else {
			msg.Minimum = msg.Average - halfJitter
		}
		msg.Variance = jitter
	}
}

func (processor *DefaultProcessor) processDelayMessage(msg *consumer.DelayMessage) {
	processor.log.Debugf("Process delay of node %s of interface %s", msg.Tags.Source, msg.Tags.InterfaceName)
	shortInterfaceName, err := processor.shortenInterfaceName(msg.Tags.InterfaceName)
	if err != nil {
		processor.log.Debugf("Failed to shorten interface name: %v", err)
		return
	}
	impairmentsPrefix := processor.helper.GetDefaultImpairmentsPrefix(msg.Tags.Source, shortInterfaceName)
	delay, jitter, err := processor.getDelayValues(impairmentsPrefix)
	if err != nil {
		processor.log.Errorf("Failed to get delay values: %v", err)
		return
	}

	// Add normalized random factor to the delay
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomFactor := (math.Log10(float64(delay+1))*0.2 - 0.1) * 0.05 * (r.Float64()*2 - 1)
	processor.setDelayValues(msg, delay, jitter, randomFactor)
	processor.log.Debugf("Adjusted delay of node %s of interface %s to: %d", msg.Tags.Source, msg.Tags.InterfaceName, delay)
	processor.processedMsgChan <- msg
}

func (processor *DefaultProcessor) getLossValue(impairmentsPrefix string) (float64, error) {
	loss := processor.config.GetValue(impairmentsPrefix + "loss")
	if loss != "" {
		if lossValue, err := strconv.ParseFloat(loss, 64); err != nil {
			return 0, fmt.Errorf("Failed to convert loss to float64: %v", err)
		} else {
			return lossValue, nil
		}
	}
	return 0.001, nil
}

func (processor *DefaultProcessor) setLossValue(msg *consumer.LossMessage, loss float64, randomFactor float64) {
	msg.LossPercentage = loss + loss*randomFactor
}
func (processor *DefaultProcessor) processLossMessage(msg *consumer.LossMessage) {
	processor.log.Debugf("Process loss of node %s of interface %s", msg.Tags.Source, msg.Tags.InterfaceName)
	shortInterfaceName, err := processor.shortenInterfaceName(msg.Tags.InterfaceName)
	if err != nil {
		processor.log.Debugf("Failed to shorten interface name: %v", err)
		return
	}
	loss, err := processor.getLossValue(processor.helper.GetDefaultImpairmentsPrefix(msg.Tags.Source, shortInterfaceName))
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

func (proessor *DefaultProcessor) getBandwidthValue(impairmentsPrefix string) (float64, error) {
	bandwidth := proessor.config.GetValue(impairmentsPrefix + "rate")
	if bandwidth != "" {
		if bandwidthValue, err := strconv.ParseFloat(bandwidth, 64); err != nil {
			return 0, fmt.Errorf("Failed to convert bandwidth to float64: %v", err)
		} else {
			return bandwidthValue, nil
		}
	}
	return 1000000, nil // 1Gbps
}
func (processor *DefaultProcessor) processBandwidthMessage(msg *consumer.BandwidthMessage) {
	processor.log.Debugf("Process bandwidth of node %s of interface %s", msg.Tags.Source, msg.Tags.InterfaceName)
	shortInterfaceName, err := processor.shortenInterfaceName(msg.Tags.InterfaceName)
	if err != nil {
		processor.log.Debugf("Failed to shorten interface name: %v", err)
		return
	}
	bandwidth, err := processor.getBandwidthValue(processor.helper.GetDefaultImpairmentsPrefix(msg.Tags.Source, shortInterfaceName))
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
