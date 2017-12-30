package inputs

import (
	"regexp"
	"strconv"
	"time"

	"fmt"

	"github.com/r3boot/logshipper/lib/events"
)

// 2016-02-08T03:50:51.927029+01:00 shell sshd[27725]: Failed password for root from 125.88.177.93 port 34261 ssh2
const (
	re_syslog_1 = "^(?P<ts>[0-9T:+-.]+)\\ (?P<hostname>[a-zA-Z0-9_-]+)\\ (?P<program>[a-zA-Z0-9-_]+)\\[(?P<pid>[0-9]+)\\]:\\ +(?P<message>.*)"

	// 2016-02-08T11:39:27.166653+01:00 shell sudo:   r3boot : TTY=pts/44 ; PWD=/people/r3boot ; USER=root ; COMMAND=/bin/bash
	re_syslog_2 = "^(?P<ts>[0-9T:+-.]+)\\ (?P<hostname>[a-zA-Z0-9_-]+)\\ (?P<program>[a-zA-Z0-9-_]+):\\ +(?P<message>[a-zA-Z0-9].*)"

	// 016-02-25T01:05:05.521035+01:00 nic1 kernel: [2687937.343882] DROP INVALID: IN=eth0 OUT= MAC=52:54:00:d3:e7:23:00:14:f6:0b:ff:48:08:00 SRC=10.42.15.32 DST=10.42.0.63 LEN=40 TOS=0x00 PREC=0x00 TTL=63 ID=58661 DF PROTO=TCP SPT=58754 DPT=636 WINDOW=0 RES=0x00 RST URGP=0
	re_syslog_3 = "^(?P<ts>[0-9T:+-.]+)\\ (?P<hostname>[a-zA-Z0-9_-]+)\\ (?P<program>[a-zA-Z0-9-_]+):\\ \\[[0-9.]+\\]\\ (?P<message>.*)"
)

var syslog_regexps []*regexp.Regexp

func SetupSyslog() (err error) {
	var re *regexp.Regexp

	if re, err = regexp.Compile(re_syslog_1); err != nil {
		return
	}
	syslog_regexps = append(syslog_regexps, re)

	if re, err = regexp.Compile(re_syslog_2); err != nil {
		return
	}
	syslog_regexps = append(syslog_regexps, re)

	if re, err = regexp.Compile(re_syslog_3); err != nil {
		return
	}
	syslog_regexps = append(syslog_regexps, re)

	return
}

func SyslogParseLine(line, fname, tsformat string) ([]byte, time.Time, error) {
	match := []string{}
	keys := []string{}
	for _, re := range syslog_regexps {
		all_matches := re.FindAllStringSubmatch(line, -1)
		if len(all_matches) > 0 {
			match = all_matches[0]
			keys = re.SubexpNames()
			break
		}
	}

	if len(match) == 0 {
		return nil, time.Now(), fmt.Errorf("SyslogParseLine: No match found")
	}

	r := map[string]string{}
	for i, v := range match {
		r[keys[i]] = v
	}

	ts, err := time.Parse(tsformat, r["ts"])
	if err != nil {
		return nil, time.Now(), fmt.Errorf("SyslogParseLine time.Parse: %v", err)
	}

	se := events.NewSyslogEvent()
	se.Timestamp = ts
	se.Path = fname
	se.Syslog.Hostname = r["hostname"]
	se.Syslog.Program = r["program"]
	se.Syslog.Message = r["message"]

	if pid_s, ok := r["pid"]; ok {
		pid, _ := strconv.ParseInt(pid_s, 10, 0)
		se.Syslog.Pid = pid
	}

	e, err := se.Serialize()

	return e, ts, nil
}
