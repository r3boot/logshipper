package inputs

import (
	"github.com/r3boot/logshipper/3rdparty/tail"
	"github.com/r3boot/logshipper/config"
	"github.com/r3boot/rlib/logger"
)

var Config config.Config
var Log logger.Log
var TailConfig tail.Config = tail.Config{Follow: true, MustExist: false}

var MonitoredFiles []*MonitoredFile

func Setup(l logger.Log, c config.Config) (err error) {
	var mf *MonitoredFile

	Log = l
	Config = c

	if err = SetupSyslog(); err != nil {
		return
	}

	if err = SetupCFL(); err != nil {
		return
	}

	for _, input := range Config.Inputs {
		mf, err = NewMonitoredFile(input.Name, input.Path, input.Type)
		if err != nil {
			return
		}
		MonitoredFiles = append(MonitoredFiles, mf)
	}

	return
}
