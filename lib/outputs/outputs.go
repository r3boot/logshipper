package outputs

import (
	"fmt"

	"github.com/r3boot/logshipper/lib/config"
	"github.com/r3boot/logshipper/lib/logger"
)

func NewOutputs(l *logger.Logger, c *config.Config) error {
	var err error

	log = l
	cfg = c

	Multiplexer = NewOutputMultiplexer()

	if cfg.Redis.Name != "" {
		Redis, err = NewRedisShipper()
		if err != nil {
			return fmt.Errorf("NewOutputs: %v", err)
		}
		log.Debugf("NewOutputs: Redis log shipper enabled")
	} else {
		log.Debugf("NewOutputs: Redis log shipper not enabled")
	}

	if cfg.Amqp.Name != "" {
		if Amqp, err = NewAmqpShipper(); err != nil {
			return fmt.Errorf("NewOutputs: %v", err)
		}
		log.Debugf("NewOutputs: Amqp log shipper enabled")
	} else {
		log.Debugf("NewOutputs: Amqp log shipper not enabled")
	}

	return nil
}
