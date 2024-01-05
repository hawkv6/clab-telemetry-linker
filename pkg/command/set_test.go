package command

import (
	"os/exec"
	"testing"

	"github.com/hawkv6/clab-telemetry-linker/pkg/logging"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewDefaultSetCommand(t *testing.T) {
	type args struct {
		node       string
		interface_ string
		clabName   string
	}
	tests := []struct {
		name string
		args args
		want *DefaultSetCommand
	}{
		{
			name: "Test create basic command",
			args: args{
				node:       "XR-1",
				interface_: "Gi0-0-0-0",
				clabName:   "clab-hawkv6",
			},
			want: &DefaultSetCommand{
				BaseCommand: BaseCommand{
					log:         logging.DefaultLogger.WithField("subsystem", "command"),
					execCommand: exec.Command("containerlab", "tools", "netem", "set", "-n", "clab-hawkv6-XR-1", "-i", "Gi0-0-0-0"),
				},
				resetCommand: exec.Command("containerlab", "tools", "netem", "set", "-n", "clab-hawkv6-XR-1", "-i", "Gi0-0-0-0"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewDefaultSetCommand(tt.args.node, tt.args.interface_, tt.args.clabName))
		})
	}
}

func TestDefaultSetCommand_AddDelay(t *testing.T) {
	type fields struct {
		log         *logrus.Entry
		fullCommand *exec.Cmd
	}
	type args struct {
		delay uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *exec.Cmd
	}{
		{
			name: "Test add 100ms delay to command",
			fields: fields{
				log:         logging.DefaultLogger.WithField("subsystem", "command"),
				fullCommand: exec.Command("containerlab", "tools", "netem", "set", "-n", "clab-hawkv6-XR-1", "-i", "Gi0-0-0-0"),
			},
			args: args{
				delay: 100,
			},
			want: exec.Command("containerlab", "tools", "netem", "set", "-n", "clab-hawkv6-XR-1", "-i", "Gi0-0-0-0", "--delay", "100ms"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command := &DefaultSetCommand{
				BaseCommand: BaseCommand{
					log:         tt.fields.log,
					execCommand: tt.fields.fullCommand,
				},
			}
			command.AddDelay(tt.args.delay)
			assert.Equal(t, command.execCommand, tt.want)
		})
	}
}

func TestDefaultSetCommand_AddJitter(t *testing.T) {
	type fields struct {
		log         *logrus.Entry
		fullCommand *exec.Cmd
	}
	type args struct {
		jitter uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *exec.Cmd
	}{
		{
			name: "Test add 100ms jitter to command",
			fields: fields{
				log:         logging.DefaultLogger.WithField("subsystem", "command"),
				fullCommand: exec.Command("containerlab", "tools", "netem", "set", "-n", "clab-hawkv6-XR-1", "-i", "Gi0-0-0-0"),
			},
			args: args{
				jitter: 100,
			},
			want: exec.Command("containerlab", "tools", "netem", "set", "-n", "clab-hawkv6-XR-1", "-i", "Gi0-0-0-0", "--jitter", "100ms"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command := &DefaultSetCommand{
				BaseCommand: BaseCommand{
					log:         tt.fields.log,
					execCommand: tt.fields.fullCommand,
				},
			}
			command.AddJitter(tt.args.jitter)
			assert.Equal(t, command.execCommand, tt.want)
		})
	}
}

func TestDefaultSetCommand_AddLoss(t *testing.T) {
	type fields struct {
		log         *logrus.Entry
		fullCommand *exec.Cmd
	}
	type args struct {
		loss float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *exec.Cmd
	}{
		{
			name: "Test add 100% loss to command",
			fields: fields{
				log:         logging.DefaultLogger.WithField("subsystem", "command"),
				fullCommand: exec.Command("containerlab", "tools", "netem", "set", "-n", "clab-hawkv6-XR-1", "-i", "Gi0-0-0-0"),
			},
			args: args{
				loss: 100,
			},
			want: exec.Command("containerlab", "tools", "netem", "set", "-n", "clab-hawkv6-XR-1", "-i", "Gi0-0-0-0", "--loss", "100.000000"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command := &DefaultSetCommand{
				BaseCommand: BaseCommand{
					log:         tt.fields.log,
					execCommand: tt.fields.fullCommand,
				},
			}
			command.AddLoss(tt.args.loss)
			assert.Equal(t, command.execCommand, tt.want)
		})
	}
}

func TestDefaultSetCommand_AddRate(t *testing.T) {
	type fields struct {
		log         *logrus.Entry
		fullCommand *exec.Cmd
	}
	type args struct {
		rate uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *exec.Cmd
	}{
		{
			name: "Test add 100mbps rate to command",
			fields: fields{
				log:         logging.DefaultLogger.WithField("subsystem", "command"),
				fullCommand: exec.Command("containerlab", "tools", "netem", "set", "-n", "clab-hawkv6-XR-1", "-i", "Gi0-0-0-0"),
			},
			args: args{
				rate: 100000,
			},
			want: exec.Command("containerlab", "tools", "netem", "set", "-n", "clab-hawkv6-XR-1", "-i", "Gi0-0-0-0", "--rate", "100000"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command := &DefaultSetCommand{
				BaseCommand: BaseCommand{
					log:         tt.fields.log,
					execCommand: tt.fields.fullCommand,
				},
			}
			command.AddRate(tt.args.rate)
			assert.Equal(t, command.execCommand, tt.want)
		})
	}
}
func TestDefaultSetCommand_executeCommand(t *testing.T) {
	type fields struct {
		log *logrus.Entry
	}
	type args struct {
		cmd *exec.Cmd
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantError bool
	}{
		{
			name: "Test execute command",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", "command"),
			},
			args: args{
				cmd: exec.Command("echo", "test"),
			},
			wantError: false,
		},
		{
			name: "Test execute command with error",
			fields: fields{
				log: logging.DefaultLogger.WithField("subsystem", "command"),
			},
			args: args{
				cmd: exec.Command("containerlab", "noarg"),
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command := &DefaultSetCommand{
				BaseCommand: BaseCommand{
					log: tt.fields.log,
				},
			}
			if tt.wantError {
				assert.Error(t, command.ExecuteCommand(tt.args.cmd), "exit status 1")
			} else {
				assert.NoError(t, command.ExecuteCommand(tt.args.cmd))
			}
		})
	}
}
