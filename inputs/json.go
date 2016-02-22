package inputs

import (
	"encoding/json"
	"time"
)

type JSONTime struct {
	RawTimestamp string `json:"timestamp"`
}

func JSONParseLine(line string) (e []byte, ts time.Time, err error) {
	var jt JSONTime
	var tf string

	tf = "2006-01-02T15:04:05.000000"

	if err = json.Unmarshal([]byte(line), &jt); err != nil {
		Log.Debug("Failed to unmarshall JSON")
		return
	}

	if ts, err = time.Parse(tf, jt.RawTimestamp); err != nil {
		Log.Debug("Failed to parse time: " + jt.RawTimestamp)
		return
	}

	e = []byte(line)

	return
}
