package config

import (
	"github.com/r3boot/rlib/logger"
	"gopkg.in/yaml.v2"
	"os"
)

const MAX_CFG_SIZE int = 65534

var Log logger.Log

type Config struct {
	Hostname string
	Inputs   []struct {
		Name string `yaml:"name"`
		Path string `yaml:"path"`
		Type int    `yaml:"type"`
	}
	Redis struct {
		Name     string `yaml:"name"`
		Uri      string `yaml:"uri"`
		Key      string `yaml:"key"`
		Type     string `yaml:type"`
		Password string `yaml:type"`
		Database int64  `yaml:database"`
	}
}

func Setup(l logger.Log, fname string) (cfg Config, err error) {
	var fd *os.File
	var fs os.FileInfo
	var data []byte

	Log = l

	if fs, err = os.Stat(fname); err != nil {
		return
	}

	data = make([]byte, fs.Size())

	if fd, err = os.Open(fname); err != nil {
		return
	}

	if _, err = fd.Read(data); err != nil {
		return
	}

	cfg = Config{}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return
	}

	cfg.Hostname, _ = os.Hostname()

	Log.Verbose("Configuration initialized")

	return
}
