package inputs

import (
	"errors"
	"github.com/r3boot/logshipper/events"
	"regexp"
	"strconv"
	"time"
)

// 2016-02-08T03:50:51.927029+01:00 shell sshd[27725]: Failed password for root from 125.88.177.93 port 34261 ssh2
const re_syslog_1 string = "^(?P<ts>[0-9T:+-.]+)\\ (?P<hostname>[a-zA-Z0-9_-]+)\\ (?P<program>[a-zA-Z0-9-_]+)\\[(?P<pid>[0-9]+)\\]:\\ +(?P<message>.*)"

// 2016-02-08T11:39:27.166653+01:00 shell sudo:   r3boot : TTY=pts/44 ; PWD=/people/r3boot ; USER=root ; COMMAND=/bin/bash
const re_syslog_2 string = "^(?P<ts>[0-9T:+-.]+)\\ (?P<hostname>[a-zA-Z0-9_-]+)\\ (?P<program>[a-zA-Z0-9-_]+):\\ +(?P<message>[a-zA-Z0-9].*)"

// 016-02-25T01:05:05.521035+01:00 nic1 kernel: [2687937.343882] DROP INVALID: IN=eth0 OUT= MAC=52:54:00:d3:e7:23:00:14:f6:0b:ff:48:08:00 SRC=10.42.15.32 DST=10.42.0.63 LEN=40 TOS=0x00 PREC=0x00 TTL=63 ID=58661 DF PROTO=TCP SPT=58754 DPT=636 WINDOW=0 RES=0x00 RST URGP=0
const re_syslog_3 string = "^(?P<ts>[0-9T:+-.]+)\\ (?P<hostname>[a-zA-Z0-9_-]+)\\ (?P<program>[a-zA-Z0-9-_]+):\\ \\[[0-9.]+\\]\\ (?P<message>.*)"

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

func SyslogParseLine(line string, fname string, tsformat string) (e []byte, ts time.Time, err error) {
	var all_matches [][]string
	var match []string
	var keys []string

	for _, re := range syslog_regexps {
		all_matches = re.FindAllStringSubmatch(line, -1)
		if len(all_matches) > 0 {
			match = all_matches[0]
			keys = re.SubexpNames()
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
	}

	if ts, err = time.Parse(tsformat, r["ts"]); err != nil {
		return
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

	e, err = se.Serialize()

	return
}
