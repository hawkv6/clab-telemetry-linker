package processor

import (
	"testing"
	"time"

	"github.com/hawkv6/clab-telemetry-linker/pkg/config"
	"github.com/hawkv6/clab-telemetry-linker/pkg/consumer"
	"github.com/hawkv6/clab-telemetry-linker/pkg/helpers"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewDefaultProcessor(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Test Creating Default Processor",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			config := config.NewMockConfig(ctrl)
			unprocessedMsgChan := make(chan consumer.Message)
			processedMsgChan := make(chan consumer.Message)
			helper := helpers.NewMockHelper(ctrl)
			defaultProcessor := NewDefaultProcessor(config, unprocessedMsgChan, processedMsgChan, helper)
			assert.NotNil(t, defaultProcessor)
		})
	}
}

func TestDefaultProcessor_shortenInterfaceName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test Shorten Interface Name with valid name",
			args: args{
				name: "GigabitEthernet0/0/0/0",
			},
			want:    "Gi0-0-0-0",
			wantErr: false,
		},
		{
			name: "Test Shorten Interface Name with invalid name",
			args: args{
				name: "FastEthernet0/0/0/0",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			config := config.NewMockConfig(ctrl)
			unprocessedMsgChan := make(chan consumer.Message)
			processedMsgChan := make(chan consumer.Message)
			helper := helpers.NewMockHelper(ctrl)
			processor := NewDefaultProcessor(config, unprocessedMsgChan, processedMsgChan, helper)
			err, got := processor.shortenInterfaceName(tt.args.name)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDefaultProcessor_getDelayValues(t *testing.T) {
	tests := []struct {
		name        string
		delay       string
		delayValue  float64
		jitter      string
		jitterValue float64
		wantErr     bool
	}{
		{
			name:        "Test Get Delay Values with delay set and no jitter ",
			delay:       "10",
			delayValue:  10000,
			jitter:      "",
			jitterValue: 0,
			wantErr:     false,
		},
		{
			name:        "Test Get Delay Values no delay set and no jitter ",
			delay:       "",
			delayValue:  0,
			jitter:      "",
			jitterValue: 0,
			wantErr:     false,
		},
		{
			name:        "Test Get Delay Values delay and jitter set",
			delay:       "10",
			delayValue:  10000,
			jitter:      "1",
			jitterValue: 1000,
			wantErr:     false,
		},
		{
			name:    "Test Get Delay Values with invalid delay ",
			delay:   "invalid",
			wantErr: true,
		},
		{
			name:    "Test Get Delay Values with invalid jitter",
			delay:   "10",
			jitter:  "invalid",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			config := config.NewMockConfig(ctrl)
			unprocessedMsgChan := make(chan consumer.Message)
			processedMsgChan := make(chan consumer.Message)
			helper := helpers.NewMockHelper(ctrl)
			processor := NewDefaultProcessor(config, unprocessedMsgChan, processedMsgChan, helper)
			helper.EXPECT().GetDefaultImpairmentsPrefix("XR-1", "Gi0-0-0-0").Return("nodes.XR-1.config.Gi0-0-0-0.impairments.")
			impairmentsPrefix := helper.GetDefaultImpairmentsPrefix("XR-1", "Gi0-0-0-0")
			config.EXPECT().GetValue(impairmentsPrefix + "delay").Return(tt.delay)
			config.EXPECT().GetValue(impairmentsPrefix + "jitter").Return(tt.jitter).AnyTimes()
			err, delay, jitter := processor.getDelayValues(impairmentsPrefix)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, 0.0, delay)
				assert.Equal(t, 0.0, jitter)
			} else {
				assert.Equal(t, tt.delayValue, delay)
				assert.Equal(t, tt.jitterValue, jitter)
				assert.Nil(t, err)
			}
		})
	}
}

func TestDefaultProcessor_setDelayValues(t *testing.T) {
	type want struct {
		Average  float64
		Maximum  float64
		Minimum  float64
		Variance float64
	}
	tests := []struct {
		name         string
		delay        float64
		jitter       float64
		randomFactor float64
		want         want
	}{
		{
			name:         "Test Get Delay Values with delay set and no jitter ",
			delay:        7000,
			jitter:       0,
			randomFactor: 0.1,
			want: want{
				Average:  10700, // 3000 + (7000 + 7000*0.1)
				Maximum:  11770, // 10700 + 10700 *0.1
				Minimum:  9630,  // 10700 - 10700 *0.1
				Variance: 2140,  // 11770 - 9630
			},
		},
		{
			name:         "Test Get Delay Values with delay set and jitter ",
			delay:        7000,
			jitter:       1000,
			randomFactor: 0.1,
			want: want{
				Average:  10700, // 3000 + (7000 + 7000*0.1)
				Maximum:  11200, // 10700 + 0.5 * 1000
				Minimum:  10200, // 10700 - 0.5 * 1000
				Variance: 1000,  // 2000
			},
		},
		{
			name:         "Test Get Delay Values with no delay set and no jitter ",
			delay:        0,
			jitter:       0,
			randomFactor: 0.1,
			want: want{
				Average:  3300, // 3000 + 3000*0.1
				Maximum:  3630, // 3300 + 3300 *0.1
				Minimum:  2970, // 3300 - 3300 *0.1
				Variance: 660,  // 3630 - 2970
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := consumer.DelayMessage{
				TelemetryMessage: consumer.TelemetryMessage{
					Name: "performance_monitoring",
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
			}
			ctrl := gomock.NewController(t)
			config := config.NewMockConfig(ctrl)
			unprocessedMsgChan := make(chan consumer.Message)
			processedMsgChan := make(chan consumer.Message)
			helper := helpers.NewMockHelper(ctrl)
			processor := NewDefaultProcessor(config, unprocessedMsgChan, processedMsgChan, helper)
			processor.setDelayValues(&msg, tt.delay, tt.jitter, tt.randomFactor)
			assert.Equal(t, tt.want.Average, msg.Average)
			assert.Equal(t, tt.want.Maximum, msg.Maximum)
			assert.Equal(t, tt.want.Minimum, msg.Minimum)
			assert.Equal(t, tt.want.Variance, msg.Variance)
		})
	}
}

func TestDefaultProcessor_processDelayMessage(t *testing.T) {
	type fields struct {
		Interface string
	}
	type want struct {
		Average  float64
		Maximum  float64
		Minimum  float64
		Variance float64
		Err      bool
	}
	tests := []struct {
		fields fields
		name   string
		delay  string
		jitter string
		want   want
	}{
		{
			name: "Test with invalid name",
			fields: fields{
				Interface: "Loopback0",
			},
			want: want{
				Err: true,
			},
		},
		{
			name:  "Test with invalid delay",
			delay: "invalid",
			fields: fields{
				Interface: "GigabitEthernet0/0/0/0",
			},
			want: want{
				Err: true,
			},
		},
		{
			name:   "Test with valid message",
			delay:  "7000",
			jitter: "0",
			fields: fields{
				Interface: "GigabitEthernet0/0/0/0",
			},
			want: want{
				Err: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := consumer.DelayMessage{
				TelemetryMessage: consumer.TelemetryMessage{
					Name: "performance_monitoring",
					Tags: consumer.MessageTags{
						Host:          "telegraf",
						InterfaceName: tt.fields.Interface,
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
			}
			ctrl := gomock.NewController(t)
			config := config.NewMockConfig(ctrl)
			unprocessedMsgChan := make(chan consumer.Message)
			processedMsgChan := make(chan consumer.Message)
			helper := helpers.NewMockHelper(ctrl)
			processor := NewDefaultProcessor(config, unprocessedMsgChan, processedMsgChan, helper)
			impairmentsPrefix := "nodes.XR-1.config.Gi0-0-0-0.impairments."
			helper.EXPECT().GetDefaultImpairmentsPrefix("XR-1", "Gi0-0-0-0").Return(impairmentsPrefix).AnyTimes()
			config.EXPECT().GetValue(impairmentsPrefix + "delay").Return(tt.delay).AnyTimes()
			config.EXPECT().GetValue(impairmentsPrefix + "jitter").Return(tt.jitter).AnyTimes()
			go processor.processDelayMessage(&msg)
			time.Sleep(time.Second * 1)
			if tt.want.Err {
				select {
				case <-processor.processedMsgChan:
					assert.Fail(t, "Message should not be sent to processedMsgChan")
				default:
					assert.NoError(t, nil)
				}
			} else {
				select {
				case <-processor.processedMsgChan:
					assert.NoError(t, nil)
				default:
					assert.Fail(t, "Message should be sent to processedMsgChan")
				}
			}
		})
	}
}

func TestDefaultProcessor_getLossValue(t *testing.T) {
	tests := []struct {
		name      string
		loss      string
		lossValue float64
		wantErr   bool
	}{
		{
			name:      "Test Get Loss Values with valid loss",
			loss:      "10",
			lossValue: 10,
			wantErr:   false,
		},
		{
			name:      "Test Get Loss Values with no loss",
			loss:      "",
			lossValue: 0,
			wantErr:   false,
		},
		{
			name:    "Test Get Loss Values with invalid loss",
			loss:    "invalid",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			config := config.NewMockConfig(ctrl)
			unprocessedMsgChan := make(chan consumer.Message)
			processedMsgChan := make(chan consumer.Message)
			helper := helpers.NewMockHelper(ctrl)
			processor := NewDefaultProcessor(config, unprocessedMsgChan, processedMsgChan, helper)
			helper.EXPECT().GetDefaultImpairmentsPrefix("XR-1", "Gi0-0-0-0").Return("nodes.XR-1.config.Gi0-0-0-0.impairments.")
			impairmentsPrefix := helper.GetDefaultImpairmentsPrefix("XR-1", "Gi0-0-0-0")
			config.EXPECT().GetValue(impairmentsPrefix + "loss").Return(tt.loss)
			err, loss := processor.getLossValue(impairmentsPrefix)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, 0.0, loss)
			} else {
				assert.Equal(t, tt.lossValue, loss)
				assert.Nil(t, err)
			}
		})
	}
}

func TestDefaultProcessor_setLossValue(t *testing.T) {
	type want struct {
		Loss float64
	}
	tests := []struct {
		name         string
		loss         float64
		randomFactor float64
		want         want
	}{
		{
			name:         "Test Set Loss Values with valid loss",
			loss:         1,
			randomFactor: 0.1,
			want: want{
				Loss: 1.1,
			},
		},
		{
			name:         "Test Set Loss Values with no loss",
			loss:         0,
			randomFactor: 0.1,
			want: want{
				Loss: 0.0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := consumer.LossMessage{
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
				LossPercentage: 0.0,
			}
			ctrl := gomock.NewController(t)
			config := config.NewMockConfig(ctrl)
			unprocessedMsgChan := make(chan consumer.Message)
			processedMsgChan := make(chan consumer.Message)
			helper := helpers.NewMockHelper(ctrl)
			processor := NewDefaultProcessor(config, unprocessedMsgChan, processedMsgChan, helper)
			processor.setLossValue(&msg, tt.loss, tt.randomFactor)
			assert.Equal(t, tt.want.Loss, msg.LossPercentage)
		})
	}
}

func TestDefaultProcessor_processLossMessage(t *testing.T) {
	type fields struct {
		Interface string
	}
	type want struct {
		Loss float64
		Err  bool
	}
	tests := []struct {
		fields fields
		name   string
		loss   string
		want   want
	}{
		{
			name: "Test with invalid name",
			fields: fields{
				Interface: "Loopback0",
			},
			want: want{
				Err: true,
			},
		},
		{
			name: "Test with invalid loss",
			loss: "invalid",
			fields: fields{
				Interface: "GigabitEthernet0/0/0/0",
			},
			want: want{
				Err: true,
			},
		},
		{
			name: "Test with valid message",
			loss: "10",
			fields: fields{
				Interface: "GigabitEthernet0/0/0/0",
			},
			want: want{
				Err: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := consumer.LossMessage{
				TelemetryMessage: consumer.TelemetryMessage{
					Name: "isis",
					Tags: consumer.MessageTags{
						Host:          "telegraf",
						InterfaceName: tt.fields.Interface,
						Path:          "Cisco-IOS-XR-clns-isis-oper:isis/instances/instance/interfaces/interface",
						Source:        "XR-1",
						Subscription:  "hawk-metrics",
					},
					Timestamp: 1704728135,
				},
				LossPercentage: 0.0,
			}
			ctrl := gomock.NewController(t)
			config := config.NewMockConfig(ctrl)
			unprocessedMsgChan := make(chan consumer.Message)
			processedMsgChan := make(chan consumer.Message)
			helper := helpers.NewMockHelper(ctrl)
			processor := NewDefaultProcessor(config, unprocessedMsgChan, processedMsgChan, helper)
			impairmentsPrefix := "nodes.XR-1.config.Gi0-0-0-0.impairments."
			helper.EXPECT().GetDefaultImpairmentsPrefix("XR-1", "Gi0-0-0-0").Return(impairmentsPrefix).AnyTimes()
			config.EXPECT().GetValue(impairmentsPrefix + "loss").Return(tt.loss).AnyTimes()
			go processor.processLossMessage(&msg)
			time.Sleep(time.Second * 1)
			if tt.want.Err {
				select {
				case <-processor.processedMsgChan:
					assert.Fail(t, "Message should not be sent to processedMsgChan")
				default:
					assert.NoError(t, nil)
				}
			} else {
				select {
				case <-processor.processedMsgChan:
					assert.NoError(t, nil)
				default:
					assert.Fail(t, "Message should be sent to processedMsgChan")
				}
			}
		})
	}
}

func TestDefaultProcessor_getBandwidthValue(t *testing.T) {
	tests := []struct {
		name      string
		rate      string
		rateValue float64
		wantErr   bool
	}{
		{
			name:      "Test Get BW Values with valid BW",
			rate:      "1000000",
			rateValue: 1000000,
			wantErr:   false,
		},
		{
			name:      "Test Get Loss Values with no loss",
			rate:      "",
			rateValue: 1000000,
			wantErr:   false,
		},
		{
			name:    "Test Get Loss Values with invalid loss",
			rate:    "invalid",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			config := config.NewMockConfig(ctrl)
			unprocessedMsgChan := make(chan consumer.Message)
			processedMsgChan := make(chan consumer.Message)
			helper := helpers.NewMockHelper(ctrl)
			processor := NewDefaultProcessor(config, unprocessedMsgChan, processedMsgChan, helper)
			helper.EXPECT().GetDefaultImpairmentsPrefix("XR-1", "Gi0-0-0-0").Return("nodes.XR-1.config.Gi0-0-0-0.impairments.")
			impairmentsPrefix := helper.GetDefaultImpairmentsPrefix("XR-1", "Gi0-0-0-0")
			config.EXPECT().GetValue(impairmentsPrefix + "rate").Return(tt.rate)
			err, rate := processor.getBandwidthValue(impairmentsPrefix)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tt.rateValue, rate)
		})
	}
}
func TestDefaultProcessor_processBandwidthMessage(t *testing.T) {
	type fields struct {
		Interface string
	}
	type want struct {
		Err bool
	}
	tests := []struct {
		name   string
		fields fields
		rate   string
		want   want
	}{
		{
			name: "Test with invalid name",
			fields: fields{
				Interface: "Loopback0",
			},
			want: want{
				Err: true,
			},
		},
		{
			name: "Test with invalid rate",
			rate: "invalid",
			fields: fields{
				Interface: "GigabitEthernet0/0/0/0",
			},
			want: want{
				Err: true,
			},
		},
		{
			name: "Test with valid message",
			rate: "100000",
			fields: fields{
				Interface: "GigabitEthernet0/0/0/0",
			},
			want: want{
				Err: false,
			},
		},
		{
			name: "Test with zero value",
			rate: "0",
			fields: fields{
				Interface: "GigabitEthernet0/0/0/0",
			},
			want: want{
				Err: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := consumer.BandwidthMessage{
				TelemetryMessage: consumer.TelemetryMessage{
					Name: "isis",
					Tags: consumer.MessageTags{
						Host:          "telegraf",
						InterfaceName: tt.fields.Interface,
						Path:          "Cisco-IOS-XR-clns-isis-oper:isis/instances/instance/interfaces/interface",
						Source:        "XR-1",
						Subscription:  "hawk-metrics",
					},
					Timestamp: 1704728135,
				},
				Bandwidth: 0,
			}
			ctrl := gomock.NewController(t)
			config := config.NewMockConfig(ctrl)
			unprocessedMsgChan := make(chan consumer.Message)
			processedMsgChan := make(chan consumer.Message)
			helper := helpers.NewMockHelper(ctrl)
			processor := NewDefaultProcessor(config, unprocessedMsgChan, processedMsgChan, helper)
			impairmentsPrefix := "nodes.XR-1.config.Gi0-0-0-0.impairments."
			helper.EXPECT().GetDefaultImpairmentsPrefix("XR-1", "Gi0-0-0-0").Return(impairmentsPrefix).AnyTimes()
			config.EXPECT().GetValue(impairmentsPrefix + "rate").Return(tt.rate).AnyTimes()
			go processor.processBandwidthMessage(&msg)
			time.Sleep(time.Second * 1)
			if tt.want.Err {
				select {
				case <-processor.processedMsgChan:
					assert.Fail(t, "Message should not be sent to processedMsgChan")
				default:
					assert.NoError(t, nil)
				}
			} else {
				select {
				case <-processor.processedMsgChan:
					assert.NoError(t, nil)
				default:
					assert.Fail(t, "Message should be sent to processedMsgChan")
				}
			}
		})
	}
}

func TestDefaultProcessor_processMessage(t *testing.T) {
	tests := []struct {
		name    string
		msg     consumer.Message
		wantErr bool
	}{
		{
			name: "Test with valid Delay message",
			msg: &consumer.DelayMessage{
				TelemetryMessage: consumer.TelemetryMessage{
					Name: "performance_monitoring",
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
			wantErr: false,
		},
		{
			name: "Test with valid Loss message",
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
				LossPercentage: 0.0,
			},
			wantErr: false,
		},
		{
			name: "Test with valid Bandwidth message",
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
				Bandwidth: 0,
			},
			wantErr: false,
		},
		{
			name: "Test with invalid message",
			msg: &consumer.TelemetryMessage{
				Name: "unknown",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			config := config.NewMockConfig(ctrl)
			unprocessedMsgChan := make(chan consumer.Message)
			processedMsgChan := make(chan consumer.Message)
			helper := helpers.NewMockHelper(ctrl)
			processor := NewDefaultProcessor(config, unprocessedMsgChan, processedMsgChan, helper)
			impairmentsPrefix := "nodes.XR-1.config.Gi0-0-0-0.impairments."
			helper.EXPECT().GetDefaultImpairmentsPrefix("XR-1", "Gi0-0-0-0").Return(impairmentsPrefix).AnyTimes()
			config.EXPECT().GetValue(impairmentsPrefix + "delay").Return("").AnyTimes()
			config.EXPECT().GetValue(impairmentsPrefix + "jitter").Return("").AnyTimes()
			config.EXPECT().GetValue(impairmentsPrefix + "loss").Return("").AnyTimes()
			config.EXPECT().GetValue(impairmentsPrefix + "rate").Return("").AnyTimes()
			go processor.processMessage(tt.msg)
			time.Sleep(time.Second * 1)
			if tt.wantErr {
				select {
				case <-processor.processedMsgChan:
					assert.Fail(t, "Message should not be sent to processedMsgChan")
				default:
					assert.NoError(t, nil)
				}
			} else {
				select {
				case <-processor.processedMsgChan:
					assert.NoError(t, nil)
				default:
					assert.Fail(t, "Message should be sent to processedMsgChan")
				}
			}
		})
	}
}

func TestDefaultProcessor_Stop(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Test Stop",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			config := config.NewMockConfig(ctrl)
			unprocessedMsgChan := make(chan consumer.Message)
			processedMsgChan := make(chan consumer.Message)
			helper := helpers.NewMockHelper(ctrl)
			processor := NewDefaultProcessor(config, unprocessedMsgChan, processedMsgChan, helper)
			go processor.Start()
			time.Sleep(time.Second * 1)
			go processor.Stop()
			assert.NoError(t, nil)
		})
	}
}
