package inputs

import (
	"tail"

	"github.com/r3boot/logshipper/lib/config"
	"github.com/r3boot/logshipper/lib/logger"
)

var (
	TailConfig = tail.Config{Follow: true, MustExist: false}

	cfg *config.Config
	log *logger.Logger

	MonitoredFiles []*MonitoredFile
	InputAmqp      []*AmqpSlurper
)
