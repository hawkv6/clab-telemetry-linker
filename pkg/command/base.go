package command

import (
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"
)

type Command interface {
	ExecuteCommand(*exec.Cmd) error
}

type BaseCommand struct {
	log         *logrus.Entry
	execCommand *exec.Cmd
}

func (command *BaseCommand) ExecuteCommand(cmd *exec.Cmd) error {
	command.log.Debugf("Execute Command: %s\n", cmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Aborting... Following Error happened: %v", err)
	}
	fmt.Printf("%s\n", output)
	return nil
}
