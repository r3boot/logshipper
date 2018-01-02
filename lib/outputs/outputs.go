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

	for _, output := range cfg.Outputs {
		switch output.Type {
		case OUTPUT_AMQP:
			{
				if Amqp, err = NewAmqpShipper(); err != nil {
					return fmt.Errorf("NewOutputs: %v", err)
				}
				log.Debugf("NewOutputs: Amqp log shipper enabled")
			}
		case OUTPUT_REDIS:
			{
				Redis, err = NewRedisShipper()
				if err != nil {
					return fmt.Errorf("NewOutputs: %v", err)
				}
				log.Debugf("NewOutputs: Redis log shipper enabled")
			}
		case OUTPUT_ES:
			{

				ES, err = NewESShipper()
				if err != nil {
					return fmt.Errorf("NewOutputs: %v", err)
				}
				log.Debugf("NewOutputs: Elasticsearch log shipper enabled")
			}
		case OUTPUT_STDOUT:
			{
				Stdout = NewStdoutShipper()
				log.Debugf("NewOutputs: Stdout log shipper enabled")
			}
		}
	}

	return nil
}
