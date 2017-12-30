package logger

const (
	LOG_INFO    = "I"
	LOG_DEBUG   = "D"
	LOG_WARNING = "W"
	LOG_FATAL   = "F"
)

type Logger struct {
	UseTimestamp bool
	UseVerbose   bool
	UseDebug     bool
}
