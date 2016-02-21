package outputs

import (
	"github.com/r3boot/logshipper/config"
	"github.com/r3boot/rlib/logger"
)

var Log logger.Log
var Config config.Config

var Redis *RedisShipper
var Amqp *AmqpShipper
var Multiplexer *OutputMultiplexer

func Setup(l logger.Log, c config.Config) (err error) {
	Log = l
	Config = c

	if Multiplexer, err = NewOutputMultiplexer(); err != nil {
		return
	}

	if Config.Redis.Name != "" {
		if Redis, err = NewRedisShipper(); err != nil {
			return
		}
		Log.Verbose("Redis log shipper enabled")
	} else {
		Log.Debug("Redis log shipper not enabled")
	}

	if Config.Amqp.Name != "" {
		if Amqp, err = NewAmqpShipper(); err != nil {
			return
		}
	}

	return
}
