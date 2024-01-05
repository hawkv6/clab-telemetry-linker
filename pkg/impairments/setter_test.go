package impairments

import (
	"errors"
	"testing"

	"github.com/hawkv6/clab-telemetry-linker/pkg/command"
	"github.com/hawkv6/clab-telemetry-linker/pkg/config"
	"github.com/hawkv6/clab-telemetry-linker/pkg/helpers"
	"github.com/hawkv6/clab-telemetry-linker/pkg/logging"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewDefaultSetter(t *testing.T) {
	type args struct {
		node       string
		interface_ string
		helper     helpers.Helper
	}
	tests := []struct {
		name string
		args args
		want *DefaultSetter
	}{
		{
			name: "TestNewDefaultImpairmentsManager",
			args: args{
				node:       "XR-1",
				interface_: "Gi0-0-0-0",
				helper:     helpers.NewDefaultHelper(),
			},
			want: &DefaultSetter{
				ImpairmentsManager: ImpairmentsManager{
					log: logging.DefaultLogger.WithField("subsystem", Subsystem),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			config := config.NewMockConfig(ctrl)
			helper := helpers.NewDefaultHelper()
			config.EXPECT().GetValue(helper.GetDefaultClabNameKey()).Return("clab-name").AnyTimes()
			tt.want.config = config
			tt.want.impairmentsPrefix = helper.SetDefaultImpairmentsPrefix(tt.args.node, tt.args.interface_)
			tt.want.command = command.NewDefaultSetCommand(tt.args.node, tt.args.interface_, config.GetValue(helper.GetDefaultClabNameKey()))
			assert.Equal(t, tt.want, NewDefaultSetter(config, tt.args.node, tt.args.interface_, tt.args.helper, tt.want.command))

		})
	}
}

func TestDefaultSetter_SetDelay(t *testing.T) {
	type args struct {
		delay uint64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test with positive delay",
			args: args{
				delay: 100,
			},
			wantErr: false,
		},
		{
			name: "Test with return error from config",
			args: args{
				delay: 100,
			},
			wantErr: true,
		},
		{
			name: "Test with delay 0 (delete delay)",
			args: args{
				delay: 0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockConfig := config.NewMockConfig(gomock.NewController(t))
			mockCommand := command.NewMockSetCommand(gomock.NewController(t))
			manager := &DefaultSetter{
				ImpairmentsManager: ImpairmentsManager{
					log:    logging.DefaultLogger.WithField("subsystem", Subsystem),
					config: mockConfig,
				},
				command:           mockCommand,
				impairmentsPrefix: "nodes.XR-1.config.Gi0-0-0-0.impairments.",
			}
			if tt.wantErr {
				mockConfig.EXPECT().SetValue(manager.impairmentsPrefix+"delay", tt.args.delay).Return(errors.New("error"))
				assert.Error(t, manager.SetDelay(tt.args.delay))
			} else {
				if tt.args.delay == 0 {
					mockConfig.EXPECT().DeleteValue(manager.impairmentsPrefix + "delay").Return()
					assert.NoError(t, manager.SetDelay(tt.args.delay))
				} else {
					mockConfig.EXPECT().SetValue(manager.impairmentsPrefix+"delay", tt.args.delay).Return(nil)
					mockCommand.EXPECT().AddDelay(tt.args.delay)
					assert.NoError(t, manager.SetDelay(tt.args.delay))
				}
			}
		})
	}
}

func TestDefaultSetter_SetJitter(t *testing.T) {
	type args struct {
		jitter uint64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test with positive jitter",
			args: args{
				jitter: 100,
			},
			wantErr: false,
		},
		{
			name: "Test with return error from config",
			args: args{
				jitter: 100,
			},
			wantErr: true,
		},
		{
			name: "Test with jitter 0 (delete jitter)",
			args: args{
				jitter: 0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockConfig := config.NewMockConfig(gomock.NewController(t))
			mockCommand := command.NewMockSetCommand(gomock.NewController(t))
			manager := &DefaultSetter{
				ImpairmentsManager: ImpairmentsManager{
					log:    logging.DefaultLogger.WithField("subsystem", Subsystem),
					config: mockConfig,
				},
				command:           mockCommand,
				impairmentsPrefix: "nodes.XR-1.config.Gi0-0-0-0.impairments.",
			}
			if tt.wantErr {
				mockConfig.EXPECT().SetValue(manager.impairmentsPrefix+"jitter", tt.args.jitter).Return(errors.New("error"))
				assert.Error(t, manager.SetJitter(tt.args.jitter))
			} else {
				if tt.args.jitter == 0 {
					mockConfig.EXPECT().DeleteValue(manager.impairmentsPrefix + "jitter").Return()
					assert.NoError(t, manager.SetJitter(tt.args.jitter))
				} else {
					mockConfig.EXPECT().SetValue(manager.impairmentsPrefix+"jitter", tt.args.jitter).Return(nil)
					mockCommand.EXPECT().AddJitter(tt.args.jitter)
					assert.NoError(t, manager.SetJitter(tt.args.jitter))
				}
			}
		})
	}
}

func TestDefaultSetter_SetLoss(t *testing.T) {
	type args struct {
		loss float64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test with positive loss",
			args: args{
				loss: 100,
			},
			wantErr: false,
		},
		{
			name: "Test with return error from config",
			args: args{
				loss: 100,
			},
			wantErr: true,
		},
		{
			name: "Test with loss 0 (delete loss)",
			args: args{
				loss: 0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockConfig := config.NewMockConfig(gomock.NewController(t))
			mockCommand := command.NewMockSetCommand(gomock.NewController(t))
			manager := &DefaultSetter{
				ImpairmentsManager: ImpairmentsManager{
					log:    logging.DefaultLogger.WithField("subsystem", Subsystem),
					config: mockConfig,
				},
				command:           mockCommand,
				impairmentsPrefix: "nodes.XR-1.config.Gi0-0-0-0.impairments.",
			}
			if tt.wantErr {
				mockConfig.EXPECT().SetValue(manager.impairmentsPrefix+"loss", tt.args.loss).Return(errors.New("error"))
				assert.Error(t, manager.SetLoss(tt.args.loss))
			} else {
				if tt.args.loss == 0 {
					mockConfig.EXPECT().DeleteValue(manager.impairmentsPrefix + "loss").Return()
					assert.NoError(t, manager.SetLoss(tt.args.loss))
				} else {
					mockConfig.EXPECT().SetValue(manager.impairmentsPrefix+"loss", tt.args.loss).Return(nil)
					mockCommand.EXPECT().AddLoss(tt.args.loss)
					assert.NoError(t, manager.SetLoss(tt.args.loss))
				}
			}
		})
	}
}

func TestDefaultSetter_SetRate(t *testing.T) {
	type args struct {
		rate uint64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test with positive rate",
			args: args{
				rate: 100000,
			},
			wantErr: false,
		},
		{
			name: "Test with return error from config",
			args: args{
				rate: 100000,
			},
			wantErr: true,
		},
		{
			name: "Test with rate 0 (delete rate)",
			args: args{
				rate: 0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockConfig := config.NewMockConfig(gomock.NewController(t))
			mockCommand := command.NewMockSetCommand(gomock.NewController(t))
			manager := &DefaultSetter{
				ImpairmentsManager: ImpairmentsManager{
					log:    logging.DefaultLogger.WithField("subsystem", Subsystem),
					config: mockConfig,
				},
				command:           mockCommand,
				impairmentsPrefix: "nodes.XR-1.config.Gi0-0-0-0.impairments.",
			}
			if tt.wantErr {
				mockConfig.EXPECT().SetValue(manager.impairmentsPrefix+"rate", tt.args.rate).Return(errors.New("error"))
				assert.Error(t, manager.SetRate(tt.args.rate))
			} else {
				if tt.args.rate == 0 {
					mockConfig.EXPECT().DeleteValue(manager.impairmentsPrefix + "rate").Return()
					assert.NoError(t, manager.SetRate(tt.args.rate))
				} else {
					mockConfig.EXPECT().SetValue(manager.impairmentsPrefix+"rate", tt.args.rate).Return(nil)
					mockCommand.EXPECT().AddRate(tt.args.rate)
					assert.NoError(t, manager.SetRate(tt.args.rate))
				}
			}
		})
	}
}

func TestDefaultSetter_ApplyImpairments(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Test without error",
			wantErr: false,
		},
		{
			name:    "Test with error",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockConfig := config.NewMockConfig(gomock.NewController(t))
			mockCommand := command.NewMockSetCommand(gomock.NewController(t))
			manager := &DefaultSetter{
				ImpairmentsManager: ImpairmentsManager{
					log:    logging.DefaultLogger.WithField("subsystem", Subsystem),
					config: mockConfig,
				},
				command:           mockCommand,
				impairmentsPrefix: "nodes.XR-1.config.Gi0-0-0-0.impairments.",
			}
			if tt.wantErr {
				mockCommand.EXPECT().ApplyImpairments().Return(errors.New("error"))
				assert.Error(t, manager.ApplyImpairments())
			} else {
				mockCommand.EXPECT().ApplyImpairments().Return(nil)
				assert.NoError(t, manager.ApplyImpairments())
			}
		})
	}

}

func TestDefaultSetter_DeleteImpairments(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Test without error",
			wantErr: false,
		},
		{
			name:    "Test with error",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockConfig := config.NewMockConfig(gomock.NewController(t))
			mockCommand := command.NewMockSetCommand(gomock.NewController(t))
			manager := &DefaultSetter{
				ImpairmentsManager: ImpairmentsManager{
					log:    logging.DefaultLogger.WithField("subsystem", Subsystem),
					config: mockConfig,
				},
				command:           mockCommand,
				impairmentsPrefix: "nodes.XR-1.config.Gi0-0-0-0.impairments.",
			}
			if tt.wantErr {
				mockCommand.EXPECT().DeleteImpairments().Return(errors.New("error"))
				assert.Error(t, manager.DeleteImpairments())
			} else {
				mockCommand.EXPECT().DeleteImpairments().Return(nil)
				assert.NoError(t, manager.DeleteImpairments())
			}
		})
	}
}

func TestDefaultSetter_WriteConfig(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Test without error",
			wantErr: false,
		},
		{
			name:    "Test with error",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockConfig := config.NewMockConfig(gomock.NewController(t))
			mockCommand := command.NewMockSetCommand(gomock.NewController(t))
			manager := &DefaultSetter{
				ImpairmentsManager: ImpairmentsManager{
					log:    logging.DefaultLogger.WithField("subsystem", Subsystem),
					config: mockConfig,
				},
				command:           mockCommand,
				impairmentsPrefix: "nodes.XR-1.config.Gi0-0-0-0.impairments.",
			}
			if tt.wantErr {
				mockConfig.EXPECT().WriteConfig().Return(errors.New("error"))
				assert.Error(t, manager.WriteConfig())
			} else {
				mockConfig.EXPECT().WriteConfig().Return(nil)
				assert.NoError(t, manager.WriteConfig())
			}
		})
	}
}
