package inputs

import (
	"github.com/r3boot/logshipper/events"
	"regexp"
	"time"
)

// 51.255.65.16 - - [08/Feb/2016:22:53:42 +0100] "GET /index.php/archives/2007/10/20/letters-from-henry/trackback/ HTTP/1.1" 302 5 "-" "Mozilla/5.0 (compatible; AhrefsBot/5.0; +http://ahrefs.com/robot/)" "-"
const re_al_1 string = "^(?P<src_ip>[0-9a-fA-F:.]+)\\ (?P<ident>[a-zA-Z0-9-_]+)\\ (?P<user>[a-zA-Z0-9-_]+)\\ \\[(?P<ts>[a-zA-Z0-9/:\\ \\+\\-]+)\\]\\ \"(?P<method>[A-Z]+)\\ (?P<resource>.*)\\ HTTP/(?P<proto>[0-2.]{3})\"\\ (?P<resp>[0-9]{3})\\ (?P<size>[0-9]+)\\ \"(?P<ref>.*)\" \"(?P<ua>.*)\" \".*\"$"

var al_regexps []*regexp.Regexp

func SetupCFL() (err error) {
	var re *regexp.Regexp

	if re, err = regexp.Compile(re_al_1); err != nil {
		return
	}
	al_regexps = append(al_regexps, re)

	return
}

func CLFParseLine(line string, fname string, tsformat string) (e []byte, ts time.Time, err error) {
	var all_matches [][]string
	var match []string
	var keys []string

	for _, re := range al_regexps {
		all_matches = re.FindAllStringSubmatch(line, -1)
		if len(all_matches) > 0 {
			match = all_matches[0]
			keys = re.SubexpNames()
			break
		}
	}

	r := map[string]string{}
	for i, v := range match {
		r[keys[i]] = v
	}

	if ts, err = time.Parse(tsformat, r["ts"]); err != nil {
		return
	}

	he := events.NewHttpEvent()
	he.Timestamp = ts
	he.Path = fname
	he.SrcIp = r["src_ip"]
	he.Http.Ident = r["ident"]
	he.Http.User = r["user"]
	he.Http.Method = r["method"]
	he.Http.Resource = r["resource"]
	he.Http.Protocol = r["proto"]
	he.Http.Response = r["resp"]
	he.Http.Size = r["size"]
	he.Http.Referrer = r["ref"]
	he.Http.UserAgent = r["ua"]

	e, err = he.Serialize()

	return
}
