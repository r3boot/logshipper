package logshipper

import (
	"flag"
	"os"
	"os/signal"

	"github.com/r3boot/logshipper/lib/config"
	"github.com/r3boot/logshipper/lib/events"
	"github.com/r3boot/logshipper/lib/inputs"
	"github.com/r3boot/logshipper/lib/logger"
	"github.com/r3boot/logshipper/lib/outputs"
)

const (
	D_VERBOSE    = false
	D_DEBUG      = false
	D_TIMESTAMP  = false
	D_CFGFILE    = "logshipper.yml"
	D_CONFIGTEST = false
)

var (
	verbose   = flag.Bool("v", D_VERBOSE, "Enable verbose output")
	debug     = flag.Bool("D", D_DEBUG, "Enable debugging output")
	timestamp = flag.Bool("T", D_TIMESTAMP, "Enable timestamps in output")
	cfgfile   = flag.String("f", D_CFGFILE, "Configuration file to use")

	Logger *logger.Logger
	Config *config.Config
)

func signalHandler(signals chan os.Signal, done chan bool) {
	for _ = range signals {
		for _, input := range inputs.MonitoredFiles {
			Logger.Debugf("main: Sending cleanup signal to handler for " + input.Name)
			input.Control <- config.CMD_CLEANUP
			<-input.Done
		}

		Logger.Debugf("main: Sending cleanup signal multiplexer")
		outputs.Multiplexer.Control <- config.CMD_CLEANUP
		<-outputs.Multiplexer.Done

		if outputs.Redis != nil {
			Logger.Debugf("main: Sending cleanup signal to " + outputs.Redis.Name)
			outputs.Redis.Control <- config.CMD_CLEANUP
			<-outputs.Redis.Done
		}
		if outputs.Amqp != nil {
			Logger.Debugf("main: Sending cleanup signal to " + outputs.Amqp.Name)
			outputs.Amqp.Control <- config.CMD_CLEANUP
			<-outputs.Amqp.Done
		}
		done <- true
	}
}

func init() {
	var err error

	flag.Parse()

	Logger = logger.NewLogger(*timestamp, *debug)
	Logger.Debugf("Logging initialized")

	Config, err = config.NewConfig(Logger, *cfgfile)
	if err != nil {
		Logger.Fatalf("init: %v", err)
	}

	events.NewEvents(Logger, Config)

	err = inputs.NewInputs(Logger, Config)
	if err != nil {
		Logger.Fatalf("init: %v", err)
	}
	Logger.Debugf("init: Initialized %d inputs", len(inputs.MonitoredFiles))

	err = outputs.NewOutputs(Logger, Config)
	if err != nil {
		Logger.Fatalf("init: %v", err)
	}
}

func main() {
	var logdata chan []byte

	logdata = make(chan []byte, 1)

	signals := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)

	signal.Notify(signals, os.Interrupt, os.Kill)
	go signalHandler(signals, cleanupDone)

	go outputs.Multiplexer.Run(logdata)
	Logger.Debugf("main: Started output multiplexer")

	for _, input := range inputs.MonitoredFiles {
		go input.Parse(logdata)
	}
	Logger.Debugf("main: Started input readers")

	<-cleanupDone
	Logger.Debugf("main: Program finished")
}
