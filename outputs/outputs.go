package outputs

import (
	"github.com/r3boot/logshipper/config"
	"github.com/r3boot/rlib/logger"
)

var Log logger.Log
var Config config.Config

func Setup(l logger.Log, c config.Config) (err error) {
	Log = l
	Config = c

	return
}
