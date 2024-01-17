package publisher

import (
	"testing"
	"time"

	"github.com/IBM/sarama/mocks"
	"github.com/hawkv6/clab-telemetry-linker/pkg/consumer"
	"github.com/stretchr/testify/assert"
)

func TestNewKafkaPublisher(t *testing.T) {
	type args struct {
		kafkaBroker string
		kafkaTopic  string
		msgChan     chan consumer.Message
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestNewKafkaPublisher",
			args: args{
				kafkaBroker: "localhost:9092",
				kafkaTopic:  "test",
				msgChan:     make(chan consumer.Message),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kafkaPublisher := NewKafkaPublisher(tt.args.kafkaBroker, tt.args.kafkaTopic, tt.args.msgChan)
			assert.NotNil(t, kafkaPublisher)
		})
	}
}

func TestKafkaPublisher_Init(t *testing.T) {
	type fields struct {
		kafkaBroker string
		kafkaTopic  string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "TestKafkaPublisher_Init",
			fields: fields{
				kafkaBroker: "localhost:9092",
				kafkaTopic:  "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			publisher := NewKafkaPublisher(tt.fields.kafkaBroker, tt.fields.kafkaTopic, make(chan consumer.Message))
			assert.Error(t, publisher.Init())
		})
	}
}

func TestKafkaPublisher_encodeTags(t *testing.T) {
	type fields struct {
		kafkaBroker string
		kafkaTopic  string
	}
	type args struct {
		msg  consumer.TelemetryMessage
		tags consumer.MessageTags
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Test encode tags with Node",
			fields: fields{
				kafkaBroker: "localhost:9092",
				kafkaTopic:  "test",
			},
			args: args{
				tags: consumer.MessageTags{
					Host:          "telegraf",
					InterfaceName: "GigabitEthernet0/0/0/0",
					Node:          "0/RP0/CPU0",
					Path:          "Cisco-IOS-XR-perf-meas-oper:performance-measurement/nodes/node/interfaces/interface-details/interface-detail",
					Source:        "XR-1",
					Subscription:  "hawk-metrics",
				},
				msg: consumer.TelemetryMessage{
					Name:      "performance-measurement",
					Timestamp: 1704728135,
				},
			},
			want: "performance-measurement,host=telegraf,interface_name=GigabitEthernet0/0/0/0,node=0/RP0/CPU0,path=Cisco-IOS-XR-perf-meas-oper:performance-measurement/nodes/node/interfaces/interface-details/interface-detail,source=XR-1,subscription=hawk-metrics",
		},
		{
			name: "Test encode tags without Node",
			fields: fields{
				kafkaBroker: "localhost:9092",
				kafkaTopic:  "test",
			},
			args: args{
				tags: consumer.MessageTags{
					Host:          "telegraf",
					InterfaceName: "GigabitEthernet0/0/0/0",
					Path:          "Cisco-IOS-XR-perf-meas-oper:performance-measurement/nodes/node/interfaces/interface-details/interface-detail",
					Source:        "XR-1",
					Subscription:  "hawk-metrics",
				},
				msg: consumer.TelemetryMessage{
					Name:      "performance-measurement",
					Timestamp: 1704728135,
				},
			},
			want: "performance-measurement,host=telegraf,interface_name=GigabitEthernet0/0/0/0,path=Cisco-IOS-XR-perf-meas-oper:performance-measurement/nodes/node/interfaces/interface-details/interface-detail,source=XR-1,subscription=hawk-metrics",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			publisher := NewKafkaPublisher(tt.fields.kafkaBroker, tt.fields.kafkaTopic, make(chan consumer.Message))
			tt.args.msg.Tags = tt.args.tags
			enc := publisher.createEncoder(tt.args.msg)
			publisher.encodeTags(&enc, tt.args.tags)
			assert.Equal(t, tt.want, string(enc.Bytes()))
		})
	}
}

func TestKafkaPublisher_encodeDelayMessage(t *testing.T) {
	type fields struct {
		kafkaBroker string
		kafkaTopic  string
	}
	type args struct {
		msg consumer.DelayMessage
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test encode delay message without error",
			fields: fields{
				kafkaBroker: "localhost:9092",
				kafkaTopic:  "test",
			},
			args: args{
				msg: consumer.DelayMessage{
					TelemetryMessage: consumer.TelemetryMessage{
						Name: "performance-measurement",
						Tags: consumer.MessageTags{
							Host:          "telegraf",
							InterfaceName: "GigabitEthernet0/0/0/0",
							Node:          "0/RP0/CPU0",
							Path:          "Cisco-IOS-XR-perf-meas-oper:performance-measurement/nodes/node/interfaces/interface-details/interface-detail",
							Source:        "XR-1",
							Subscription:  "hawk-metrics",
						},
						Timestamp: 1704728135,
					},
					Average:  3000,
					Maximum:  3000,
					Minimum:  3000,
					Variance: 0,
				},
			},
			want: "performance-measurement,host=telegraf,interface_name=GigabitEthernet0/0/0/0,node=0/RP0/CPU0,path=Cisco-IOS-XR-perf-meas-oper:performance-measurement/nodes/node/interfaces/interface-details/interface-detail,source=XR-1,subscription=hawk-metrics delay_measurement_session/last_advertisement_information/advertised_values/average=3000,delay_measurement_session/last_advertisement_information/advertised_values/maximum=3000,delay_measurement_session/last_advertisement_information/advertised_values/minimum=3000,delay_measurement_session/last_advertisement_information/advertised_values/variance=0 1704728135000000000\n",
		},
		{
			name: "Test create Encoder without Node",
			fields: fields{
				kafkaBroker: "localhost:9092",
				kafkaTopic:  "test",
			},
			args: args{
				msg: consumer.DelayMessage{
					TelemetryMessage: consumer.TelemetryMessage{},
				},
			},
			want:    "performance-measurement,host=telegraf,interface_name=GigabitEthernet0/0/0/0,path=Cisco-IOS-XR-perf-meas-oper:performance-measurement/nodes/node/interfaces/interface-details/interface-detail,source=XR-1,subscription=hawk-metrics",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			publisher := NewKafkaPublisher(tt.fields.kafkaBroker, tt.fields.kafkaTopic, make(chan consumer.Message))
			err, byteMsg := publisher.encodeDelayMessage(tt.args.msg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, string(byteMsg))
			}
		})
	}
}

func TestKafkaPublisher_encodeLossMessage(t *testing.T) {
	type fields struct {
		kafkaBroker string
		kafkaTopic  string
	}
	type args struct {
		msg consumer.LossMessage
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test encode loss message without error",
			fields: fields{
				kafkaBroker: "localhost:9092",
				kafkaTopic:  "test",
			},
			args: args{
				msg: consumer.LossMessage{
					TelemetryMessage: consumer.TelemetryMessage{
						Name: "isis",
						Tags: consumer.MessageTags{
							Host:          "telegraf",
							InterfaceName: "GigabitEthernet0/0/0/0",
							Path:          "Cisco-IOS-XR-clns-isis-oper:isis/instances/instance/interfaces/interface",
							Source:        "XR-1",
							Subscription:  "hawk-metrics",
						},
						Timestamp: 1704728135,
					},
					LossPercentage: 10.0,
				},
			},
			want: "isis,host=telegraf,interface_name=GigabitEthernet0/0/0/0,path=Cisco-IOS-XR-clns-isis-oper:isis/instances/instance/interfaces/interface,source=XR-1,subscription=hawk-metrics interface_status_and_data/enabled/packet_loss_percentage=10 1704728135000000000\n",
		},
		{
			name: "Test encode loss message with error",
			fields: fields{
				kafkaBroker: "localhost:9092",
				kafkaTopic:  "test",
			},
			args: args{
				msg: consumer.LossMessage{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			publisher := NewKafkaPublisher(tt.fields.kafkaBroker, tt.fields.kafkaTopic, make(chan consumer.Message))
			err, byteMsg := publisher.encodeLossMessage(tt.args.msg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, string(byteMsg))
			}
		})
	}
}

func TestKafkaPublisher_encodeBandwidthMessage(t *testing.T) {
	type fields struct {
		kafkaBroker string
		kafkaTopic  string
	}
	type args struct {
		msg consumer.BandwidthMessage
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test encode bw message without error",
			fields: fields{
				kafkaBroker: "localhost:9092",
				kafkaTopic:  "test",
			},
			args: args{
				msg: consumer.BandwidthMessage{
					TelemetryMessage: consumer.TelemetryMessage{
						Name: "isis",
						Tags: consumer.MessageTags{
							Host:          "telegraf",
							InterfaceName: "GigabitEthernet0/0/0/0",
							Path:          "Cisco-IOS-XR-clns-isis-oper:isis/instances/instance/interfaces/interface",
							Source:        "XR-1",
							Subscription:  "hawk-metrics",
						},
						Timestamp: 1704728135,
					},
					Bandwidth: 100000,
				},
			},
			want: "isis,host=telegraf,interface_name=GigabitEthernet0/0/0/0,path=Cisco-IOS-XR-clns-isis-oper:isis/instances/instance/interfaces/interface,source=XR-1,subscription=hawk-metrics interface_status_and_data/enabled/bandwidth=100000 1704728135000000000\n",
		},
		{
			name: "Test encode bw message with error",
			fields: fields{
				kafkaBroker: "localhost:9092",
				kafkaTopic:  "test",
			},
			args: args{
				msg: consumer.BandwidthMessage{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			publisher := NewKafkaPublisher(tt.fields.kafkaBroker, tt.fields.kafkaTopic, make(chan consumer.Message))
			err, byteMsg := publisher.encodeBandwidthMessage(tt.args.msg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, string(byteMsg))
			}
		})
	}
}

func TestKafkaPublisher_encodeMessage(t *testing.T) {
	type fields struct {
		kafkaBroker string
		kafkaTopic  string
	}
	type args struct {
		msg consumer.Message
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test encode message with BW message",
			fields: fields{
				kafkaBroker: "localhost:9092",
				kafkaTopic:  "test",
			},
			args: args{
				msg: &consumer.BandwidthMessage{
					TelemetryMessage: consumer.TelemetryMessage{
						Name: "isis",
						Tags: consumer.MessageTags{
							Host:          "telegraf",
							InterfaceName: "GigabitEthernet0/0/0/0",
							Path:          "Cisco-IOS-XR-clns-isis-oper:isis/instances/instance/interfaces/interface",
							Source:        "XR-1",
							Subscription:  "hawk-metrics",
						},
						Timestamp: 1704728135,
					},
					Bandwidth: 100000,
				},
			},
			want: "isis,host=telegraf,interface_name=GigabitEthernet0/0/0/0,path=Cisco-IOS-XR-clns-isis-oper:isis/instances/instance/interfaces/interface,source=XR-1,subscription=hawk-metrics interface_status_and_data/enabled/bandwidth=100000 1704728135000000000\n",
		},
		{
			name: "Test encode message with loss message",
			fields: fields{
				kafkaBroker: "localhost:9092",
				kafkaTopic:  "test",
			},
			args: args{
				msg: &consumer.LossMessage{
					TelemetryMessage: consumer.TelemetryMessage{
						Name: "isis",
						Tags: consumer.MessageTags{
							Host:          "telegraf",
							InterfaceName: "GigabitEthernet0/0/0/0",
							Path:          "Cisco-IOS-XR-clns-isis-oper:isis/instances/instance/interfaces/interface",
							Source:        "XR-1",
							Subscription:  "hawk-metrics",
						},
						Timestamp: 1704728135,
					},
					LossPercentage: 10.0,
				},
			},
			want: "isis,host=telegraf,interface_name=GigabitEthernet0/0/0/0,path=Cisco-IOS-XR-clns-isis-oper:isis/instances/instance/interfaces/interface,source=XR-1,subscription=hawk-metrics interface_status_and_data/enabled/packet_loss_percentage=10 1704728135000000000\n",
		},
		{
			name: "Test encode message with delay message",
			fields: fields{
				kafkaBroker: "localhost:9092",
				kafkaTopic:  "test",
			},
			args: args{
				msg: &consumer.DelayMessage{
					TelemetryMessage: consumer.TelemetryMessage{
						Name: "performance-measurement",
						Tags: consumer.MessageTags{
							Host:          "telegraf",
							InterfaceName: "GigabitEthernet0/0/0/0",
							Node:          "0/RP0/CPU0",
							Path:          "Cisco-IOS-XR-perf-meas-oper:performance-measurement/nodes/node/interfaces/interface-details/interface-detail",
							Source:        "XR-1",
							Subscription:  "hawk-metrics",
						},
						Timestamp: 1704728135,
					},
					Average:  3000,
					Maximum:  3000,
					Minimum:  3000,
					Variance: 0,
				},
			},
			want: "performance-measurement,host=telegraf,interface_name=GigabitEthernet0/0/0/0,node=0/RP0/CPU0,path=Cisco-IOS-XR-perf-meas-oper:performance-measurement/nodes/node/interfaces/interface-details/interface-detail,source=XR-1,subscription=hawk-metrics delay_measurement_session/last_advertisement_information/advertised_values/average=3000,delay_measurement_session/last_advertisement_information/advertised_values/maximum=3000,delay_measurement_session/last_advertisement_information/advertised_values/minimum=3000,delay_measurement_session/last_advertisement_information/advertised_values/variance=0 1704728135000000000\n",
		},
		{
			name: "Test encode message with error",
			args: args{
				msg: &consumer.TelemetryMessage{
					Name: "unknown",
					Tags: consumer.MessageTags{
						Host:          "telegraf",
						InterfaceName: "GigabitEthernet0/0/0/0",
						Node:          "0/RP0/CPU0",
						Path:          "Cisco-IOS-XR-perf-meas-oper:performance-measurement/nodes/node/interfaces/interface-details/interface-detail",
						Source:        "XR-1",
					},
					Timestamp: 1704728135,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			publisher := NewKafkaPublisher(tt.fields.kafkaBroker, tt.fields.kafkaTopic, make(chan consumer.Message))
			err, byteMsg := publisher.encodeMessage(tt.args.msg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, string(byteMsg))
			}
		})
	}
}

func TestKafkaPublisher_publishMessage(t *testing.T) {
	type fields struct {
		kafkaBroker string
		kafkaTopic  string
	}
	type args struct {
		msg consumer.Message
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test publish message with BW message",
			fields: fields{
				kafkaBroker: "localhost:9092",
				kafkaTopic:  "test",
			},
			args: args{
				msg: &consumer.BandwidthMessage{
					TelemetryMessage: consumer.TelemetryMessage{
						Name: "isis",
						Tags: consumer.MessageTags{
							Host:          "telegraf",
							InterfaceName: "GigabitEthernet0/0/0/0",
							Path:          "Cisco-IOS-XR-clns-isis-oper:isis/instances/instance/interfaces/interface",
							Source:        "XR-1",
							Subscription:  "hawk-metrics",
						},
						Timestamp: 1704728135,
					},
					Bandwidth: 100000,
				},
			},
			wantErr: false,
		},
		{
			name: "Test publish message with error",
			args: args{
				msg: &consumer.TelemetryMessage{
					Name: "unknown",
					Tags: consumer.MessageTags{
						Host:          "telegraf",
						InterfaceName: "GigabitEthernet0/0/0/0",
						Node:          "0/RP0/CPU0",
						Path:          "Cisco-IOS-XR-perf-meas-oper:performance-measurement/nodes/node/interfaces/interface-details/interface-detail",
						Source:        "XR-1",
					},
					Timestamp: 1704728135,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			publisher := NewKafkaPublisher(tt.fields.kafkaBroker, tt.fields.kafkaTopic, make(chan consumer.Message))
			if tt.wantErr {
				publisher.publishMessage(tt.args.msg)
			} else {
				publisher.producer = mocks.NewAsyncProducer(t, nil).ExpectInputAndSucceed()
				publisher.publishMessage(tt.args.msg)
			}
		})
	}
}
func TestKafkaPublisher_Stop(t *testing.T) {
	type fields struct {
		kafkaBroker string
		kafkaTopic  string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Test Stop function",
			fields: fields{
				kafkaBroker: "localhost:9092",
				kafkaTopic:  "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			publisher := NewKafkaPublisher(tt.fields.kafkaBroker, tt.fields.kafkaTopic, make(chan consumer.Message))
			publisher.producer = mocks.NewAsyncProducer(t, nil)
			go publisher.Start()
			time.Sleep(1 * time.Second)
			assert.NoError(t, publisher.Stop())
		})
	}
}
