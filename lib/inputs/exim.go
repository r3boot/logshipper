package inputs

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"time"
)

const (
	// 2016-02-25 20:08:55 1aZ1HS-0006iF-Rk <= m-2vlkybcnxljidjw6iay5827f5opbillgehggdelein1h7sxg15@bounce.linkedin.com H=mail01.pyzuka.nl [2a01:7c8:c037:6::4]:47889 I=[2001:7b8:3:47:213:154:229:26]:25 P=esmtps X=TLS1.2:DHE_RSA_AES_128_CBC_SHA1:128 CV=no S=39472 id=1087531239.521310.1456427301019.JavaMail.app@ltx1-app7953.prod.linkedin.com
	RE_EXIM_1 = 0
	re_exim_1 = "^(?P<date>[0-9-]+)\\ (?P<time>[0-9:]+)\\ (?P<lid>[0-9a-zA-Z-]+)\\ <=\\ (?P<from>[a-zA-Z0-9-_+.@]+)\\ H=(?P<srchost>[a-zA-Z0-9._-]+)\\ \\[(?P<srcip>[0-9a-f.:]+)\\]:(?P<srcport>[0-9]+)\\ I=\\[(?P<dstip>[0-9a-f.:]+)\\]:(?P<dstport>[0-9]+)\\ P=(?P<proto>[a-zA-Z0-9]+)\\ X=(?P<tlsproto>[A-Z0-9_.:]+)\\ CV=(?P<tlsverify>[a-z]+)\\ S=(?P<size>[0-9]+)\\ id=(?P<rid>.*)"

	// 2016-02-25 20:02:24 1aZ1B6-0006g0-ND => r3boot <r3boot@r3blog.nl> R=mailbox T=dovecot S=7727 QT=4s DT=0s
	RE_EXIM_2 = 1
	re_exim_2 = "^(?P<date>[0-9-]+)\\ (?P<time>[0-9:]+)\\ (?P<lid>[0-9a-zA-Z-]+)\\ =>\\ (?P<to>.*)\\ R=(?P<dest>[a-zA-Z0-9-_]+)\\ T=(?P<transport>[a-zA-Z0-9-_]+)\\ S=(?P<size>[0-9]+)\\ QT=(?P<qt>[0-9]+)s\\ DT=(?P<dt>[0-9]+)s"

	// 2016-02-25 20:02:24 1aZ1B6-0006g0-ND Completed
	RE_EXIM_3 = 2
	re_exim_3 = "^(?P<date>[0-9-]+)\\ (?P<time>[0-9:]+)\\ (?P<lid>[0-9a-zA-Z-]+)\\ (?P<status>[a-zA-Z]+)$"

	// 2016-02-25 20:02:21 1aZ1B6-0006g0-ND DKIM: d=gmail.com s=20120113 c=relaxed/relaxed a=rsa-sha256 [verification succeeded]
	RE_EXIM_4 = 3
	re_exim_4 = "^(?P<date>[0-9-]+)\\ (?P<time>[0-9:]+)\\ (?P<lid>[0-9a-zA-Z-]+)\\ DKIM:\\ d=(?P<domain>[a-zA-Z0-9-.]+)\\ s=(?P<size>[0-9]+)\\ c=(?P<canon>[a-z/]+)\\ a=(?P<algo>[a-z0-9-]+)\\ \\[(?P<status>.*)\\]"

	// 2016-02-25 20:00:08 1aYGBw-0002FJ-Jt == bridget_meggett@hzhuixin.top R=dnslookup T=remote_smtp defer (-53): retry time not reached for any host
	RE_EXIM_5 = 4
	re_exim_5 = "^(?P<date>[0-9-]+)\\ (?P<time>[0-9:]+)\\ (?P<lid>[0-9a-zA-Z-]+)\\ ==\\ (?P<to>.*)\\ R=(?P<router>[a-zA-Z0-9]+)\\ T=(?P<transport>[a-zA-Z0-9-_]+)\\ (?P<status>[a-zA-Z0-9]+)\\ \\((?P<errcode>[0-9-]+)\\):\\ (?P<message>.*)"

	// 2016-02-25 19:53:39 H=ar-ix.net (mail.ar-ix.net) [2a00:1bd0:197:2:1::42]:53981 I=[2001:7b8:3:47:213:154:229:26]:25 X=TLS1.2:ECDHE_RSA_AES_256_GCM_SHA384:256 CV=no temporarily rejected MAIL <MullinsSonja52353@webmailcourrier.com>: Could not complete sender verify
	RE_EXIM_6 = 5
	re_exim_6 = "^(?P<date>[0-9-]+)\\ (?P<time>[0-9:]+)\\ H=(?P<host>[.a-zA-Z0-9-_]+)\\ \\((?P<helo>[.a-zA-Z0-9-_]+)\\)\\ \\[(?P<srcip>[0-9a-f.:]+)\\]:(?P<srcport>[0-9]+)\\ I=\\[(?P<dstip>[0-9a-f.:]+)\\]:(?P<dstport>[0-9]+)\\ X=(?P<tlsproto>[A-Z0-9_.:]+)\\ CV=(?P<tlsverify>[a-z]+)\\ (?P<action>[a-zA-Z0-9\\ ]+)\\ \\<(?P<to>.*)\\>:\\ (?P<message>.*)"
)

type EximIncomingMailEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Host      string    `json:"host"`
	Path      string    `json:"path"`
	EventType string    `json:"event_type"`
	DstIp     string    `json:"dst_ip"`
	DstPort   int64     `json:"dst_port"`
	SrcIp     string    `json:"src_ip"`
	SrcPort   int64     `json:"src_port"`
	Mail      struct {
		MsgId       string `json:"msg_id"`
		From        string `json:"from"`
		Proto       string `json:"proto"`
		TLSProto    string `json:"tls_proto"`
		TLSVerify   bool   `json:"tls_verify"`
		Size        int64  `json:"size"`
		RemoteMsgId string `json:"rid"`
	} `json:"mail"`
}

type EximMailStoreEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Host      string    `json:"host"`
	Path      string    `json:"path"`
	EventType string    `json:"event_type"`
	Mail      struct {
		MsgId        string `json:"msg_id"`
		To           string `json:"to"`
		Dest         string `json:"dest"`
		Transport    string `json:"transport"`
		Size         int64  `json:"size"`
		QueueTime    int64  `json:"qt"`
		DeliveryTime int64  `json:"dt"`
	} `json:"mail"`
}

type EximDeliverMailEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Host      string    `json:"host"`
	Path      string    `json:"path"`
	EventType string    `json:"event_type"`
	Mail      struct {
		MsgId  string `json:"msg_id"`
		Status string `json:"status"`
	} `json:"mail"`
}

type EximDKIMEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Host      string    `json:"host"`
	Path      string    `json:"path"`
	EventType string    `json:"event_type"`
	Mail      struct {
		MsgId string `json:"msg_id"`
	} `json:"mail"`
	DKIM struct {
		Domain    string `json:"domain"`
		Size      int64  `json:"size"`
		CanonAlgo string `json:"canon_algo"`
		SignAlgo  string `json:"sign_algo"`
		Status    string `json:"status"`
	} `json:"dkim"`
}

type EximDeliveryErrorEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Host      string    `json:"host"`
	Path      string    `json:"path"`
	EventType string    `json:"event_type"`
	Mail      struct {
		MsgId     string `json:"msg_id"`
		Router    string `json:"router"`
		Transport string `json:"transport"`
		Status    string `json:"status"`
		ErrCode   int64  `json:"errcode"`
		Message   string `json:"message"`
	} `json:"mail"`
}

type EximTempRejectEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Host      string    `json:"host"`
	Path      string    `json:"path"`
	EventType string    `json:"event_type"`
	DstIp     string    `json:"dst_ip"`
	DstPort   int64     `json:"dst_port"`
	SrcIp     string    `json:"src_ip"`
	SrcPort   int64     `json:"src_port"`
	Mail      struct {
		TLSProto  string `json:"tls_proto"`
		TLSVerify bool   `json:"tls_verify"`
		Action    string `json:"action"`
		Message   string `json:"message"`
	}
}

var exim_regexps []*regexp.Regexp

func SetupExim() (err error) {
	var re *regexp.Regexp

	if re, err = regexp.Compile(re_exim_1); err != nil {
		return
	}
	exim_regexps = append(exim_regexps, re)

	if re, err = regexp.Compile(re_exim_2); err != nil {
		return
	}
	exim_regexps = append(exim_regexps, re)

	if re, err = regexp.Compile(re_exim_3); err != nil {
		return
	}
	exim_regexps = append(exim_regexps, re)

	if re, err = regexp.Compile(re_exim_4); err != nil {
		return
	}
	exim_regexps = append(exim_regexps, re)

	if re, err = regexp.Compile(re_exim_5); err != nil {
		return
	}
	exim_regexps = append(exim_regexps, re)

	if re, err = regexp.Compile(re_exim_6); err != nil {
		return
	}
	exim_regexps = append(exim_regexps, re)

	return
}

func EximParseLine(line string, fname string, tsformat string) (e []byte, ts time.Time, err error) {
	var all_matches [][]string
	var match []string
	var keys []string
	var re_idx int
	var src_port int64
	var dst_port int64
	var size int64

	for idx, re := range exim_regexps {
		all_matches = re.FindAllStringSubmatch(line, -1)
		if len(all_matches) > 0 {
			match = all_matches[0]
			keys = re.SubexpNames()
			re_idx = idx
			break
		}
	}

	if len(match) == 0 {
		err = errors.New("No match found!")
		return
	}

	r := map[string]string{}
	for i, v := range match {
		r[keys[i]] = v
		// log.Debug(keys[i] + ": " + v)
	}

	if ts, err = time.Parse(tsformat, r["date"]+" "+r["time"]); err != nil {
		return
	}

	switch re_idx {
	case RE_EXIM_1:
		{
			src_port, err = strconv.ParseInt(r["srcport"], 10, 0)
			if err != nil {
				return
			}

			dst_port, err = strconv.ParseInt(r["dstport"], 10, 0)
			if err != nil {
				return
			}

			size, err = strconv.ParseInt(r["size"], 10, 0)
			if err != nil {
				return
			}

			event := &EximIncomingMailEvent{
				Timestamp: ts,
				Host:      cfg.Hostname,
				Path:      fname,
				EventType: "exim",
				DstIp:     r["dstip"],
				DstPort:   src_port,
				SrcIp:     r["srcip"],
				SrcPort:   dst_port,
			}

			event.Mail.MsgId = r["lid"]
			event.Mail.From = r["from"]
			event.Mail.TLSProto = r["tlsproto"]
			if r["tlsverify"] == "yes" {
				event.Mail.TLSVerify = true
			} else {
				event.Mail.TLSVerify = false
			}
			event.Mail.Size = size
			event.Mail.RemoteMsgId = r["rid"]

			e, err = json.Marshal(event)
		}
	case RE_EXIM_2:
		{
			var qt int64
			var dt int64

			if size, err = strconv.ParseInt(r["size"], 10, 0); err != nil {
				return
			}

			if qt, err = strconv.ParseInt(r["qt"], 10, 0); err != nil {
				return
			}

			if dt, err = strconv.ParseInt(r["dt"], 10, 0); err != nil {
				return
			}

			event := EximMailStoreEvent{
				Timestamp: ts,
				Host:      cfg.Hostname,
				Path:      fname,
				EventType: "exim",
			}

			event.Mail.MsgId = r["lid"]
			event.Mail.To = r["to"]
			event.Mail.Dest = r["dest"]
			event.Mail.Transport = r["transport"]
			event.Mail.Size = size
			event.Mail.QueueTime = qt
			event.Mail.DeliveryTime = dt

			e, err = json.Marshal(event)
		}
	case RE_EXIM_3:
		{
			event := &EximDeliverMailEvent{
				Timestamp: ts,
				Host:      cfg.Hostname,
				Path:      fname,
				EventType: "exim",
			}

			event.Mail.MsgId = r["lid"]
			event.Mail.Status = r["status"]

			e, err = json.Marshal(event)
		}
	case RE_EXIM_4:
		{
			if size, err = strconv.ParseInt(r["size"], 10, 0); err != nil {
				return
			}

			event := &EximDKIMEvent{
				Timestamp: ts,
				Host:      cfg.Hostname,
				Path:      fname,
				EventType: "exim",
			}

			event.Mail.MsgId = r["lid"]
			event.DKIM.Domain = r["domain"]
			event.DKIM.CanonAlgo = r["canon"]
			event.DKIM.SignAlgo = r["algo"]
			event.DKIM.Size = size
			event.DKIM.Status = r["status"]

			e, err = json.Marshal(event)
		}
	case RE_EXIM_5:
		{
			var errcode int64

			if errcode, err = strconv.ParseInt(r["errcode"], 10, 0); err != nil {
				return
			}

			event := &EximDeliveryErrorEvent{
				Timestamp: ts,
				Host:      cfg.Hostname,
				Path:      fname,
				EventType: "exim",
			}

			event.Mail.MsgId = r["lid"]
			event.Mail.Router = r["router"]
			event.Mail.Transport = r["transport"]
			event.Mail.Status = r["status"]
			event.Mail.ErrCode = errcode
			event.Mail.Message = r["message"]

			e, err = json.Marshal(event)
		}
	case RE_EXIM_6:
		{
			dst_port, err = strconv.ParseInt(r["srcport"], 10, 0)
			if err != nil {
				return
			}

			src_port, err = strconv.ParseInt(r["dstport"], 10, 0)
			if err != nil {
				return
			}

			event := &EximTempRejectEvent{
				Timestamp: ts,
				Host:      cfg.Hostname,
				Path:      fname,
				EventType: "exim",
				DstIp:     r["dstip"],
				DstPort:   dst_port,
				SrcIp:     r["srcip"],
				SrcPort:   src_port,
			}

			event.Mail.TLSProto = r["tlsproto"]
			if r["tlsverify"] == "yes" {
				event.Mail.TLSVerify = true
			} else {
				event.Mail.TLSVerify = false
			}
			event.Mail.Action = r["action"]
			event.Mail.Message = r["message"]

			e, err = json.Marshal(event)
		}
	}

	return
}
