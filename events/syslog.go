package events

import (
	"encoding/json"
	"time"
)

type SyslogData struct {
	Hostname string `json:"hostname"`
	Program  string `json:"program"`
	Pid      int64  `json:"pid,omitifempty"`
	Message  string `json:"message"`
}

type SyslogEvent struct {
	Timestamp time.Time  `json:"timestamp"`
	EventType string     `json:"event_type"`
	Path      string     `json:"path"`
	Host      string     `json:"host"`
	Syslog    SyslogData `json:"syslog"`
}

func (se *SyslogEvent) Serialize() (result []byte, err error) {
	result, err = json.Marshal(se)
	return
}

func NewSyslogEvent() (se SyslogEvent) {
	se = SyslogEvent{}
	se.EventType = "syslog"
	se.Host = Config.Hostname
	return
}
