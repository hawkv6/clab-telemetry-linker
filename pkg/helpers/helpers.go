package helpers

import (
	"fmt"
	"os"
	"os/user"
)

func IsRoot() bool {
	return os.Geteuid() == 0
}

func GetUserHome() (error, string) {
	username := os.Getenv("SUDO_USER")
	user, err := user.Lookup(username)
	if err != nil {
		return fmt.Errorf("Unable to find userhome for user %q: %v", username, err), ""
	}
	return nil, user.HomeDir
}

func GetDefaultClabNameKey() string {
	return "clab-name"
}

func GetDefaultClabName() string {
	return "clab-hawkv6"
}

func SetDefaultImpairmentsPrefix(node, interface_ string) string {
	return "nodes." + node + ".config." + interface_ + ".impairments."
}
