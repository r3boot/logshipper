package config

import "github.com/r3boot/logshipper/lib/logger"

/*
 * Various constants used to configure the input system
 */
const (
	T_SYSLOG     = "syslog"     // Syslog logging with RFC3339 timestamps
	T_CLF        = "clf"        // HTTP Common Log Format
	T_SURICATA   = "suricata"   // JSON log format
	T_EXIM       = "exim"       // Exim4 log format
	T_AMQPBRIDGE = "amqpbridge" // Bridge from Amqp to Elasticsearch

	/*
	 * Various constants used to configure the output system
	 */
	T_STDOUT = 0 // Write to stdout
	T_REDIS  = 1 // Write to redis

	/*
	 * Various constants used to control the input/output threads
	 */
	CMD_CLEANUP = 0 // Stop whatever you're doing and cleanup

	/*
	 * Constants used to determine the suricata event type
	 */
	S_ALERT    string = "alert"
	S_FILEINFO string = "fileinfo"
	S_HTTP     string = "http"
	S_TLS      string = "tls"

	/*
	 * Constants defining different time formats for different logfiles
	 */
	TF_CLF string = "02/Jan/2006:15:04:05 -0700"

	TF_SURICATA string = "2006-01-02T15:04:05.000000-0700"
	TF_EXIM     string = "2006-01-02 15:04:05"

	MAX_CFG_SIZE int = 65534
)

type Config struct {
	Hostname string
	Inputs   []struct {
		Name     string `yaml:"name"`
		Path     string `yaml:"path"`
		Type     string `yaml:"type"`
		TsFormat string `yaml:"ts_format"`
	} `yaml:"inputs"`
	Outputs []struct {
		Name string `yaml:"name"`
		Type string `yaml:"type"`
	} `yaml:"outputs"`
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
	ES struct {
		Name  string `yaml:"name"`
		Host  string `yaml:"host"`
		Port  int    `yaml:"port"`
		Index string `yaml:"index"`
	} `yaml:"elasticsearch"`
}

var TIME_FORMATS = []string{
	"ANSIC",
	"UnixDate",
	"RubyDate",
	"RFC822",
	"RFC822Z",
	"RFC850",
	"RFC1123",
	"RFC1123Z",
	"RFC3339",
	"RFC3339Nano",
	"Kitchen",
	"Stamp",
	"StampMilli",
	"StampMicro",
	"StampNano",
}

var log *logger.Logger
