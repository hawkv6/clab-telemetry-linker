package config

import (
	"os"

	"github.com/knadh/koanf/parsers/yaml"
)

var impairmentsPrefix string

func GetClabName() string {
	name := koanfInstance.String(clabNameKey)
	if name == "" {
		if err := koanfInstance.Set(clabNameKey, clabName); err != nil {
			log.Fatalln(err)
		}
	}
	return clabName
}

func SetPrefix(node, interface_ string) {
	impairmentsPrefix = "nodes." + node + ".interfaces." + interface_ + ".impairments."
}

func SetDelay(delay uint64) {
	if delay == 0 {
		log.Debugln("Remove delay from config if set")
		koanfInstance.Delete(impairmentsPrefix + "delay")
	} else {
		log.Debugf("Set delay in config to %d\n", delay)
		if err := koanfInstance.Set(impairmentsPrefix+"delay", delay); err != nil {
			log.Fatalln(err)
		}
	}
}

func SetJitter(jitter uint64) {
	if jitter == 0 {
		log.Debugln("Remove jitter from config if set")
		koanfInstance.Delete(impairmentsPrefix + "jitter")
	} else {
		log.Debugf("Set jitter in config to %d\n", jitter)
		if err := koanfInstance.Set(impairmentsPrefix+"jitter", jitter); err != nil {
			log.Fatalln(err)
		}
	}
}
func SetLoss(loss float64) {
	if loss == 0 {
		log.Debugln("Remove loss from config if set")
		koanfInstance.Delete(impairmentsPrefix + "loss")
	} else {
		log.Debugf("Set loss in config to %f\n", loss)
		if err := koanfInstance.Set(impairmentsPrefix+"loss", loss); err != nil {
			log.Fatalln(err)
		}
	}
}

func SetRate(rate uint64) {
	if rate == 0 {
		log.Debugln("Remove rate from config if set")
		koanfInstance.Delete(impairmentsPrefix + "rate")
	} else {
		log.Debugf("Set rate in config to %d\n", rate)
		if err := koanfInstance.Set(impairmentsPrefix+"rate", rate); err != nil {
			log.Fatalln(err)
		}
	}
}

func WriteConfig() error {
	log.Debugln("Try to write config file: ", configFile)
	data, err := koanfInstance.Marshal(yaml.Parser())
	if err != nil {
		log.Errorf("error marshalling config: %v", err)
		return err
	}
	if err := os.WriteFile(configFile, data, 0644); err != nil {
		log.Errorf("error writing config file: %v", err)
		return err
	}
	return nil
}
