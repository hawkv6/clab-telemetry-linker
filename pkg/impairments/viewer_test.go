package impairments

import (
	"fmt"
	"testing"

	"github.com/hawkv6/clab-telemetry-linker/pkg/command"
	"github.com/hawkv6/clab-telemetry-linker/pkg/logging"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewDefaultViewer(t *testing.T) {
	type args struct {
		node     string
		clabName string
	}
	tests := []struct {
		name string
		args args
		want *DefaultViewer
	}{
		{
			name: "Test function TestNewDefaultViewer",
			args: args{
				node:     "XR-1",
				clabName: "clab-hawkv6",
			},
			want: &DefaultViewer{
				ImpairmentsManager: ImpairmentsManager{
					log: logging.DefaultLogger.WithField("subsystem", Subsystem),
				},
				command: command.NewDefaultShowCommand("XR-1", "clab-hawkv6"),
			},
		},
	}
	for _, tt := range tests {
		command := command.NewDefaultShowCommand(tt.args.node, tt.args.clabName)
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewDefaultViewer(tt.args.node, command))
		})
	}
}

func TestDefaultViewer_ShowImpairments(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Test ShowImpairments with no errors",
			wantErr: false,
		},
		{
			name:    "Test ShowImpairments with no errors",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &DefaultViewer{
				ImpairmentsManager: ImpairmentsManager{
					log: logging.DefaultLogger.WithField("subsystem", Subsystem),
				},
			}
			ctrl := gomock.NewController(t)
			command := command.NewMockShowCommand(ctrl)
			manager.command = command
			if tt.wantErr {
				command.EXPECT().ShowImpairments().Return(fmt.Errorf("error showing impairments"))
				assert.Error(t, manager.ShowImpairments())
			} else {
				command.EXPECT().ShowImpairments().Return(nil)
				assert.NoError(t, manager.ShowImpairments())
			}
		})
	}
}
