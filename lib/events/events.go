package events

import (
	"github.com/r3boot/logshipper/lib/config"
	"github.com/r3boot/logshipper/lib/logger"
)

func NewEvents(l *logger.Logger, c *config.Config) {
	cfg = c
	log = l

	log.Debugf("NewEvents: Events initialized")

	return
}
