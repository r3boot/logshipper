package main

import (
	"flag"
	"github.com/r3boot/logshipper/config"
	"github.com/r3boot/logshipper/events"
	"github.com/r3boot/logshipper/inputs"
	"github.com/r3boot/logshipper/outputs"
	"github.com/r3boot/rlib/logger"
	"os"
	"os/signal"
	"strconv"
)

const D_VERBOSE bool = false
const D_DEBUG bool = false
const D_TIMESTAMP bool = false
const D_CFGFILE string = "logshipper.yml"

var verbose = flag.Bool("v", D_VERBOSE, "Enable verbose output")
var debug = flag.Bool("D", D_DEBUG, "Enable debugging output")
var timestamp = flag.Bool("T", D_TIMESTAMP, "Enable timestamps in output")
var cfgfile = flag.String("f", D_CFGFILE, "Configuration file to use")

var Log logger.Log
var Config config.Config

func signalHandler(signals chan os.Signal, done chan bool) {
	for _ = range signals {
		for _, input := range inputs.MonitoredFiles {
			Log.Debug("Sending cleanup signal to handler for " + input.Name)
			input.Control <- config.CMD_CLEANUP
			<-input.Done
		}

		Log.Debug("Sending cleanup signal multiplexer")
		outputs.Multiplexer.Control <- config.CMD_CLEANUP
		<-outputs.Multiplexer.Done

		if outputs.Redis.Name != "" {
			Log.Debug("Sending cleanup signal to " + outputs.Redis.Name)
			outputs.Redis.Control <- config.CMD_CLEANUP
			<-outputs.Redis.Done
		}
		if outputs.Amqp.Name != "" {
			Log.Debug("Sending cleanup signal to " + outputs.Amqp.Name)
			outputs.Amqp.Control <- config.CMD_CLEANUP
			<-outputs.Amqp.Done
		}
		done <- true
	}
}

func init() {
	var err error

	flag.Parse()

	Log.UseDebug = *debug
	Log.UseVerbose = *verbose
	Log.UseTimestamp = *timestamp

	Log.Verbose("Logging initialized")

	if Config, err = config.Setup(Log, *cfgfile); err != nil {
		Log.Fatal("Failed to initialize configuration: " + err.Error())
	}

	if err = events.Setup(Log, Config); err != nil {
		Log.Fatal("Failed to initialize events: " + err.Error())
	}

	if err = inputs.Setup(Log, Config); err != nil {
		Log.Fatal("Failed to initialize inputs: " + err.Error())
	}
	Log.Verbose("Initialized " + strconv.Itoa(len(inputs.MonitoredFiles)) + " inputs")

	if err = outputs.Setup(Log, Config); err != nil {
		Log.Fatal("Failed to initialize outputs: " + err.Error())
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
	Log.Debug("Started output multiplexer")

	for _, input := range inputs.MonitoredFiles {
		go input.Parse(logdata)
	}
	Log.Debug("Started input readers")

	<-cleanupDone
	Log.Debug("Program finished")
}
