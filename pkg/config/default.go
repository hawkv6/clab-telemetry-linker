package config

import (
	"log"
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
	fullFileLocation string
	clabName         string
	clabNameKey      string
}

func (config *DefaultConfig) setUserHome() {
	if err, userHome := helpers.GetUserHome(); err != nil {
		config.log.Fatalln(err)
	} else {
		config.userHome = userHome
	}
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

func (config *DefaultConfig) setClabName(clabName string) {
	name := config.koanfInstance.String(config.clabNameKey)
	if name == "" {
		config.log.Debugln("No clab name found in config, set to default: clab-hawkv6")
		if err := config.koanfInstance.Set(config.clabNameKey, config.clabName); err != nil {
			log.Fatalln(err)
		}
	} else if name != clabName {
		config.log.Debugf("Clab name in config is different from the one provided as flag, use the one from the config: %s", name)
		config.clabName = name
	} else {
		config.log.Debugf("Clab name in config and flag identical: %s", name)
		config.clabName = clabName
	}
}

func (config *DefaultConfig) setConfigPath() {
	config.configPath = config.userHome + "/.clab-telemetry-linker"
}
func (config *DefaultConfig) setConfigFileLocation() {
	config.fullFileLocation = config.configPath + "/" + config.fileName
}

func (config *DefaultConfig) doesConfigExist() bool {
	config.log.Debugln("Check if config file exists:", config.fullFileLocation)
	if _, err := os.Stat(config.fullFileLocation); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		config.log.Fatalf("Unable to check if config file exists: %v", err)
	}
	return true
}

func (config *DefaultConfig) createConfig() {
	config.log.Debugln("Config file not found - creating a new one")
	if err := os.MkdirAll(config.configPath, 0755); err != nil {
		config.log.Fatalln(err)
	}
	configFile, err := os.Create(config.fullFileLocation)
	if err != nil {
		config.log.Fatalln(err)
	}
	defer configFile.Close()
}

func (config *DefaultConfig) readConfig() {
	config.log.Debugln("Read config file: ", config.fullFileLocation)
	if err := config.koanfInstance.Load(file.Provider(config.fullFileLocation), yaml.Parser()); err != nil {
		config.log.Fatalf("Unable to read config file: %v", err)
	}
}

func (config *DefaultConfig) InitConfig() {
	if !config.doesConfigExist() {
		config.createConfig()
	}
	config.readConfig()
}

func (config *DefaultConfig) DeleteValue(key string) {
	config.koanfInstance.Delete(key)
}

func (config *DefaultConfig) SetValue(key string, value interface{}) {
	if err := config.koanfInstance.Set(key, value); err != nil {
		log.Fatalln(err)
	}
}

func (config *DefaultConfig) GetValue(key string) string {
	return config.koanfInstance.String(key)
}

func (config *DefaultConfig) WriteConfig() error {
	config.log.Debugln("Write config file: ", config.fullFileLocation)
	data, err := config.koanfInstance.Marshal(yaml.Parser())
	if err != nil {
		config.log.Errorf("error marshalling config: %v", err)
		return err
	}
	if err := os.WriteFile(config.fullFileLocation, data, 0644); err != nil {
		config.log.Errorf("error writing config: %v", err)
		return err
	}
	return nil
}

func createDefaultConfig(configFileName, clabName, clabNameKey string) *DefaultConfig {
	defaultConfig := &DefaultConfig{
		log:           logging.DefaultLogger.WithField("subsystem", Subsystem),
		koanfInstance: koanf.New("."),
		clabName:      helpers.GetDefaultClabName(),
		clabNameKey:   helpers.GetDefaultClabNameKey(),
	}
	defaultConfig.setUserHome()
	defaultConfig.setConfigFileName(configFileName)
	defaultConfig.setConfigPath()
	defaultConfig.setConfigFileLocation()
	defaultConfig.setClabName(clabName)
	defaultConfig.InitConfig()
	return defaultConfig
}

func NewDefaultConfig() *DefaultConfig {
	return createDefaultConfig("", "", "")
}
