package command

import (
	"fmt"
	"os/exec"

	"github.com/hawkv6/clab-telemetry-linker/pkg/config"
	"github.com/hawkv6/clab-telemetry-linker/pkg/logging"
)

var basicCommand = exec.Command("containerlab")
var log = logging.DefaultLogger.WithField("subsystem", "command")

func CreateBasicCommand(node string, interface_ string) {
	basicCommand.Args = append(basicCommand.Args, "tools", "netem", "set")
	clabNode := config.GetClabName() + "-" + node
	basicCommand.Args = append(basicCommand.Args, "-n", clabNode)
	basicCommand.Args = append(basicCommand.Args, "-i", interface_)
	log.Debugln("Create basic command: ", basicCommand)
}

func AddDelay(delay uint64) {
	if delay != 0 {
		log.Debugf("Add '--delay %dms' to command\n", delay)
		basicCommand.Args = append(basicCommand.Args, "--delay", fmt.Sprintf("%dms", delay))
	}
}

func AddJitter(jitter uint64) {
	if jitter != 0 {
		log.Debugf("Add '--jitter %dms' to command\n", jitter)
		basicCommand.Args = append(basicCommand.Args, "--jitter", fmt.Sprintf("%dms", jitter))
	}
}

func AddLoss(loss float64) {
	if loss != 0 {
		log.Debugf("Add '--loss %f' to command\n", loss)
		basicCommand.Args = append(basicCommand.Args, "--loss", fmt.Sprintf("%f", loss))
	}
}

func AddRate(rate uint64) {
	if rate != 0 {
		log.Debugf("Add '--rate %d' to command\n", rate)
		basicCommand.Args = append(basicCommand.Args, "--rate", fmt.Sprintf("%d", rate))
	}
}

func ExecuteCommand() {
	log.Debugf("Execute Command: %s\n", basicCommand)
	output, err := basicCommand.CombinedOutput()
	if err != nil {
		log.Fatalf("Error executing command: %v\n", err)
	}
	fmt.Printf("%s\n", output)
}
