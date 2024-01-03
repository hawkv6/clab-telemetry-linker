package impairments

import (
	"github.com/hawkv6/clab-telemetry-linker/pkg/command"
	"github.com/hawkv6/clab-telemetry-linker/pkg/config"
	"github.com/hawkv6/clab-telemetry-linker/pkg/helpers"
	"github.com/hawkv6/clab-telemetry-linker/pkg/logging"
	"github.com/sirupsen/logrus"
)

const Subsystem = "impairments"

type ImpairmentsManager interface {
	SetDelay(uint64)
	SetJitter(uint64)
	SetLoss(float64)
	SetRate(uint64)
	ApplyImpairments()
	DeleteImpairments()
	WriteConfig() error
}

type DefaultImpairmentsManager struct {
	log               *logrus.Entry
	config            config.Config
	command           command.SetCommand
	impairmentsPrefix string
}

func NewDefaultImpairmentsManager(config config.Config, node, interface_ string) *DefaultImpairmentsManager {
	defaultImpairmentsManager := &DefaultImpairmentsManager{
		log:               logging.DefaultLogger.WithField("subsystem", Subsystem),
		config:            config,
		impairmentsPrefix: helpers.SetDefaultImpairmentsPrefix(node, interface_),
		command:           command.NewBasicCommand(node, interface_, config.GetValue(helpers.GetDefaultClabNameKey())),
	}
	return defaultImpairmentsManager
}

func (manager *DefaultImpairmentsManager) SetDelay(delay uint64) {
	if delay == 0 {
		manager.log.Debugln("Remove delay from config if set")
		manager.config.DeleteValue(manager.impairmentsPrefix + "delay")
	} else {
		manager.log.Debugf("Set delay in config to %d\n", delay)
		manager.config.SetValue(manager.impairmentsPrefix+"delay", delay)
		manager.command.AddDelay(delay)
	}
}

func (manager *DefaultImpairmentsManager) SetJitter(jitter uint64) {
	if jitter == 0 {
		manager.log.Debugln("Remove jitter from config if set")
		manager.config.DeleteValue(manager.impairmentsPrefix + "jitter")
	} else {
		manager.log.Debugf("Set jitter in config to %d\n", jitter)
		manager.config.SetValue(manager.impairmentsPrefix+"jitter", jitter)
		manager.command.AddJitter(jitter)
	}
}
func (manager *DefaultImpairmentsManager) SetLoss(loss float64) {
	if loss == 0 {
		manager.log.Debugln("Remove loss from config if set")
		manager.config.DeleteValue(manager.impairmentsPrefix + "loss")
	} else {
		manager.log.Debugf("Set loss in config to %f\n", loss)
		manager.config.SetValue(manager.impairmentsPrefix+"loss", loss)
		manager.command.AddLoss(loss)
	}
}

func (manager *DefaultImpairmentsManager) SetRate(rate uint64) {
	if rate == 0 {
		manager.log.Debugln("Remove rate from config if set")
		manager.config.DeleteValue(manager.impairmentsPrefix + "rate")
	} else {
		manager.log.Debugf("Set rate in config to %d\n", rate)
		manager.config.SetValue(manager.impairmentsPrefix+"rate", rate)
		manager.command.AddRate(rate)
	}
}

func (manager *DefaultImpairmentsManager) ApplyImpairments() {
	manager.command.ApplyImpairments()
}

func (manager *DefaultImpairmentsManager) DeleteImpairments() {
	manager.command.DeleteImpairments()
}

func (manager *DefaultImpairmentsManager) WriteConfig() error {
	return manager.config.WriteConfig()
}
