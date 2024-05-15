package consumer

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/hawkv6/clab-telemetry-linker/pkg/logging"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	log                     *logrus.Entry
	kafkaBroker             string
	kafkaTopic              string
	unprocessedMsgChan      chan Message
	quitChan                chan bool
	saramaConfig            *sarama.Config
	saramaConsumer          sarama.Consumer
	saramaPartitionConsumer sarama.PartitionConsumer
}

func NewKafkaConsumer(kafkaBroker, kafkaTopic string, msgChan chan Message) *KafkaConsumer {
	return &KafkaConsumer{
		log:                logging.DefaultLogger.WithField("subsystem", subsystem),
		kafkaBroker:        kafkaBroker,
		kafkaTopic:         kafkaTopic,
		unprocessedMsgChan: msgChan,
		quitChan:           make(chan bool),
	}
}

func (consumer *KafkaConsumer) createConfig() {
	consumer.saramaConfig = sarama.NewConfig()
	consumer.saramaConfig.Net.DialTimeout = time.Second * 5
}

func (consumer *KafkaConsumer) createConsumer() error {
	consumer.createConfig()
	saramaConsumer, err := sarama.NewConsumer([]string{consumer.kafkaBroker}, consumer.saramaConfig)
	if err != nil {
		consumer.log.Debugln("Error creating consumer: ", err)
		return err
	}
	consumer.log.Debugln("Successfully created Kafka consumer for broker: ", consumer.kafkaBroker)
	consumer.saramaConsumer = saramaConsumer
	return nil
}
func (consumer *KafkaConsumer) createParitionConsumer() error {
	partitionConsumer, err := consumer.saramaConsumer.ConsumePartition(consumer.kafkaTopic, 0, sarama.OffsetNewest)
	if err != nil {
		consumer.log.Debugln("Error partition consumer: ", err)
		return err
	}
	consumer.saramaPartitionConsumer = partitionConsumer
	consumer.log.Debugln("Successfully created Kafka partition consumer for topic: ", consumer.kafkaTopic)
	return nil
}

func (consumer *KafkaConsumer) Init() error {
	if err := consumer.createConsumer(); err != nil {
		return err
	}
	if err := consumer.createParitionConsumer(); err != nil {
		return err
	}
	return nil
}

func (consumer *KafkaConsumer) UnmarshalTelemetryMessage(message *sarama.ConsumerMessage) (*TelemetryMessage, error) {
	consumer.log.Debugln("Received JSON message: ", string(message.Value))
	var telemetryMessage TelemetryMessage
	if err := json.Unmarshal([]byte(message.Value), &telemetryMessage); err != nil {
		consumer.log.Debugln("Error unmarshalling message: ", err)
		return nil, err
	}
	return &telemetryMessage, nil
}

func (consumer *KafkaConsumer) UnmarshalDelayMessage(telemetryMessage TelemetryMessage) (*DelayMessage, error) {
	delayMessage := DelayMessage{TelemetryMessage: telemetryMessage}

	fields := map[string]*uint32{
		"delay_measurement_session/last_advertisement_information/advertised_values/average":  &delayMessage.Average,
		"delay_measurement_session/last_advertisement_information/advertised_values/minimum":  &delayMessage.Minimum,
		"delay_measurement_session/last_advertisement_information/advertised_values/maximum":  &delayMessage.Maximum,
		"delay_measurement_session/last_advertisement_information/advertised_values/variance": &delayMessage.Variance,
	}
	for key, field := range fields {
		value, ok := telemetryMessage.Fields[key].(float64)
		if !ok {
			return nil, fmt.Errorf("unable to convert %s to float64", key)
		}
		*field = uint32(value)
	}
	return &delayMessage, nil
}

func (consumer *KafkaConsumer) UnmarshalIsisMessage(telemetryMessage TelemetryMessage) ([]Message, error) {
	var messages []Message
	var err error
	var msg Message

	if telemetryMessage.Fields["interface_status_and_data/enabled/packet_loss_percentage"] != nil {
		msg, err = consumer.UnmarshalLossMessage(telemetryMessage)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	if telemetryMessage.Fields["interface_status_and_data/enabled/bandwidth"] != nil {
		msg, err = consumer.UnmarshalBandwidthMessage(telemetryMessage)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	if len(messages) == 0 {
		return nil, fmt.Errorf("Received unknown ISIS message: %v", telemetryMessage)
	}

	return messages, nil
}

func (consumer *KafkaConsumer) UnmarshalLossMessage(telemetryMessage TelemetryMessage) (*LossMessage, error) {
	lossMessage := LossMessage{TelemetryMessage: telemetryMessage}
	value, ok := telemetryMessage.Fields["interface_status_and_data/enabled/packet_loss_percentage"].(float64)
	if !ok {
		return nil, fmt.Errorf("unable to convert packet_loss_percentage to float")
	}
	lossMessage.LossPercentage = value
	return &lossMessage, nil
}

func (consumer *KafkaConsumer) UnmarshalBandwidthMessage(telemetryMessage TelemetryMessage) (*BandwidthMessage, error) {
	bandwidthMessage := BandwidthMessage{TelemetryMessage: telemetryMessage}
	value, ok := telemetryMessage.Fields["interface_status_and_data/enabled/bandwidth"].(float64)
	if !ok {
		return nil, fmt.Errorf("unable to convert bandwidth to float64")
	}
	bandwidthMessage.Bandwidth = value
	return &bandwidthMessage, nil
}

func (consumer *KafkaConsumer) processMessage(message *sarama.ConsumerMessage) {
	telemetryMessage, err := consumer.UnmarshalTelemetryMessage(message)
	if err != nil {
		return
	}
	if telemetryMessage.Name == "performance-measurement" {
		delayMessage, err := consumer.UnmarshalDelayMessage(*telemetryMessage)
		if err != nil {
			return
		}
		consumer.unprocessedMsgChan <- delayMessage
	} else if telemetryMessage.Name == "isis" {
		isisMessages, err := consumer.UnmarshalIsisMessage(*telemetryMessage)
		if err != nil {
			return
		}
		for _, isisMessage := range isisMessages {
			consumer.unprocessedMsgChan <- isisMessage
		}
	} else {
		consumer.log.Debugf("Skipping unknown message: %v", telemetryMessage)
		return
	}
}

func (consumer *KafkaConsumer) Start() {
	consumer.log.Infof("Start consuming messages from broker %s and topic %s", consumer.kafkaBroker, consumer.kafkaTopic)
	for {
		select {
		case message := <-consumer.saramaPartitionConsumer.Messages():
			consumer.processMessage(message)
		case <-consumer.quitChan:
			consumer.log.Infoln("Stop consumer with values: ", consumer.kafkaBroker, consumer.kafkaTopic)
			return
		}
	}
}

func (consumer *KafkaConsumer) Stop() error {
	consumer.quitChan <- true
	if err := consumer.saramaPartitionConsumer.Close(); err != nil {
		consumer.log.Errorln("Error closing partition consumer: ", err)
		return err
	}
	if err := consumer.saramaConsumer.Close(); err != nil {
		consumer.log.Errorln("Error closing consumer: ", err)
		return err
	}
	return nil
}
