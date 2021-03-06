package events

import (
	"encoding/json"
	"time"

	"github.com/r3boot/logshipper/lib/config"
)

type HttpData struct {
	Ident     string `json:"ident"`
	User      string `json:"user"`
	Method    string `json:"method"`
	Resource  string `json:"resource"`
	Protocol  string `json:"protocol"`
	Response  string `json:"response"`
	Size      string `json:"size"`
	Referrer  string `json:"referrer"`
	UserAgent string `json:"useragent"`
}

type HttpEvent struct {
	Timestamp time.Time `json:"timestamp"`
	EventType string    `json:"event_type"`
	Type      string    `json:"type"`
	Path      string    `json:"path"`
	Host      string    `json:"host"`
	SrcIp     string    `json:"src_ip"`
	Http      HttpData  `json:"http"`
}

func (se *HttpEvent) Serialize() (result []byte, err error) {
	result, err = json.Marshal(se)
	return
}

func NewHttpEvent() (se HttpEvent) {
	se = HttpEvent{}
	se.EventType = config.T_CLF
	se.Host = cfg.Hostname
	se.Type = cfg.ELK.Type
	return
}
