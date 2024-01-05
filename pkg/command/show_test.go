package command

import (
	"os/exec"
	"testing"

	"github.com/hawkv6/clab-telemetry-linker/pkg/logging"
	"github.com/stretchr/testify/assert"
)

func TestNewDefaultShowCommand(t *testing.T) {
	type args struct {
		node     string
		clabName string
	}
	tests := []struct {
		name string
		args args
		want *DefaultShowCommand
	}{
		{
			name: "Test create basic show command",
			args: args{
				node:     "XR-1",
				clabName: "clab-hawkv6",
			},
			want: &DefaultShowCommand{
				BaseCommand: BaseCommand{
					log:         logging.DefaultLogger.WithField("subsystem", "command"),
					execCommand: exec.Command("containerlab", "tools", "netem", "show", "-n", "clab-hawkv6-XR-1"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaultShowCommand := NewDefaultShowCommand(tt.args.node, tt.args.clabName)
			assert.Equal(t, tt.want, defaultShowCommand)
		})
	}
}

func TestDefaultShowCommand_ShowImpairments(t *testing.T) {
	type fields struct {
		name     string
		clabName string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Test show impairments",
			fields: fields{
				name:     "non-existing-node",
				clabName: "non-existing-clab",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command := NewDefaultShowCommand(tt.fields.name, tt.fields.clabName)
			if tt.wantErr {
				assert.Error(t, command.ShowImpairments())
			} else {
				assert.NoError(t, command.ShowImpairments())
			}
		})
	}
}
