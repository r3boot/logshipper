package inputs

import (
	"encoding/json"
	"github.com/r3boot/logshipper/config"
	"time"
)

type SuricataEventDiscovery struct {
	EventType string `json:"event_type"`
}

type JSONTime struct {
	RawTimestamp string `json:"timestamp"`
}

type SuricataAlertEvent struct {
	Timestamp string `json:"timestamp"`
	Host      string `json:"host"`
	Path      string `json:"path"`
	EventType string `json:"event_type"`
	SrcIp     string `json:"src_ip"`
	SrcPort   int    `json:"src_port"`
	DstIp     string `json:"dst_ip"`
	DstPort   int    `json:"dst_port"`
	Proto     string `json:"proto"`
	Alert     struct {
		Action      string `json:"action"`
		Gid         int    `json"gid"`
		SignatureId int64  `json:"signature_id"`
		Rev         int    `json:"rev"`
		Signature   string `json:"signature"`
		Category    string `json:"category"`
		Severity    int    `json:"severity"`
	} `json:"alert"`
}

type SuricataFileInfoEvent struct {
	Timestamp string `json:"timestamp"`
	Host      string `json:"host"`
	Path      string `json:"path"`
	EventType string `json:"event_type"`
	SrcIp     string `json:"src_ip"`
	SrcPort   int    `json:"src_port"`
	DstIp     string `json:"dst_ip"`
	DstPort   int    `json:"dst_port"`
	Proto     string `json:"proto"`
	Http      struct {
		Url           string `json:"url"`
		Hostname      string `json:"hostname"`
		HttpUserAgent string `json:"http_user_agent"`
	} `json:"http"`
	FileInfo struct {
		Filename string `json:"filename"`
		Magic    string `json:"magic"`
		Md5      string `json"md5"`
		Stored   bool   `json"stored"`
		Size     int64  `json:"size"`
		Type     string `json:"type"`
	} `json:"fileinfo"`
}

type SuricataHttpEvent struct {
	Timestamp string `json:"timestamp"`
	Host      string `json:"host"`
	Path      string `json:"path"`
	EventType string `json:"event_type"`
	SrcIp     string `json:"src_ip"`
	SrcPort   int    `json:"src_port"`
	DstIp     string `json:"dst_ip"`
	DstPort   int    `json:"dst_port"`
	Proto     string `json:"proto"`
	Http      struct {
		Hostname      string `json:"hostname"`
		Url           string `json:"url"`
		HttpUserAgent string `json:"http_user_agent"`
		HttpMethod    string `json:"http_method"`
		Protocol      string `json:"protocol"`
		length        int64  `json:"length"`
	} `json:"http"`
}

type SuricataTlsEvent struct {
	Timestamp string `json:"timestamp"`
	Host      string `json:"host"`
	Path      string `json:"path"`
	EventType string `json:"event_type"`
	SrcIp     string `json:"src_ip"`
	SrcPort   int    `json:"src_port"`
	DstIp     string `json:"dst_ip"`
	DstPort   int    `json:"dst_port"`
	Proto     string `json:"proto"`
	Tls       struct {
		Subject     string `json:"subject"`
		IssuerDN    string `json:"issuerdn"`
		Fingerprint string `json:"fingerprint"`
		Version     string `json:"version"`
	}
}

func SuricataParseLine(line string, fname string, tsformat string) (e []byte, ts time.Time, err error) {
	var jt JSONTime
	var sed SuricataEventDiscovery

	if err = json.Unmarshal([]byte(line), &sed); err != nil {
		return
	}

	switch sed.EventType {
	case config.S_ALERT:
		{
			var sae SuricataAlertEvent
			if err = json.Unmarshal([]byte(line), &sae); err != nil {
				return
			}
			sae.Host = Config.Hostname
			sae.Path = fname
			if e, err = json.Marshal(sae); err != nil {
				return
			}
		}
	case config.S_FILEINFO:
		{
			var sfe SuricataFileInfoEvent
			if err = json.Unmarshal([]byte(line), &sfe); err != nil {
				return
			}
			sfe.Host = Config.Hostname
			sfe.Path = fname
			if e, err = json.Marshal(sfe); err != nil {
				return
			}
		}
	case config.S_HTTP:
		{
			var she SuricataHttpEvent
			if err = json.Unmarshal([]byte(line), &she); err != nil {
				return
			}
			she.Host = Config.Hostname
			she.Path = fname
			if e, err = json.Marshal(she); err != nil {
				return
			}
		}
	case config.S_TLS:
		{
			var ste SuricataTlsEvent
			if err = json.Unmarshal([]byte(line), &ste); err != nil {
				return
			}
			ste.Host = Config.Hostname
			ste.Path = fname
			if e, err = json.Marshal(ste); err != nil {
				return
			}
		}
	}

	if err = json.Unmarshal([]byte(line), &jt); err != nil {
		Log.Debug("Failed to unmarshall JSON")
		return
	}

	if ts, err = time.Parse(tsformat, jt.RawTimestamp); err != nil {
		Log.Debug("Failed to parse time: " + jt.RawTimestamp)
		return
	}

	return
}
