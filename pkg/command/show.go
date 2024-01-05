package command

import (
	"os/exec"

	"github.com/hawkv6/clab-telemetry-linker/pkg/logging"
)

type ShowCommand interface {
	ShowImpairments() error
}

type DefaultShowCommand struct {
	BaseCommand
}

func NewDefaultShowCommand(node, clabName string) *DefaultShowCommand {
	command := &DefaultShowCommand{
		BaseCommand: BaseCommand{
			log:         logging.DefaultLogger.WithField("subsystem", subsystem),
			execCommand: exec.Command("containerlab"),
		},
	}
	command.execCommand.Args = append(command.execCommand.Args, "tools", "netem", "show")
	clabNode := clabName + "-" + node
	command.execCommand.Args = append(command.execCommand.Args, "-n", clabNode)
	command.log.Debugln("Create basic command: ", command.execCommand)
	return command
}

func (command *DefaultShowCommand) ShowImpairments() error {
	return command.ExecuteCommand(command.execCommand)
}
