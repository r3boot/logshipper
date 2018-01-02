package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/r3boot/logshipper/lib/logger"
)

func NewConfig(l *logger.Logger, fname string) (*Config, error) {
	log = l

	cfg, err := LoadAndCheckConfig(fname)
	if err != nil {
		return nil, fmt.Errorf("NewConfig: %v", err)
	}

	log.Debugf("NewConfig: Configuration initialized")

	return cfg, nil
}

func LoadConfig(fname string) (*Config, error) {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, fmt.Errorf("LoadConfig ioutil.ReadFile: %v", err)
	}

	// Parse the yaml into a struct
	cfg := &Config{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, fmt.Errorf("LoadConfig yaml.Unmarshal: %v", err)
	}

	// Set the hostname
	cfg.Hostname, _ = os.Hostname()

	return cfg, nil
}

func LoadAndCheckConfig(fname string) (*Config, error) {
	cfg, err := LoadConfig(fname)
	if err != nil {
		return nil, fmt.Errorf("LoadAndCheckConfig: %v", err)
	}

	// Check if a hostname has been specified
	if cfg.Hostname == "" {
		return nil, fmt.Errorf("LoadAndCheckConfig: hostname must be set")
	}

	// Check for missing inputs
	if len(cfg.Inputs) == 0 {
		return nil, fmt.Errorf("LoadAndCheckConfig: no inputs configured")
	}

	i := 0
	for _, input := range cfg.Inputs {
		// Check if the type of parser to use has been defined
		if input.Type == "" {
			return nil, fmt.Errorf("LoadAndCheckConfig: no type found for %s", input.Name)
		}
		valid_input_type := false
		for _, t := range []string{T_SYSLOG, T_CLF, T_SURICATA, T_EXIM, T_AMQPBRIDGE} {
			if input.Type == t {
				valid_input_type = true
				break
			}
		}
		if !valid_input_type {
			err = errors.New("Unknown type specified for " + input.Name + ": " + input.Type)
			return nil, fmt.Errorf("LoadAndCheckConfig: unknown type specified for %s: %s", input.Name, input.Type)
		}

		// We need a name for the input
		if input.Name == "" {
			return nil, fmt.Errorf("LoadAndCheckConfig: no name specified for input")
		}

		// We also need a path towards a file we can monitor
		if input.Path == "" && input.Type != T_AMQPBRIDGE {
			err = errors.New("No path specified for " + input.Name)
			return nil, fmt.Errorf("LoadAndCheckConfig: no path specified for %s", input.Name)
		}

		// .. and it needs to be readably preferrably, but we can open the
		// file on a later time.
		if input.Type != T_AMQPBRIDGE {
			if _, err = os.Stat(input.Path); err != nil {
				log.Warningf("LoadAndCheckConfig: %s does not exist", input.Path)
			}
		}

		// Check if a timestamp format has been specified. If not, define
		// a default one
		if input.TsFormat == "" {
			switch input.Type {
			case T_SYSLOG:
				{
					cfg.Inputs[i].TsFormat = time.RFC3339Nano
				}
			case T_CLF:
				{
					cfg.Inputs[i].TsFormat = TF_CLF
				}
			case T_SURICATA:
				{
					cfg.Inputs[i].TsFormat = TF_SURICATA
				}
			case T_EXIM:
				{
					cfg.Inputs[i].TsFormat = TF_EXIM
				}
			case T_AMQPBRIDGE:
				{
					cfg.Inputs[i].TsFormat = time.RFC3339
				}
			default:
				{
					log.Warningf("LoadAndCheckConfig: no timestamp format set for " + input.Name + ", defaulting to RFC3339")
					cfg.Inputs[i].TsFormat = time.RFC3339
				}
			}
		} else {
			cfg.Inputs[i].TsFormat = parseTimeFormat(input.TsFormat)
		}

		i += 1
	}

	return cfg, nil
}
