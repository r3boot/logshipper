package events

import (
	"github.com/r3boot/logshipper/config"
	"github.com/r3boot/rlib/logger"
)

var Config config.Config
var Log logger.Log

func Setup(l logger.Log, c config.Config) (err error) {
	Config = c
	Log = l

	Log.Verbose("Events initialized")

	return
}
