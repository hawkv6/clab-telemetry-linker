package command

import (
	"fmt"
	"os/exec"

	"github.com/hawkv6/clab-telemetry-linker/pkg/logging"
)

var subsystem = "command"

type SetCommand interface {
	AddDelay(uint64)
	AddJitter(uint64)
	AddLoss(float64)
	AddRate(uint64)
	ApplyImpairments() error
	DeleteImpairments() error
}

type DefaultSetCommand struct {
	BaseCommand
	resetCommand *exec.Cmd
}

func NewDefaultSetCommand(node, interface_, clabName string) *DefaultSetCommand {
	command := &DefaultSetCommand{
		BaseCommand: BaseCommand{
			log:         logging.DefaultLogger.WithField("subsystem", subsystem),
			execCommand: exec.Command("containerlab"),
		},
	}
	command.execCommand.Args = append(command.execCommand.Args, "tools", "netem", "set")
	clabNode := clabName + "-" + node
	command.execCommand.Args = append(command.execCommand.Args, "-n", clabNode)
	command.execCommand.Args = append(command.execCommand.Args, "-i", interface_)
	command.log.Debugln("Create basic command: ", command.execCommand)
	command.resetCommand = command.execCommand
	return command
}

func (command *DefaultSetCommand) AddDelay(delay uint64) {
	if delay != 0 {
		command.log.Debugf("Add '--delay %dms' to command\n", delay)
		command.execCommand.Args = append(command.execCommand.Args, "--delay", fmt.Sprintf("%dms", delay))
	}
}

func (command *DefaultSetCommand) AddJitter(jitter uint64) {
	if jitter != 0 {
		command.log.Debugf("Add '--jitter %dms' to command\n", jitter)
		command.execCommand.Args = append(command.execCommand.Args, "--jitter", fmt.Sprintf("%dms", jitter))
	}
}

func (command *DefaultSetCommand) AddLoss(loss float64) {
	if loss != 0 {
		command.log.Debugf("Add '--loss %f' to command\n", loss)
		command.execCommand.Args = append(command.execCommand.Args, "--loss", fmt.Sprintf("%f", loss))
	}
}

func (command *DefaultSetCommand) AddRate(rate uint64) {
	if rate != 0 {
		command.log.Debugf("Add '--rate %d' to command\n", rate)
		command.execCommand.Args = append(command.execCommand.Args, "--rate", fmt.Sprintf("%d", rate))
	}
}

func (command *DefaultSetCommand) ApplyImpairments() error {
	return command.ExecuteCommand(command.execCommand)
}

func (command *DefaultSetCommand) DeleteImpairments() error {
	return command.ExecuteCommand(command.resetCommand)
}
