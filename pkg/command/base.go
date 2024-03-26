package command

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"
)

var subsystem = "command"

type Command interface {
	ExecuteCommand(*exec.Cmd) error
}

type BaseCommand struct {
	log         *logrus.Entry
	execCommand *exec.Cmd
}

func (command *BaseCommand) ExecuteCommand(cmd *exec.Cmd) error {
	command.log.Debugf("Execute Command: %s\n", cmd)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		if stderr.String() != "" {
			return fmt.Errorf("Aborting... Following Error happened: %v", stderr.String())
		} else {
			return fmt.Errorf("Aborting... Following Error happened: %v", err)
		}
	}
	command.log.Debugln("Result: ", out.String())
	return nil
}
