package outputs

import (
	"github.com/r3boot/logshipper/config"
	"github.com/r3boot/rlib/logger"
)

var Log logger.Log
var Config config.Config

var LogShippers []*LogShipper

func Setup(l logger.Log, c config.Config) (err error) {
	var ls *LogShipper

	Log = l
	Config = c

	for _, output := range Config.Outputs {
		ls, err = NewLogShipper(output.Name, output.Type)
		ls.Redis.Uri = output.Redis.Uri
		ls.Redis.Key = output.Redis.Key
		ls.Redis.Type = output.Redis.Type

		LogShippers = append(LogShippers, ls)
	}
	return
}
