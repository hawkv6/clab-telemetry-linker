package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/hawkv6/clab-telemetry-linker/pkg/helpers"
	"github.com/hawkv6/clab-telemetry-linker/pkg/logging"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/sirupsen/logrus"
)

const Subsystem = "config"

type DefaultConfig struct {
	log              *logrus.Entry
	koanfInstance    *koanf.Koanf
	userHome         string
	fileName         string
	configPath       string
	fullfileLocation string
	clabName         string
	clabNameKey      string
	fileProvider     *file.File
	helper           helpers.Helper
}

func (config *DefaultConfig) setUserHome() error {
	if err, userHome := config.helper.GetUserHome(); err != nil {
		return err
	} else {
		config.userHome = userHome
	}
	return nil
}

func (config *DefaultConfig) setConfigFileName(configName string) {
	if configName == "" {
		config.fileName = "config.yaml"
	} else {
		if !strings.HasSuffix(configName, ".yaml") {
			configName = configName + ".yaml"
		}
		config.fileName = configName
	}
}

func (config *DefaultConfig) GetConfigPath() string {
	return config.fullfileLocation
}

func (config *DefaultConfig) setConfigPath() {
	config.configPath = config.userHome + "/.clab-telemetry-linker"
}
func (config *DefaultConfig) setConfigFileLocation() {
	config.fullfileLocation = config.configPath + "/" + config.fileName
}

func (config *DefaultConfig) setClabName(clabName string) error {
	name := config.koanfInstance.String(config.clabNameKey)
	if name == "" {
		config.log.Debugln("No clab name found in config, set to default: clab-hawkv6")
		config.clabName = config.helper.GetDefaultClabName()
		if err := config.koanfInstance.Set(config.clabNameKey, config.clabName); err != nil {
			return err
		}
	} else if name != clabName {
		config.log.Debugf("Clab name in config is different from the one provided as flag, use the one from the config: %s", name)
		config.clabName = name
	} else {
		config.log.Debugf("Clab name in config and flag identical: %s", name)
		config.clabName = clabName
	}
	return nil
}

func (config *DefaultConfig) doesConfigExist() (error, bool) {
	config.log.Debugln("Check if config file exists:", config.fullfileLocation)
	if _, err := os.Stat(config.fullfileLocation); err != nil {
		if os.IsNotExist(err) {
			return nil, false
		}
		return fmt.Errorf("Unable to check if config file exists: %v", err), false
	}
	return nil, true
}

func (config *DefaultConfig) createConfig() error {
	config.log.Debugln("Config file not found - creating a new one")
	if err := os.MkdirAll(config.configPath, 0755); err != nil {
		return err
	}
	configFile, err := os.Create(config.fullfileLocation)
	if err != nil {
		return err
	}
	defer configFile.Close()
	return nil
}

func (config *DefaultConfig) readConfig() error {
	config.log.Infoln("Read config file: ", config.fullfileLocation)
	config.fileProvider = file.Provider(config.fullfileLocation)
	if err := config.koanfInstance.Load(config.fileProvider, yaml.Parser()); err != nil {
		return err
	}
	return nil
}

func (config *DefaultConfig) InitConfig() error {
	err, exist := config.doesConfigExist()
	if err != nil {
		return err
	}
	if !exist {
		if err := config.createConfig(); err != nil {
			return err
		}
	}
	if err := config.readConfig(); err != nil {
		return err
	}
	return nil
}
func (config *DefaultConfig) WatchConfigChange() error {
	if err := config.fileProvider.Watch(func(event interface{}, err error) {

		if err != nil {
			config.log.Errorf("Error watching config file: %v", err)
		}
		config.log.Debugln("Config file changed")
		if err := config.readConfig(); err != nil {
			config.log.Errorf("Error reading config file: %v", err)
		}
	}); err != nil {
		return err
	}
	return nil
}

func (config *DefaultConfig) DeleteValue(key string) {
	config.log.Debugln("Delete value from config: ", key)
	config.koanfInstance.Delete(key)
}

func (config *DefaultConfig) SetValue(key string, value interface{}) error {
	config.log.Debugln("Set value in config: ", key, value)
	if err := config.koanfInstance.Set(key, value); err != nil {
		return err
	}
	return nil
}

func (config *DefaultConfig) GetValue(key string) string {
	value := config.koanfInstance.String(key)
	if value == "" {
		config.log.Debugf("No value found in config for key: %s", key)
	} else {
		config.log.Debugf("value from config: %s = %s", key, value)
	}
	return value
}

func (config *DefaultConfig) WriteConfig() error {
	config.log.Debugln("Write config file: ", config.fullfileLocation)
	data, err := config.koanfInstance.Marshal(yaml.Parser())
	if err != nil {
		config.log.Errorf("error marshalling config: %v", err)
		return err
	}
	if err := os.WriteFile(config.fullfileLocation, data, 0644); err != nil {
		config.log.Errorf("error writing config: %v", err)
		return err
	}
	return nil
}

func CreateDefaultConfig(configFileName, clabName, clabNameKey string, helper helpers.Helper) (error, *DefaultConfig) {
	defaultConfig := &DefaultConfig{
		log:           logging.DefaultLogger.WithField("subsystem", Subsystem),
		koanfInstance: koanf.New("."),
		clabNameKey:   helper.GetDefaultClabNameKey(),
		helper:        helper,
	}
	if err := defaultConfig.setUserHome(); err != nil {
		return err, nil
	}
	defaultConfig.setConfigFileName(configFileName)
	defaultConfig.setConfigPath()
	defaultConfig.setConfigFileLocation()
	if err := defaultConfig.setClabName(clabName); err != nil {
		return err, nil
	}
	if err := defaultConfig.InitConfig(); err != nil {
		return err, nil
	}
	return nil, defaultConfig
}

func NewDefaultConfig() (error, *DefaultConfig) {
	err, defaultConfig := CreateDefaultConfig("", "", "", helpers.NewDefaultHelper())
	if err != nil {
		return err, nil
	}
	return nil, defaultConfig
}
