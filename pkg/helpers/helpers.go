package helpers

import (
	"fmt"
	"os"
	"os/user"
)

type Helper interface {
	IsRoot() bool
	GetUserHome() (error, string)
	GetDefaultClabNameKey() string
	GetDefaultClabName() string
	GetDefaultImpairmentsPrefix(node, interface_ string) string
}

type DefaultHelper struct{}

func NewDefaultHelper() *DefaultHelper {
	return &DefaultHelper{}
}

func (helper *DefaultHelper) IsRoot() bool {
	return os.Geteuid() == 0
}

func (helper *DefaultHelper) GetUserHome() (error, string) {
	username := os.Getenv("SUDO_USER")
	user, err := user.Lookup(username)
	if err != nil {
		return fmt.Errorf("Unable to find userhome for user %q: %v", username, err), ""
	}
	return nil, user.HomeDir
}

func (helper *DefaultHelper) GetDefaultClabNameKey() string {
	return "clab-name"
}

func (helper *DefaultHelper) GetDefaultClabName() string {
	return "clab-hawkv6"
}

func (helper *DefaultHelper) GetDefaultImpairmentsPrefix(node, interface_ string) string {
	return "nodes." + node + ".config." + interface_ + ".impairments."
}
