package inputs

import (
	"encoding/json"
	"fmt"
	"time"
)

func SuricataParseLine(line, fname, tsformat string) ([]byte, time.Time, error) {
	var event map[string]interface{}

	err := json.Unmarshal([]byte(line), &event)
	if err != nil {
		return nil, time.Now(), fmt.Errorf("SuricataParseLine json.Unmarshal: %v", err)
	}
	event["host"] = cfg.Hostname
	event["path"] = fname

	ts, err := time.Parse(tsformat, event["timestamp"].(string))
	if err != nil {
		return nil, time.Now(), fmt.Errorf("SuricataParseLine time.Parse: %v", err)
	}

	e, err := json.Marshal(event)
	if err != nil {
		return nil, time.Now(), fmt.Errorf("SuricataParseLine json.Marshal: %v", err)
	}

	return e, ts, nil
}
