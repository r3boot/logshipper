package outputs

import (
	"github.com/r3boot/logshipper/lib/config"
	"github.com/r3boot/logshipper/lib/logger"
)

const (
	OUTPUT_REDIS  = "redis"
	OUTPUT_AMQP   = "amqp"
	OUTPUT_ES     = "elasticsearch"
	OUTPUT_STDOUT = "stdout"
)

var (
	cfg *config.Config
	log *logger.Logger

	Redis       *RedisShipper
	Amqp        *AmqpShipper
	ES          *ESShipper
	Stdout      *StdoutShipper
	Multiplexer *OutputMultiplexer
)
