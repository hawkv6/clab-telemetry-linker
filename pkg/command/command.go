package command

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/hawkv6/clab-telemetry-linker/pkg/logging"
	"github.com/sirupsen/logrus"
)

var subsystem = "command"

type Command interface {
	AddDelay(uint64)
	AddJitter(uint64)
	AddLoss(float64)
	AddRate(uint64)
	ApplyImpairments()
	DeleteImpairments()
}

type DefaultCommand struct {
	log          *logrus.Entry
	basicCommand *exec.Cmd
	fullCommand  *exec.Cmd
}

func NewBasicCommand(node, interface_, clabName string) *DefaultCommand {
	command := &DefaultCommand{
		log:         logging.DefaultLogger.WithField("subsystem", subsystem),
		fullCommand: exec.Command("containerlab"),
	}
	command.fullCommand.Args = append(command.fullCommand.Args, "tools", "netem", "set")
	clabNode := clabName + "-" + node
	command.fullCommand.Args = append(command.fullCommand.Args, "-n", clabNode)
	command.fullCommand.Args = append(command.fullCommand.Args, "-i", interface_)
	command.log.Debugln("Create basic command: ", command.fullCommand)
	command.basicCommand = command.fullCommand
	return command
}

func (command *DefaultCommand) AddDelay(delay uint64) {
	if delay != 0 {
		command.log.Debugf("Add '--delay %dms' to command\n", delay)
		command.fullCommand.Args = append(command.fullCommand.Args, "--delay", fmt.Sprintf("%dms", delay))
	}
}

func (command *DefaultCommand) AddJitter(jitter uint64) {
	if jitter != 0 {
		command.log.Debugf("Add '--jitter %dms' to command\n", jitter)
		command.fullCommand.Args = append(command.fullCommand.Args, "--jitter", fmt.Sprintf("%dms", jitter))
	}
}

func (command *DefaultCommand) AddLoss(loss float64) {
	if loss != 0 {
		command.log.Debugf("Add '--loss %f' to command\n", loss)
		command.fullCommand.Args = append(command.fullCommand.Args, "--loss", fmt.Sprintf("%f", loss))
	}
}

func (command *DefaultCommand) AddRate(rate uint64) {
	if rate != 0 {
		command.log.Debugf("Add '--rate %d' to command\n", rate)
		command.fullCommand.Args = append(command.fullCommand.Args, "--rate", fmt.Sprintf("%d", rate))
	}
}

func (command *DefaultCommand) executeCommand(cmd *exec.Cmd) {
	command.log.Debugf("Execute Command: %s\n", cmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Aborting... Following Error happened: %v", err)
	}
	fmt.Printf("%s\n", output)
}

func (command *DefaultCommand) ApplyImpairments() {
	command.executeCommand(command.fullCommand)
}

func (command *DefaultCommand) DeleteImpairments() {
	command.executeCommand(command.basicCommand)
}
