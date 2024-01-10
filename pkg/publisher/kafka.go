package publisher

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/hawkv6/clab-telemetry-linker/pkg/consumer"
	"github.com/hawkv6/clab-telemetry-linker/pkg/logging"
	"github.com/influxdata/line-protocol/v2/lineprotocol"
	"github.com/sirupsen/logrus"
)

type KafkaPublisher struct {
	log              *logrus.Entry
	kafkaBroker      string
	kafkaTopic       string
	processedMsgChan chan consumer.Message
	quitChan         chan bool
	producer         sarama.AsyncProducer
}

func NewKafkaPublisher(kafkaBroker, kafkaTopic string, msgChan chan consumer.Message) *KafkaPublisher {
	return &KafkaPublisher{
		log:              logging.DefaultLogger.WithField("subsystem", subsystem),
		kafkaBroker:      kafkaBroker,
		kafkaTopic:       kafkaTopic,
		processedMsgChan: msgChan,
		quitChan:         make(chan bool),
	}
}

func (publisher *KafkaPublisher) Init() error {
	producer, err := sarama.NewAsyncProducer([]string{publisher.kafkaBroker}, nil)
	if err != nil {
		publisher.log.Debugln("Error creating producer: ", err)
		return err
	}
	publisher.producer = producer
	return nil
}

func (publisher *KafkaPublisher) createEncoder(msg consumer.TelemetryMessage) lineprotocol.Encoder {
	var enc lineprotocol.Encoder
	enc.SetPrecision(lineprotocol.Nanosecond)
	enc.StartLine(msg.Name)
	return enc
}

func (publisher *KafkaPublisher) encodeTags(enc *lineprotocol.Encoder, tags consumer.MessageTags) {
	enc.AddTag("host", tags.Host)
	enc.AddTag("interface_name", tags.InterfaceName)
	if tags.Node != "" {
		enc.AddTag("node", tags.Node)
	}
	enc.AddTag("path", tags.Path)
	enc.AddTag("source", tags.Source)
	enc.AddTag("subscription", tags.Subscription)
}

func (publisher *KafkaPublisher) encodeDelayMessage(msg consumer.DelayMessage) (error, []byte) {
	enc := publisher.createEncoder(msg.TelemetryMessage)
	publisher.encodeTags(&enc, msg.Tags)
	enc.AddField("delay_measurement_session/last_advertisement_information/advertised_values/average", lineprotocol.MustNewValue(msg.Average))
	enc.AddField("delay_measurement_session/last_advertisement_information/advertised_values/maximum", lineprotocol.MustNewValue(msg.Maximum))
	enc.AddField("delay_measurement_session/last_advertisement_information/advertised_values/minimum", lineprotocol.MustNewValue(msg.Minimum))
	enc.AddField("delay_measurement_session/last_advertisement_information/advertised_values/variance", lineprotocol.MustNewValue(msg.Variance))
	enc.EndLine(time.Unix(msg.Timestamp, 0))
	if err := enc.Err(); err != nil {
		return err, nil
	}
	return nil, enc.Bytes()
}

func (publisher *KafkaPublisher) encodeLossMessage(msg consumer.LossMessage) (error, []byte) {
	enc := publisher.createEncoder(msg.TelemetryMessage)
	publisher.encodeTags(&enc, msg.Tags)
	enc.AddField("interface_status_and_data/enabled/packet_loss_percentage", lineprotocol.MustNewValue(msg.LossPercentage))
	enc.EndLine(time.Unix(msg.Timestamp, 0))
	if err := enc.Err(); err != nil {
		return err, nil
	}
	return nil, enc.Bytes()
}

func (publisher *KafkaPublisher) encodeBandwidthMessage(msg consumer.BandwidthMessage) (error, []byte) {
	enc := publisher.createEncoder(msg.TelemetryMessage)
	publisher.encodeTags(&enc, msg.Tags)
	enc.AddField("interface_status_and_data/enabled/bandwidth", lineprotocol.MustNewValue(msg.Bandwidth))
	enc.EndLine(time.Unix(msg.Timestamp, 0))
	if err := enc.Err(); err != nil {
		return err, nil
	}
	return nil, enc.Bytes()
}

func (publisher *KafkaPublisher) encodeMessage(msg consumer.Message) (error, []byte) {
	switch msg := msg.(type) {
	case *consumer.DelayMessage:
		return publisher.encodeDelayMessage(*msg)
	case *consumer.LossMessage:
		return publisher.encodeLossMessage(*msg)
	case *consumer.BandwidthMessage:
		return publisher.encodeBandwidthMessage(*msg)
	default:
		return fmt.Errorf("Skipping unknown message type: %v", msg), nil
	}
}
func (publisher *KafkaPublisher) publishMessage(msg consumer.Message) {
	err, encodedMsg := publisher.encodeMessage(msg)
	if err != nil {
		publisher.log.Errorln("Error encoding message: ", err)
		return
	}
	select {
	case publisher.producer.Input() <- &sarama.ProducerMessage{Topic: publisher.kafkaTopic, Key: nil, Value: sarama.ByteEncoder(encodedMsg)}:
		publisher.log.Debugf("Successfully enqueued message %v on topic %s\n", string(encodedMsg), publisher.kafkaTopic)
	case err := <-publisher.producer.Errors():
		log.Println("Failed to produce message", err)
	}
}

func (publisher *KafkaPublisher) Start() {
	publisher.log.Infoln("Starting publishing messages to broker", publisher.kafkaBroker, "and topic", publisher.kafkaTopic)
	for {
		select {
		case msg := <-publisher.processedMsgChan:
			publisher.publishMessage(msg)
		case <-publisher.quitChan:
			publisher.log.Infoln("Stopping publisher with broker ", publisher.kafkaBroker, " and topic ", publisher.kafkaTopic, " and quitChan: ", publisher.quitChan)
			return
		}
	}
}

func (publisher *KafkaPublisher) Stop() error {
	publisher.quitChan <- true
	if err := publisher.producer.Close(); err != nil {
		return err
	}
	return nil
}
