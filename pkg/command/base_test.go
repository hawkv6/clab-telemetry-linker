package command

import (
	"os/exec"
	"testing"

	"github.com/hawkv6/clab-telemetry-linker/pkg/logging"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestBaseCommand_ExecuteCommand(t *testing.T) {
	type fields struct {
		log         *logrus.Entry
		execCommand *exec.Cmd
	}
	type args struct {
		cmd *exec.Cmd
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test valid execute command",
			fields: fields{
				log:         logging.DefaultLogger.WithField("subsystem", "command"),
				execCommand: exec.Command("echo", "test"),
			},
			args: args{
				cmd: exec.Command("echo", "test"),
			},
			wantErr: false,
		},
		{
			name: "Test invalid execute command",
			fields: fields{
				log:         logging.DefaultLogger.WithField("subsystem", "command"),
				execCommand: exec.Command("nocommand"),
			},
			args: args{
				cmd: exec.Command("echo", "test"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command := &BaseCommand{
				log:         tt.fields.log,
				execCommand: tt.fields.execCommand,
			}
			if tt.wantErr {
				assert.Error(t, command.ExecuteCommand(tt.args.cmd))
			} else {
				assert.NoError(t, command.ExecuteCommand(tt.args.cmd))
			}
		})
	}
}
