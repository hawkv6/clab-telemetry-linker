package consumer

import (
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/IBM/sarama/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewKafkaConsumer(t *testing.T) {
	type args struct {
		kafkaBroker string
		kafkaTopic  string
	}
	tests := []struct {
		name string
		args args
		want *KafkaConsumer
	}{
		{
			name: "Test New Kafka Consumer",
			args: args{
				kafkaBroker: "localhost:9092",
				kafkaTopic:  "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msgChan := make(chan Message)
			kafkaConsumer := NewKafkaConsumer(tt.args.kafkaBroker, tt.args.kafkaTopic, msgChan)
			assert.NotNil(t, kafkaConsumer)
		})
	}
}

func TestKafkaConsumer_createConfig(t *testing.T) {
	type fields struct {
		kafkaBroker        string
		kafkaTopic         string
		unprocessedMsgChan chan Message
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Test Create Config",
			fields: fields{
				kafkaBroker:        "localhost:9092",
				kafkaTopic:         "test",
				unprocessedMsgChan: make(chan Message),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kafkaConsumer := NewKafkaConsumer(tt.fields.kafkaBroker, tt.fields.kafkaTopic, tt.fields.unprocessedMsgChan)
			kafkaConsumer.createConfig()
			assert.NotNil(t, kafkaConsumer.saramaConfig)
		})
	}
}

func TestKafkaConsumer_createConsumer(t *testing.T) {
	type fields struct {
		kafkaBroker        string
		kafkaTopic         string
		unprocessedMsgChan chan Message
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Test create consumer with invalid broker",
			fields: fields{
				kafkaBroker:        "localhost:9092",
				kafkaTopic:         "test",
				unprocessedMsgChan: make(chan Message),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kafkaConsumer := NewKafkaConsumer(tt.fields.kafkaBroker, tt.fields.kafkaTopic, tt.fields.unprocessedMsgChan)
			assert.Error(t, kafkaConsumer.createConsumer())
		})
	}
}
func TestKafkaConsumer_createParitionConsumer(t *testing.T) {
	type fields struct {
		kafkaBroker        string
		kafkaTopic         string
		unprocessedMsgChan chan Message
	}
	tests := []struct {
		name     string
		fields   fields
		wantsErr bool
	}{
		{
			name: "Test create partition consumer with valid consumer",
			fields: fields{
				kafkaBroker:        "localhost:9092",
				kafkaTopic:         "test",
				unprocessedMsgChan: make(chan Message),
			},
			wantsErr: false,
		},
		{
			name: "Test create partition consumer with invalid consumer",
			fields: fields{
				kafkaBroker:        "localhost:9092",
				kafkaTopic:         "test",
				unprocessedMsgChan: make(chan Message),
			},
			wantsErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kafkaConsumer := NewKafkaConsumer(tt.fields.kafkaBroker, tt.fields.kafkaTopic, tt.fields.unprocessedMsgChan)
			saramaConsumer := mocks.NewConsumer(t, nil)
			saramaConsumer.ExpectConsumePartition(tt.fields.kafkaTopic, 0, sarama.OffsetNewest)
			kafkaConsumer.saramaConsumer = saramaConsumer
			if tt.wantsErr {
				assert.Error(t, kafkaConsumer.createParitionConsumer())
			} else {
				assert.NoError(t, kafkaConsumer.createParitionConsumer())
			}
		})
	}
}

func TestKafkaConsumer_Init(t *testing.T) {
	type fields struct {
		kafkaBroker        string
		kafkaTopic         string
		unprocessedMsgChan chan Message
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Test init function with invalid consumer",
			fields: fields{
				kafkaBroker:        "localhost:9092",
				kafkaTopic:         "test",
				unprocessedMsgChan: make(chan Message),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kafkaConsumer := NewKafkaConsumer(tt.fields.kafkaBroker, tt.fields.kafkaTopic, tt.fields.unprocessedMsgChan)
			assert.Error(t, kafkaConsumer.Init())
		})
	}
}

func TestKafkaConsumer_UnmarshalTelemetryMessage(t *testing.T) {
	type fields struct {
		kafkaBroker        string
		kafkaTopic         string
		unprocessedMsgChan chan Message
	}
	type args struct {
		message *sarama.ConsumerMessage
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test Unmarshal Telemetry Message with packet loss message",
			fields: fields{
				kafkaBroker:        "localhost:9092",
				kafkaTopic:         "test",
				unprocessedMsgChan: make(chan Message),
			},
			args: args{
				message: &sarama.ConsumerMessage{
					Value: []byte(`{
						"fields": {
							"interface_status_and_data/enabled/packet_loss_percentage": 0
						},
						"name": "isis",
						"tags": {
							"host": "telegraf",
							"instance_name": "1",
							"interface_name": "GigabitEthernet0/0/0/0",
							"path": "Cisco-IOS-XR-clns-isis-oper:isis/instances/instance/interfaces/interface",
							"source": "XR-1",
							"subscription": "hawk-metrics"
						},
						"timestamp": 1704728296
					}`),
				},
			},
			wantErr: false,
		},
		{
			name: "Test Unmarshal Telemetry Message with delay message",
			fields: fields{
				kafkaBroker:        "localhost:9092",
				kafkaTopic:         "test",
				unprocessedMsgChan: make(chan Message),
			},
			args: args{
				message: &sarama.ConsumerMessage{
					Value: []byte(`{
						"fields": {
						  "delay_measurement_session/last_advertisement_information/advertised_values/average": 10000,
						  "delay_measurement_session/last_advertisement_information/advertised_values/maximum": 10000,
						  "delay_measurement_session/last_advertisement_information/advertised_values/minimum": 10000,
						  "delay_measurement_session/last_advertisement_information/advertised_values/variance": 0
						},
						"name": "performance-measurement",
						"tags": {
						  "host": "telegraf",
						  "interface_name": "GigabitEthernet0/0/0/1",
						  "node": "0/RP0/CPU0",
						  "path": "Cisco-IOS-XR-perf-meas-oper:performance-measurement/nodes/node/interfaces/interface-details/interface-detail",
						  "source": "XR-1",
						  "subscription": "hawk-metrics"
						},
						"timestamp": 1704728135
					  }`),
				},
			},
			wantErr: false,
		},
		{
			name: "Test Unmarshal Telemetry Message with bandwidth message",
			fields: fields{
				kafkaBroker:        "localhost:9092",
				kafkaTopic:         "test",
				unprocessedMsgChan: make(chan Message),
			},
			args: args{
				message: &sarama.ConsumerMessage{
					Value: []byte(`{
						"fields": {
							"interface_status_and_data/enabled/bandwidth": 1000000
						},
						"name": "isis",
						"tags": {
							"host": "telegraf",
							"instance_name": "1",
							"interface_name": "GigabitEthernet0/0/0/0",
							"path": "Cisco-IOS-XR-clns-isis-oper:isis/instances/instance/interfaces/interface",
							"source": "XR-1",
							"subscription": "hawk-metrics"
						},
						"timestamp": 1704728369
					}`),
				},
			},
			wantErr: false,
		},
		{
			name: "Test Unmarshal Telemetry Message with invalid message",
			fields: fields{
				kafkaBroker:        "localhost:9092",
				kafkaTopic:         "test",
				unprocessedMsgChan: make(chan Message),
			},
			args: args{
				message: &sarama.ConsumerMessage{
					Value: []byte(`{
						"title": "invalid message",
					}`),
				},
			},
			wantErr: true,
		},
		{
			name: "Test Unmarshal Telemetry Message with invalid json",
			fields: fields{
				kafkaBroker:        "localhost:9092",
				kafkaTopic:         "test",
				unprocessedMsgChan: make(chan Message),
			},
			args: args{
				message: &sarama.ConsumerMessage{
					Value: []byte(`{
						"title": 
					`),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kafkaConsumer := NewKafkaConsumer(tt.fields.kafkaBroker, tt.fields.kafkaTopic, tt.fields.unprocessedMsgChan)
			err, telemetryMsg := kafkaConsumer.UnmarshalTelemetryMessage(tt.args.message)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, telemetryMsg)
			}
		})
	}
}

func TestKafkaConsumer_UnmarshalDelayMessage(t *testing.T) {
	type fields struct {
		kafkaBroker        string
		kafkaTopic         string
		unprocessedMsgChan chan Message
	}
	type args struct {
		message *sarama.ConsumerMessage
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test Unmarshal Telemetry Message with delay message",
			fields: fields{
				kafkaBroker:        "localhost:9092",
				kafkaTopic:         "test",
				unprocessedMsgChan: make(chan Message),
			},
			args: args{
				message: &sarama.ConsumerMessage{
					Value: []byte(`{
						"fields": {
						  "delay_measurement_session/last_advertisement_information/advertised_values/average": 10000,
						  "delay_measurement_session/last_advertisement_information/advertised_values/maximum": 10000,
						  "delay_measurement_session/last_advertisement_information/advertised_values/minimum": 10000,
						  "delay_measurement_session/last_advertisement_information/advertised_values/variance": 0
						},
						"name": "performance-measurement",
						"tags": {
						  "host": "telegraf",
						  "interface_name": "GigabitEthernet0/0/0/1",
						  "node": "0/RP0/CPU0",
						  "path": "Cisco-IOS-XR-perf-meas-oper:performance-measurement/nodes/node/interfaces/interface-details/interface-detail",
						  "source": "XR-1",
						  "subscription": "hawk-metrics"
						},
						"timestamp": 1704728135
					  }`),
				},
			},
			wantErr: false,
		},
		{
			name: "Test Unmarshal Telemetry Message with wrong fields",
			fields: fields{
				kafkaBroker:        "localhost:9092",
				kafkaTopic:         "test",
				unprocessedMsgChan: make(chan Message),
			},
			args: args{
				message: &sarama.ConsumerMessage{
					Value: []byte(`{
						"fields": {
						  "delay_measurement_session/last_advertisement_information/advertised_values/average": "wrong",
						  "delay_measurement_session/last_advertisement_information/advertised_values/maximum": 10000,
						  "delay_measurement_session/last_advertisement_information/advertised_values/minimum": 10000,
						  "delay_measurement_session/last_advertisement_information/advertised_values/variance": 0
						},
						"name": "performance-measurement",
						"tags": {
						  "host": "telegraf",
						  "interface_name": "GigabitEthernet0/0/0/1",
						  "node": "0/RP0/CPU0",
						  "path": "Cisco-IOS-XR-perf-meas-oper:performance-measurement/nodes/node/interfaces/interface-details/interface-detail",
						  "source": "XR-1",
						  "subscription": "hawk-metrics"
						},
						"timestamp": 1704728135
					  }`),
				},
			},
			wantErr: true,
		},
		{
			name: "Test Unmarshal Telemetry Message with wrong message",
			fields: fields{
				kafkaBroker:        "localhost:9092",
				kafkaTopic:         "test",
				unprocessedMsgChan: make(chan Message),
			},
			args: args{
				message: &sarama.ConsumerMessage{
					Value: []byte(`{
						"fields": {
							"interface_status_and_data/enabled/bandwidth": 1000000
						},
						"name": "isis",
						"tags": {
							"host": "telegraf",
							"instance_name": "1",
							"interface_name": "GigabitEthernet0/0/0/0",
							"path": "Cisco-IOS-XR-clns-isis-oper:isis/instances/instance/interfaces/interface",
							"source": "XR-1",
							"subscription": "hawk-metrics"
						},
						"timestamp": 1704728369
					}`),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kafkaConsumer := NewKafkaConsumer(tt.fields.kafkaBroker, tt.fields.kafkaTopic, tt.fields.unprocessedMsgChan)
			err, telemetryMsg := kafkaConsumer.UnmarshalTelemetryMessage(tt.args.message)
			assert.NoError(t, err)
			err, delayMsg := kafkaConsumer.UnmarshalDelayMessage(*telemetryMsg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, delayMsg)
				assert.Equal(t, delayMsg.Average, 10000.0)
				assert.Equal(t, delayMsg.Maximum, 10000.0)
				assert.Equal(t, delayMsg.Minimum, 10000.0)
				assert.Equal(t, delayMsg.Variance, 0.0)
			}
		})
	}
}

func TestKafkaConsumer_UnmarshalIsisMessage(t *testing.T) {
	type fields struct {
		kafkaBroker        string
		kafkaTopic         string
		unprocessedMsgChan chan Message
	}
	type args struct {
		message *sarama.ConsumerMessage
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test Unmarshal Telemetry Message with unknown isis message",
			fields: fields{
				kafkaBroker:        "localhost:9092",
				kafkaTopic:         "test",
				unprocessedMsgChan: make(chan Message),
			},
			args: args{
				message: &sarama.ConsumerMessage{
					Value: []byte(`{
						"fields": {
						  "delay_measurement_session/last_advertisement_information/advertised_values/average": 10000,
						  "delay_measurement_session/last_advertisement_information/advertised_values/maximum": 10000,
						  "delay_measurement_session/last_advertisement_information/advertised_values/minimum": 10000,
						  "delay_measurement_session/last_advertisement_information/advertised_values/variance": 0
						},
						"name": "performance-measurement",
						"tags": {
						  "host": "telegraf",
						  "interface_name": "GigabitEthernet0/0/0/1",
						  "node": "0/RP0/CPU0",
						  "path": "Cisco-IOS-XR-perf-meas-oper:performance-measurement/nodes/node/interfaces/interface-details/interface-detail",
						  "source": "XR-1",
						  "subscription": "hawk-metrics"
						},
						"timestamp": 1704728135
					  }`),
				},
			},
			wantErr: true,
		},
		{
			name: "Test Unmarshal Telemetry Message with Bandwidth message",
			fields: fields{
				kafkaBroker:        "localhost:9092",
				kafkaTopic:         "test",
				unprocessedMsgChan: make(chan Message),
			},
			args: args{
				message: &sarama.ConsumerMessage{
					Value: []byte(`{
						"fields": {
							"interface_status_and_data/enabled/bandwidth": 1000000
						},
						"name": "isis",
						"tags": {
							"host": "telegraf",
							"instance_name": "1",
							"interface_name": "GigabitEthernet0/0/0/0",
							"path": "Cisco-IOS-XR-clns-isis-oper:isis/instances/instance/interfaces/interface",
							"source": "XR-1",
							"subscription": "hawk-metrics"
						},
						"timestamp": 1704728369
					}`),
				},
			},
			wantErr: false,
		},
		{
			name: "Test Unmarshal Telemetry Message with packet loss message",
			fields: fields{
				kafkaBroker:        "localhost:9092",
				kafkaTopic:         "test",
				unprocessedMsgChan: make(chan Message),
			},
			args: args{
				message: &sarama.ConsumerMessage{
					Value: []byte(`{
						"fields": {
							"interface_status_and_data/enabled/packet_loss_percentage": 0
						},
						"name": "isis",
						"tags": {
							"host": "telegraf",
							"instance_name": "1",
							"interface_name": "GigabitEthernet0/0/0/0",
							"path": "Cisco-IOS-XR-clns-isis-oper:isis/instances/instance/interfaces/interface",
							"source": "XR-1",
							"subscription": "hawk-metrics"
						},
						"timestamp": 1704728296
					}`),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kafkaConsumer := NewKafkaConsumer(tt.fields.kafkaBroker, tt.fields.kafkaTopic, tt.fields.unprocessedMsgChan)
			err, telemetryMsg := kafkaConsumer.UnmarshalTelemetryMessage(tt.args.message)
			assert.NoError(t, err)
			err, delayMsg := kafkaConsumer.UnmarshalIsisMessage(*telemetryMsg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, delayMsg)
			}
		})
	}
}
func TestKafkaConsumer_UnmarshalLossMessage(t *testing.T) {
	type fields struct {
		kafkaBroker        string
		kafkaTopic         string
		unprocessedMsgChan chan Message
	}
	type args struct {
		message *sarama.ConsumerMessage
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test Unmarshal Loss Message with valid message",
			fields: fields{
				kafkaBroker:        "localhost:9092",
				kafkaTopic:         "test",
				unprocessedMsgChan: make(chan Message),
			},
			args: args{
				message: &sarama.ConsumerMessage{
					Value: []byte(`{
						"fields": {
							"interface_status_and_data/enabled/packet_loss_percentage": 0
						},
						"name": "isis",
						"tags": {
							"host": "telegraf",
							"instance_name": "1",
							"interface_name": "GigabitEthernet0/0/0/0",
							"path": "Cisco-IOS-XR-clns-isis-oper:isis/instances/instance/interfaces/interface",
							"source": "XR-1",
							"subscription": "hawk-metrics"
						},
						"timestamp": 1704728296
					}`),
				},
			},
			wantErr: false,
		},
		{
			name: "Test Unmarshal Loss Message with invalid message",
			fields: fields{
				kafkaBroker:        "localhost:9092",
				kafkaTopic:         "test",
				unprocessedMsgChan: make(chan Message),
			},
			args: args{
				message: &sarama.ConsumerMessage{
					Value: []byte(`{
						"fields": {
							"interface_status_and_data/enabled/packet_loss_percentage": "wrong"
						},
						"name": "isis",
						"tags": {
							"host": "telegraf",
							"instance_name": "1",
							"interface_name": "GigabitEthernet0/0/0/0",
							"path": "Cisco-IOS-XR-clns-isis-oper:isis/instances/instance/interfaces/interface",
							"source": "XR-1",
							"subscription": "hawk-metrics"
						},
						"timestamp": 1704728296
					}`),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kafkaConsumer := NewKafkaConsumer(tt.fields.kafkaBroker, tt.fields.kafkaTopic, tt.fields.unprocessedMsgChan)
			err, telemetryMsg := kafkaConsumer.UnmarshalTelemetryMessage(tt.args.message)
			assert.NoError(t, err)
			err, isisMsg := kafkaConsumer.UnmarshalLossMessage(*telemetryMsg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, isisMsg)
			}
		})
	}
}

func TestKafkaConsumer_UnmarshalBandwidthMessage(t *testing.T) {
	type fields struct {
		kafkaBroker        string
		kafkaTopic         string
		unprocessedMsgChan chan Message
	}
	type args struct {
		message *sarama.ConsumerMessage
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test Unmarshal valid Bandwidth message",
			fields: fields{
				kafkaBroker:        "localhost:9092",
				kafkaTopic:         "test",
				unprocessedMsgChan: make(chan Message),
			},
			args: args{
				message: &sarama.ConsumerMessage{
					Value: []byte(`{
						"fields": {
							"interface_status_and_data/enabled/bandwidth": 1000000
						},
						"name": "isis",
						"tags": {
							"host": "telegraf",
							"instance_name": "1",
							"interface_name": "GigabitEthernet0/0/0/0",
							"path": "Cisco-IOS-XR-clns-isis-oper:isis/instances/instance/interfaces/interface",
							"source": "XR-1",
							"subscription": "hawk-metrics"
						},
						"timestamp": 1704728369
					}`),
				},
			},
			wantErr: false,
		},
		{
			name: "Test Unmarshal invalid Bandwidth message",
			fields: fields{
				kafkaBroker:        "localhost:9092",
				kafkaTopic:         "test",
				unprocessedMsgChan: make(chan Message),
			},
			args: args{
				message: &sarama.ConsumerMessage{
					Value: []byte(`{
						"fields": {
							"interface_status_and_data/enabled/bandwidth": "wrong"
						},
						"name": "isis",
						"tags": {
							"host": "telegraf",
							"instance_name": "1",
							"interface_name": "GigabitEthernet0/0/0/0",
							"path": "Cisco-IOS-XR-clns-isis-oper:isis/instances/instance/interfaces/interface",
							"source": "XR-1",
							"subscription": "hawk-metrics"
						},
						"timestamp": 1704728369
					}`),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kafkaConsumer := NewKafkaConsumer(tt.fields.kafkaBroker, tt.fields.kafkaTopic, tt.fields.unprocessedMsgChan)
			err, telemetryMsg := kafkaConsumer.UnmarshalTelemetryMessage(tt.args.message)
			assert.NoError(t, err)
			err, bwMsg := kafkaConsumer.UnmarshalBandwidthMessage(*telemetryMsg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, bwMsg)
				assert.Equal(t, bwMsg.Bandwidth, 1000000.0)
			}
		})
	}
}

func TestKafkaConsumer_processMessage(t *testing.T) {
	type fields struct {
		kafkaBroker        string
		kafkaTopic         string
		unprocessedMsgChan chan Message
	}
	type args struct {
		message *sarama.ConsumerMessage
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test processMessage with packet loss message",
			fields: fields{
				kafkaBroker:        "localhost:9092",
				kafkaTopic:         "test",
				unprocessedMsgChan: make(chan Message),
			},
			args: args{
				message: &sarama.ConsumerMessage{
					Value: []byte(`{
						"fields": {
							"interface_status_and_data/enabled/packet_loss_percentage": 0
						},
						"name": "isis",
						"tags": {
							"host": "telegraf",
							"instance_name": "1",
							"interface_name": "GigabitEthernet0/0/0/0",
							"path": "Cisco-IOS-XR-clns-isis-oper:isis/instances/instance/interfaces/interface",
							"source": "XR-1",
							"subscription": "hawk-metrics"
						},
						"timestamp": 1704728296
					}`),
				},
			},
			wantErr: false,
		},
		{
			name: "Test processMessage with delay message",
			fields: fields{
				kafkaBroker:        "localhost:9092",
				kafkaTopic:         "test",
				unprocessedMsgChan: make(chan Message),
			},
			args: args{
				message: &sarama.ConsumerMessage{
					Value: []byte(`{
						"fields": {
						  "delay_measurement_session/last_advertisement_information/advertised_values/average": 10000,
						  "delay_measurement_session/last_advertisement_information/advertised_values/maximum": 10000,
						  "delay_measurement_session/last_advertisement_information/advertised_values/minimum": 10000,
						  "delay_measurement_session/last_advertisement_information/advertised_values/variance": 0
						},
						"name": "performance-measurement",
						"tags": {
						  "host": "telegraf",
						  "interface_name": "GigabitEthernet0/0/0/1",
						  "node": "0/RP0/CPU0",
						  "path": "Cisco-IOS-XR-perf-meas-oper:performance-measurement/nodes/node/interfaces/interface-details/interface-detail",
						  "source": "XR-1",
						  "subscription": "hawk-metrics"
						},
						"timestamp": 1704728135
					  }`),
				},
			},
			wantErr: false,
		},
		{
			name: "Test processMessage with bandwidth message",
			fields: fields{
				kafkaBroker:        "localhost:9092",
				kafkaTopic:         "test",
				unprocessedMsgChan: make(chan Message),
			},
			args: args{
				message: &sarama.ConsumerMessage{
					Value: []byte(`{
						"fields": {
							"interface_status_and_data/enabled/bandwidth": 1000000
						},
						"name": "isis",
						"tags": {
							"host": "telegraf",
							"instance_name": "1",
							"interface_name": "GigabitEthernet0/0/0/0",
							"path": "Cisco-IOS-XR-clns-isis-oper:isis/instances/instance/interfaces/interface",
							"source": "XR-1",
							"subscription": "hawk-metrics"
						},
						"timestamp": 1704728369
					}`),
				},
			},
			wantErr: false,
		},
		{
			name: "Test processMessage with invalid isis bandwidth message",
			fields: fields{
				kafkaBroker:        "localhost:9092",
				kafkaTopic:         "test",
				unprocessedMsgChan: make(chan Message),
			},
			args: args{
				message: &sarama.ConsumerMessage{
					Value: []byte(`{
						"fields": {
							"interface_status_and_data/enabled/bandwidth": "invalid"
						},
						"name": "isis",
						"tags": {
							"host": "telegraf",
							"instance_name": "1",
							"interface_name": "GigabitEthernet0/0/0/0",
							"path": "Cisco-IOS-XR-clns-isis-oper:isis/instances/instance/interfaces/interface",
							"source": "XR-1",
							"subscription": "hawk-metrics"
						},
						"timestamp": 1704728369
					}`),
				},
			},
			wantErr: true,
		},
		{
			name: "Test processMessage with invalid unknown message",
			fields: fields{
				kafkaBroker:        "localhost:9092",
				kafkaTopic:         "test",
				unprocessedMsgChan: make(chan Message),
			},
			args: args{
				message: &sarama.ConsumerMessage{
					Value: []byte(`{
						"fields": {
							"interface_status_and_data/enabled/bandwidth": "invalid"
						},
						"name": "unknown",
						"tags": {
							"host": "telegraf",
							"instance_name": "1",
							"interface_name": "GigabitEthernet0/0/0/0",
							"path": "Cisco-IOS-XR-clns-isis-oper:isis/instances/instance/interfaces/interface",
							"source": "XR-1",
							"subscription": "hawk-metrics"
						},
						"timestamp": 1704728369
					}`),
				},
			},
			wantErr: true,
		},
		{
			name: "Test processMessage with invalid message",
			fields: fields{
				kafkaBroker:        "localhost:9092",
				kafkaTopic:         "test",
				unprocessedMsgChan: make(chan Message),
			},
			args: args{
				message: &sarama.ConsumerMessage{
					Value: []byte(`{
						"name": "unknown",
					}`),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kafkaConsumer := NewKafkaConsumer(tt.fields.kafkaBroker, tt.fields.kafkaTopic, tt.fields.unprocessedMsgChan)
			go kafkaConsumer.processMessage(tt.args.message)
			time.Sleep(1 * time.Second)
			if tt.wantErr {
				select {
				case <-kafkaConsumer.unprocessedMsgChan:
					assert.Fail(t, "Message should not be sent to unprocessedMsgChan")
				default:
					assert.NoError(t, nil)
				}
			} else {
				select {
				case <-kafkaConsumer.unprocessedMsgChan:
					assert.NoError(t, nil)
				default:
					assert.Fail(t, "Message should be sent to unprocessedMsgChan")
				}
			}
		})
	}
}

func TestKafkaConsumer_Stop(t *testing.T) {
	type fields struct {
		kafkaBroker        string
		kafkaTopic         string
		unprocessedMsgChan chan Message
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Test Stop function",
			fields: fields{
				kafkaBroker:        "localhost:9092",
				kafkaTopic:         "test",
				unprocessedMsgChan: make(chan Message),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kafkaConsumer := NewKafkaConsumer(tt.fields.kafkaBroker, tt.fields.kafkaTopic, tt.fields.unprocessedMsgChan)
			consumer := mocks.NewConsumer(t, nil)
			consumer.ExpectConsumePartition(tt.fields.kafkaTopic, 0, sarama.OffsetNewest)
			kafkaConsumer.saramaConsumer = consumer
			assert.NoError(t, kafkaConsumer.createParitionConsumer())
			go kafkaConsumer.Start()
			time.Sleep(1 * time.Second)
			assert.NoError(t, kafkaConsumer.Stop())
		})
	}
}
