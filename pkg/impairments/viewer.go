package impairments

import (
	"github.com/hawkv6/clab-telemetry-linker/pkg/command"
	"github.com/hawkv6/clab-telemetry-linker/pkg/logging"
)

type Viewer interface {
	ShowImpairments() error
}

type DefaultViewer struct {
	ImpairmentsManager
	command command.ShowCommand
}

func NewDefaultViewer(node string, command command.ShowCommand) *DefaultViewer {
	defautlViewer := &DefaultViewer{
		ImpairmentsManager: ImpairmentsManager{
			log: logging.DefaultLogger.WithField("subsystem", Subsystem),
		},
		command: command,
	}
	return defautlViewer
}

func (manager *DefaultViewer) ShowImpairments() error {
	return manager.command.ShowImpairments()
}
