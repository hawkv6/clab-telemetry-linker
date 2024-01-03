package config

import (
	"os"

	"github.com/hawkv6/clab-telemetry-linker/pkg/logging"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

const Subsystem = "config"

var (
	log           = logging.DefaultLogger.WithField("subsystem", Subsystem)
	koanfInstance = koanf.New(".")
)

func doesConfigExist() bool {
	log.Debugln("Check if config file exists:", configFile)
	if _, err := os.Stat(configFile); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Fatalln(err)
	}
	return true
}

func createConfig() {
	log.Debugln("Config file not found - creating a new one")
	if err := os.MkdirAll(configPath, 0755); err != nil {
		log.Fatalln(err)
	}
	configFile, err := os.Create(configPath + "/config.yaml")
	if err != nil {
		log.Fatalln(err)
	}
	defer configFile.Close()
}

func readConfig() {
	log.Debugln("Try to read config file: ", configFile)
	if err := koanfInstance.Load(file.Provider(configFile), yaml.Parser()); err != nil {
		log.Fatalln(err)
	}
}

func InitConfig() {
	if !doesConfigExist() {
		createConfig()
	}
	readConfig()
}
