package config

import (
	"os"
	"os/user"
)

func checkIsRoot() {
	if os.Geteuid() != 0 {
		log.Fatalln("Hawkwing must be run as root")
	}
}
func getUserHome() string {
	checkIsRoot()
	username := os.Getenv("SUDO_USER")
	user, err := user.Lookup(username)
	if err != nil {
		log.Fatalf("Unable to find userhome for user %q: %v", username, err)
	}
	return user.HomeDir
}

const configName = "config"
const configType = "yaml"
const clabName = "clab-hawkv6"
const clabNameKey = "clab-name"

var configPath = getUserHome() + "/.clab-telemetry-linker"
var configFile = configPath + "/" + configName + "." + configType
