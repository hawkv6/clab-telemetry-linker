package impairments

import (
	"github.com/hawkv6/clab-telemetry-linker/pkg/command"
	"github.com/hawkv6/clab-telemetry-linker/pkg/config"
	"github.com/hawkv6/clab-telemetry-linker/pkg/helpers"
	"github.com/hawkv6/clab-telemetry-linker/pkg/logging"
)

type Setter interface {
	SetDelay(uint64)
	SetJitter(uint64)
	SetLoss(float64)
	SetRate(uint64)
	ApplyImpairments()
	DeleteImpairments()
	WriteConfig() error
}

type DefaultSetter struct {
	ImpairmentsManager
	command           command.SetCommand
	impairmentsPrefix string
}

func NewDefaultSetter(config config.Config, node, interface_ string, helper helpers.Helper, command command.SetCommand) *DefaultSetter {
	defaultSetter := &DefaultSetter{
		ImpairmentsManager: ImpairmentsManager{
			log:    logging.DefaultLogger.WithField("subsystem", Subsystem),
			config: config,
		},
		impairmentsPrefix: helper.SetDefaultImpairmentsPrefix(node, interface_),
		command:           command,
	}
	return defaultSetter
}

func (manager *DefaultSetter) SetDelay(delay uint64) error {
	if delay == 0 {
		manager.log.Debugln("Remove delay from config if set")
		manager.config.DeleteValue(manager.impairmentsPrefix + "delay")
	} else {
		manager.log.Debugf("Set delay in config to %d\n", delay)
		if err := manager.config.SetValue(manager.impairmentsPrefix+"delay", delay); err != nil {
			return err
		}
		manager.command.AddDelay(delay)
	}
	return nil
}

func (manager *DefaultSetter) SetJitter(jitter uint64) error {
	if jitter == 0 {
		manager.log.Debugln("Remove jitter from config if set")
		manager.config.DeleteValue(manager.impairmentsPrefix + "jitter")
	} else {
		manager.log.Debugf("Set jitter in config to %d\n", jitter)
		if err := manager.config.SetValue(manager.impairmentsPrefix+"jitter", jitter); err != nil {
			return err
		}
		manager.command.AddJitter(jitter)
	}
	return nil
}
func (manager *DefaultSetter) SetLoss(loss float64) error {
	if loss == 0 {
		manager.log.Debugln("Remove loss from config if set")
		manager.config.DeleteValue(manager.impairmentsPrefix + "loss")
	} else {
		manager.log.Debugf("Set loss in config to %f\n", loss)
		if err := manager.config.SetValue(manager.impairmentsPrefix+"loss", loss); err != nil {
			return err
		}
		manager.command.AddLoss(loss)
	}
	return nil
}

func (manager *DefaultSetter) SetRate(rate uint64) error {
	if rate == 0 {
		manager.log.Debugln("Remove rate from config if set")
		manager.config.DeleteValue(manager.impairmentsPrefix + "rate")
	} else {
		manager.log.Debugf("Set rate in config to %d\n", rate)
		if err := manager.config.SetValue(manager.impairmentsPrefix+"rate", rate); err != nil {
			return err
		}
		manager.command.AddRate(rate)
	}
	return nil
}

func (manager *DefaultSetter) ApplyImpairments() error {
	return manager.command.ApplyImpairments()
}

func (manager *DefaultSetter) DeleteImpairments() error {
	return manager.command.DeleteImpairments()
}

func (manager *DefaultSetter) WriteConfig() error {
	return manager.config.WriteConfig()
}
