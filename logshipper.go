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

var Redis *outputs.RedisShipper

func signalHandler(signals chan os.Signal, done chan bool) {
	for _ = range signals {
		for _, input := range inputs.MonitoredFiles {
			Log.Debug("Sending cleanup signal to handler for " + input.Name)
			input.Control <- config.CMD_CLEANUP
			<-input.Done
		}
		if Redis.Name != "" {
			Log.Debug("Sending cleanup signal to handler for " + Redis.Name)
			Redis.Control <- config.CMD_CLEANUP
			<-Redis.Done
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
	if Config.Redis.Name != "" {
		Redis, err = outputs.NewRedisShipper()
		if err != nil {
			Log.Fatal("Failed to initialize redis: " + err.Error())
		}
		Log.Verbose("Initialized redis log shipper")
	}
}

func main() {
	var logdata chan []byte

	logdata = make(chan []byte, 1)

	signals := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)

	signal.Notify(signals, os.Interrupt, os.Kill)
	go signalHandler(signals, cleanupDone)

	go Redis.Ship(logdata)
	Log.Debug("Started redis log shippers")

	for _, input := range inputs.MonitoredFiles {
		go input.Parse(logdata)
	}
	Log.Debug("Started input readers")

	<-cleanupDone
	Log.Debug("Program finished")
}
