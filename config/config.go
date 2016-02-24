package config

import (
	"errors"
	"github.com/r3boot/rlib/logger"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

const MAX_CFG_SIZE int = 65534

var Log logger.Log

type Config struct {
	Hostname string
	Inputs   []struct {
		Name     string `yaml:"name"`
		Path     string `yaml:"path"`
		Type     string `yaml:"type"`
		TsFormat string `yaml:"ts_format"`
	} `yaml:"inputs"`
	ELK struct {
		Type string `yaml:"type"`
	} `yaml:"elk"`
	Redis struct {
		Name     string `yaml:"name"`
		Uri      string `yaml:"uri"`
		Key      string `yaml:"key"`
		Type     string `yaml:"type"`
		Password string `yaml:"password"`
		Database int64  `yaml:"database"`
	} `yaml:"redis"`
	Amqp struct {
		Name     string `yaml:"name"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Exchange string `yaml:"exchange"`
	} `yaml:"amqp"`
}

func parseTimeFormat(fmt string) (ts_format string) {
	switch fmt {
	case "ANSIC":
		{
			ts_format = time.ANSIC
		}
	case "UnixDate":
		{
			ts_format = time.UnixDate
		}
	case "RubyDate":
		{
			ts_format = time.RubyDate
		}
	case "RFC822":
		{
			ts_format = time.RFC822
		}
	case "RFC822Z":
		{
			ts_format = time.RFC822Z
		}
	case "RFC850":
		{
			ts_format = time.RFC850
		}
	case "RFC1123":
		{
			ts_format = time.RFC1123
		}
	case "RFC1123Z":
		{
			ts_format = time.RFC1123Z
		}
	case "RFC3339":
		{
			ts_format = time.RFC3339
		}
	case "RFC3339Nano":
		{
			ts_format = time.RFC3339Nano
		}
	case "Kitchen":
		{
			ts_format = time.Kitchen
		}
	case "Stamp":
		{
			ts_format = time.Stamp
		}
	case "StampMilli":
		{
			ts_format = time.StampMilli
		}
	case "StampMicro":
		{
			ts_format = time.StampMicro
		}
	case "StampNano":
		{
			ts_format = time.StampNano
		}
	default:
		{
			ts_format = fmt
		}
	}
	return
}

func Setup(l logger.Log) (err error) {
	Log = l

	Log.Verbose("Configuration initialized")

	return
}

func LoadConfig(fname string) (cfg Config, err error) {
	var fd *os.File
	var fs os.FileInfo
	var data []byte

	// Check if the config file exists, and create a buffer to hold it's content
	if fs, err = os.Stat(fname); err != nil {
		return
	}
	data = make([]byte, fs.Size())

	// Open and read the file into the buffer
	if fd, err = os.Open(fname); err != nil {
		return
	}
	defer fd.Close()

	if _, err = fd.Read(data); err != nil {
		return
	}

	// Parse the yaml into a struct
	cfg = Config{}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return
	}

	// Set the hostname
	cfg.Hostname, _ = os.Hostname()

	return
}

func LoadAndCheckConfig(fname string) (cfg Config, err error) {
	if cfg, err = LoadConfig(fname); err != nil {
		return
	}

	// Check if a hostname has been specified
	if cfg.Hostname == "" {
		err = errors.New("Hostname must be set")
		return
	}

	// Check for missing inputs
	if len(cfg.Inputs) == 0 {
		err = errors.New("No inputs configured")
		return
	}

	i := 0
	for _, input := range cfg.Inputs {
		// We need a name for the input
		if input.Name == "" {
			err = errors.New("No name specified for input")
			return
		}

		// We also need a path towards a file we can monitor
		if input.Path == "" {
			err = errors.New("No path specified for " + input.Name)
			return
		}

		// .. and it needs to be readably preferrably, but we can open the
		// file on a later time.
		if _, err = os.Stat(input.Path); err != nil {
			Log.Warning(input.Path + " does not exist")
		}

		// Check if the type of parser to use has been defined
		if input.Type == "" {
			err = errors.New("No type found for " + input.Name)
			return
		}
		valid_input_type := false
		for _, t := range []string{T_SYSLOG, T_CLF, T_SURICATA} {
			if input.Type == t {
				valid_input_type = true
				break
			}
		}
		if !valid_input_type {
			err = errors.New("Unknown type specified for " + input.Name + ": " + input.Type)
			return
		}

		// Check if a timestamp format has been specified. If not, define
		// a default one
		if input.TsFormat == "" {
			Log.Warning("No timestamp format set for " + input.Name + ", defaulting to RFC3339")
			cfg.Inputs[i].TsFormat = time.RFC3339
		} else {
			cfg.Inputs[i].TsFormat = parseTimeFormat(input.TsFormat)
		}

		i += 1
	}

	return
}
