package inputs

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"fmt"

	"tail"

	"github.com/r3boot/logshipper/lib/config"
)

type MonitoredFile struct {
	Tail      *tail.Tail
	Name      string
	Path      string
	SinceDB   string
	Timestamp time.Time
	TsFormat  string
	Process   func(string, string, string) ([]byte, time.Time, error)
	Control   chan int
	Done      chan bool
}

func NewMonitoredFile(name, fname, ftype, tsformat string) (*MonitoredFile, error) {
	tf := &MonitoredFile{
		Name:     name,
		Path:     fname,
		TsFormat: tsformat,
		SinceDB:  "/var/lib/logshipper/" + name + ".sincedb",
		Control:  make(chan int, 1),
		Done:     make(chan bool, 1),
	}

	switch ftype {
	case config.T_SYSLOG:
		{
			tf.Process = SyslogParseLine
		}
	case config.T_CLF:
		{
			tf.Process = CLFParseLine
		}
	case config.T_SURICATA:
		{
			tf.Process = SuricataParseLine
		}
	case config.T_EXIM:
		{
			tf.Process = EximParseLine
		}
	default:
		{
			return nil, fmt.Errorf("NewMonitoredFile: unknown input type specified")
		}
	}

	_, err := os.Stat(filepath.Dir(tf.SinceDB))
	if err != nil {
		return nil, fmt.Errorf("NewMonitoredFile os.Stat: %v", err)
	}

	fs, err := os.Stat(tf.SinceDB)
	if err == nil {
		data := make([]byte, fs.Size())

		data, err := ioutil.ReadFile(tf.SinceDB)
		if err != nil {
			return nil, fmt.Errorf("NewMonitoredFile ioutil.ReadFile: %v", err)
		}

		ts := time.Time{}
		if err = json.Unmarshal(data, &ts); err != nil {
			return nil, fmt.Errorf("NewMonitoredFile json.Unmarshal: %v", err)
		}

		tf.Timestamp = ts

		log.Debugf("NewMonitoredFile[%s]: Reading logs starting from %s", name, ts.String())
	} else {
		log.Debugf("NewMonitoredFile[%s]: No sincedb found, reading logs from start", name)
	}

	return tf, nil
}

func (tf *MonitoredFile) SaveSinceDB() error {
	data, err := json.Marshal(tf.Timestamp)
	if err != nil {
		return fmt.Errorf("MonitoredFile.SaveSinceDB json.Marshal: %v", err)
	}

	err = ioutil.WriteFile(tf.SinceDB, data, 0600)
	if err != nil {
		return fmt.Errorf("MonitoredFile.SaveSinceDB ioutil.WriteFile: %v", err)
	}

	return nil
}

func (tf *MonitoredFile) Parse(output chan []byte) error {
	t, err := tail.TailFile(tf.Path, TailConfig)
	if err != nil {
		return fmt.Errorf("MonitoredFile.Parse tail.TailFile: %v", err)
	}
	tf.Tail = t

	stop_loop := false
	for {
		if stop_loop {
			break
		}

		select {
		case line, ok := <-tf.Tail.Lines:
			{
				// Check tail status
				if !ok {
					err = tf.Tail.Err()
					if err != nil {
						return fmt.Errorf("MonitoredFile.Parse[%s] tf.Tail: %v", tf.Name, err)
					} else {
						return fmt.Errorf("MonitoredFile.Parse[%s]: unknown error", tf.Name)
					}
				}

				// Process tail output into event
				event, ts, err := tf.Process(line.Text, tf.Path, tf.TsFormat)
				if err != nil {
					log.Warningf("MonitoredFile.Parse[%s]: %v", tf.Name, err)
					continue
				}

				// Check timestamps to prevent duplicate logs
				if ts.After(tf.Timestamp) || ts.Equal(tf.Timestamp) {
					output <- event
					tf.Timestamp = ts
					err = tf.SaveSinceDB()
					if err != nil {
						return fmt.Errorf("MonitoredFile.Parse[%s]: %v", tf.Name, err)
					}
				} else {
					log.Debugf("MonitoredFile.Parse[%s]: Timestamp %v is in the past", tf.Name, ts)
				}
			}
		case cmd := <-tf.Control:
			{
				switch cmd {
				case config.CMD_CLEANUP:
					{
						tf.Tail.Stop()
						tf.SaveSinceDB()
						log.Debugf("MonitoredFile.Parse[%s]: wrote sincedb", tf.Name)
						stop_loop = true
						break
					}
				default:
					{
						log.Debugf("MonitoredFile.Parse[%s]: invalid command received", tf.Name)
					}
				}
			}
		}
	}

	tf.Done <- true
	return nil
}
