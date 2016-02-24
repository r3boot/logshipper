package inputs

import (
	"encoding/json"
	"errors"
	"github.com/r3boot/logshipper/3rdparty/tail"
	"github.com/r3boot/logshipper/config"
	"os"
	"path/filepath"
	"time"
)

type MonitoredFile struct {
	Tail      *tail.Tail
	Name      string
	Path      string
	SinceDB   string
	Timestamp time.Time
	TsFormat  string
	Process   func(string, string) ([]byte, time.Time, error)
	Control   chan int
	Done      chan bool
}

func NewMonitoredFile(name string, fname string, ftype string, tsformat string) (tf *MonitoredFile, err error) {
	var fs os.FileInfo

	tf = &MonitoredFile{
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
	case config.T_JSON:
		{
			tf.Process = JSONParseLine
		}
	default:
		{
			err = errors.New("Unknown input type specified")
			return
		}
	}

	if _, err = os.Stat(filepath.Dir(tf.SinceDB)); err != nil {
		return
	}

	if fs, err = os.Stat(tf.SinceDB); err == nil {
		var fd *os.File
		var data []byte
		var ts time.Time

		data = make([]byte, fs.Size())

		if fd, err = os.Open(tf.SinceDB); err != nil {
			tf = nil
			return
		}

		if _, err = fd.Read(data); err != nil {
			tf = nil
			return
		}

		fd.Close()

		if err = json.Unmarshal(data, &ts); err != nil {
			return
		}

		tf.Timestamp = ts

		Log.Verbose("[" + name + "]: Reading logs starting from " + ts.String())
	} else {
		err = nil
		Log.Verbose("[" + name + "]: No sincedb found, reading logs from start")
	}

	return
}

func (tf *MonitoredFile) SaveSinceDB() (err error) {
	var fd *os.File
	var data []byte
	var written int

	fd, err = os.OpenFile(tf.SinceDB, (os.O_CREATE | os.O_WRONLY), 0600)
	if err != nil {
		return
	}
	defer fd.Close()

	if data, err = json.Marshal(tf.Timestamp); err != nil {
		Log.Warning("marshall(): " + err.Error())
		return
	}
	written, err = fd.Write(data)
	if err != nil {
		return
	}
	if written != len(data) {
		err = errors.New("Invalid number of bytes written")
		return
	}

	return
}

func (tf *MonitoredFile) Parse(output chan []byte) (err error) {
	var line *tail.Line
	var event []byte
	var ts time.Time
	var ok bool
	var stop_loop bool
	var t *tail.Tail

	t, err = tail.TailFile(tf.Path, TailConfig)
	if err != nil {
		return
	}
	tf.Tail = t

	stop_loop = false
	for {
		if stop_loop {
			break
		}

		select {
		case line, ok = <-tf.Tail.Lines:
			{
				// Check tail status
				if !ok {
					err = tf.Tail.Err()
					if err != nil {
						Log.Warning("Tail failed for " + tf.Tail.Filename)
						stop_loop = true
						break
					} else {
						Log.Warning("Unknown error reading " + tf.Tail.Filename)
						stop_loop = true
						break
					}
				}

				// Process tail output into event
				event, ts, err = tf.Process(line.Text, tf.TsFormat)
				if err != nil {
					Log.Debug(line.Text)
					Log.Debug(tf.TsFormat)
					Log.Warning("[" + tf.Name + "]: " + err.Error())
					continue
				}

				// Check timestamps to prevent duplicate logs
				if ts.After(tf.Timestamp) || ts.Equal(tf.Timestamp) {
					output <- event
					tf.Timestamp = ts
					if err = tf.SaveSinceDB(); err != nil {
						return
					}
				} else {
					Log.Debug("Timestamp " + ts.String() + " is in the past")
				}
			}
		case cmd := <-tf.Control:
			{
				switch cmd {
				case config.CMD_CLEANUP:
					{
						tf.Tail.Stop()
						tf.SaveSinceDB()
						Log.Debug("Wrote sincedb for " + tf.Name)
						stop_loop = true
						break
					}
				default:
					{
						Log.Debug("invalid command received")
					}
				}
			}
		}
	}

	tf.Done <- true
	return
}
