package impairments

import (
	"github.com/hawkv6/clab-telemetry-linker/pkg/config"
	"github.com/sirupsen/logrus"
)

const Subsystem = "impairments"

type ImpairmentsManager struct {
	log    *logrus.Entry
	config config.Config
}
