package outputs

import (
	"github.com/r3boot/logshipper/lib/config"
	"github.com/r3boot/logshipper/lib/logger"
)

var (
	cfg *config.Config
	log *logger.Logger

	Redis       *RedisShipper
	Amqp        *AmqpShipper
	Multiplexer *OutputMultiplexer
)
