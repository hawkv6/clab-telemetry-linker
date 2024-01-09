package consumer

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/hawkv6/clab-telemetry-linker/pkg/logging"
	"github.com/sirupsen/logrus"
)

var subsystem = "consumer"

type Consumer interface {
	Init() error
	Start()
	Stop() error
}

type KafkaConsumer struct {
	log                     *logrus.Entry
	kafkaBroker             string
	kafkaTopic              string
	msgChan                 chan Message
	saramaConfig            *sarama.Config
	saramaConsumer          sarama.Consumer
	saramaPartitionConsumer sarama.PartitionConsumer
}

func NewKafkaConsumer(kafkaBroker, kafkaTopic string, msgChan chan Message) *KafkaConsumer {
	return &KafkaConsumer{
		log:          logging.DefaultLogger.WithField("subsystem", subsystem),
		kafkaBroker:  kafkaBroker,
		kafkaTopic:   kafkaTopic,
		msgChan:      msgChan,
		saramaConfig: sarama.NewConfig(),
	}
}

func (consumer *KafkaConsumer) createConsumer() error {
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

func (consumer *KafkaConsumer) UnmarshalTelemetryMessage(message *sarama.ConsumerMessage) (error, *TelemetryMessage) {
	consumer.log.Debugln("Received JSON message: ", string(message.Value))
	var telemetryMessage TelemetryMessage
	if err := json.Unmarshal([]byte(message.Value), &telemetryMessage); err != nil {
		consumer.log.Debugln("Error unmarshalling message: ", err)
		return err, nil
	}
	consumer.log.Debugf("Telemetry message: %v", telemetryMessage)
	return nil, &telemetryMessage
}

func (consumer *KafkaConsumer) UnmarshalDelayMessage(telemetryMessage TelemetryMessage) (error, *DelayMessage) {
	delayMessage := DelayMessage{TelemetryMessage: telemetryMessage}

	fields := map[string]*float64{
		"delay_measurement_session/last_advertisement_information/advertised_values/average":  &delayMessage.Average,
		"delay_measurement_session/last_advertisement_information/advertised_values/minimum":  &delayMessage.Minimum,
		"delay_measurement_session/last_advertisement_information/advertised_values/maximum":  &delayMessage.Maximum,
		"delay_measurement_session/last_advertisement_information/advertised_values/variance": &delayMessage.Variance,
	}
	for key, field := range fields {
		value, ok := telemetryMessage.Fields[key].(float64)
		if !ok {
			return fmt.Errorf("unable to convert %s to float", key), nil
		}
		*field = value
	}
	consumer.log.Debugf("Received Delay message: %v", delayMessage)
	return nil, &delayMessage
}

func (consumer *KafkaConsumer) UnmarshalIsisMessage(telemetryMessage TelemetryMessage) (error, Message) {
	switch {
	case telemetryMessage.Fields["interface_status_and_data/enabled/packet_loss_percentage"] != nil:
		return consumer.UnmarshalLossMessage(telemetryMessage)
	case telemetryMessage.Fields["interface_status_and_data/enabled/bandwidth"] != nil:
		return consumer.UnmarshalBandwidthMessage(telemetryMessage)
	default:
		msg := fmt.Sprintf("Received unknown ISIS message: %v", telemetryMessage)
		consumer.log.Debugln(msg)
		return fmt.Errorf(msg), nil
	}
}

func (consumer *KafkaConsumer) UnmarshalLossMessage(telemetryMessage TelemetryMessage) (error, *LossMessage) {
	lossMessage := LossMessage{TelemetryMessage: telemetryMessage}
	value, ok := telemetryMessage.Fields["interface_status_and_data/enabled/packet_loss_percentage"].(float64)
	if !ok {
		return fmt.Errorf("unable to convert packet_loss_percentage to float"), nil
	}
	lossMessage.LossPercentage = value
	consumer.log.Debugf("Received Loss message: %v", lossMessage)
	return nil, &lossMessage
}

func (consumer *KafkaConsumer) UnmarshalBandwidthMessage(telemetryMessage TelemetryMessage) (error, *BandwidthMessage) {
	bandwidthMessage := BandwidthMessage{TelemetryMessage: telemetryMessage}
	value, ok := telemetryMessage.Fields["interface_status_and_data/enabled/bandwidth"].(float64)
	if !ok {
		return fmt.Errorf("unable to convert bandwidth to float64"), nil
	}
	bandwidthMessage.Bandwidth = value
	consumer.log.Debugf("Received Bandwidth message: %v", bandwidthMessage)
	return nil, &bandwidthMessage
}

func (consumer *KafkaConsumer) processMessage(message *sarama.ConsumerMessage) {
	err, telemetryMessage := consumer.UnmarshalTelemetryMessage(message)
	if err != nil {
		return
	}
	if telemetryMessage.Name == "performance_monitoring" {
		err, delayMessage := consumer.UnmarshalDelayMessage(*telemetryMessage)
		if err != nil {
			return
		}
		consumer.msgChan <- delayMessage
	} else if telemetryMessage.Name == "isis" {
		err, isisMessage := consumer.UnmarshalIsisMessage(*telemetryMessage)
		if err != nil {
			return
		}
		consumer.log.Infof("Received ISIS message: %v", isisMessage)
		consumer.msgChan <- isisMessage
	} else {
		consumer.log.Debugf("Skipping unknown message: %v", telemetryMessage)
		return
	}
}

func (consumer *KafkaConsumer) Start() {
	consumer.log.Infof("Start consuming messages from broker %s and topic %s", consumer.kafkaBroker, consumer.kafkaTopic)
	defer consumer.saramaPartitionConsumer.Close()
	for {
		message := <-consumer.saramaPartitionConsumer.Messages()
		consumer.processMessage(message)
	}
}

func (consumer *KafkaConsumer) Stop() error {
	consumer.log.Debugln("Stop consumer with values: ", consumer.kafkaBroker, consumer.kafkaTopic)
	if err := consumer.saramaPartitionConsumer.Close(); err != nil {
		consumer.log.Debugln("Error closing partition consumer: ", err)
		return err
	}
	if err := consumer.saramaConsumer.Close(); err != nil {
		consumer.log.Debugln("Error closing consumer: ", err)
		return err
	}
	return nil
}
