package config

/*
 * Various constants used to configure the input system
 */
const T_SYSLOG string = "syslog"     // Syslog logging with RFC3339 timestamps
const T_CLF string = "clf"           // HTTP Common Log Format
const T_SURICATA string = "suricata" // JSON log format
const T_EXIM string = "exim"         // Exim4 log format

/*
 * Various constants used to configure the output system
 */
const T_STDOUT int = 0 // Write to stdout
const T_REDIS int = 1  // Write to redis

/*
 * Various constants used to control the input/output threads
 */
const CMD_CLEANUP int = 0 // Stop whatever you're doing and cleanup

var TIME_FORMATS []string = []string{
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

/*
 * Constants used to determine the suricata event type
 */
const S_ALERT string = "alert"
const S_FILEINFO string = "fileinfo"
const S_HTTP string = "http"
const S_TLS string = "tls"

/*
 * Constants defining different time formats for different logfiles
 */
const TF_CLF string = "02/Jan/2006:15:04:05 -0700"
const TF_SURICATA string = "2006-01-02T15:04:05.000000-0700"
const TF_EXIM string = "2006-01-02 15:04:05"
