package outputs

import (
	"github.com/r3boot/logshipper/config"
)

type LogShipper struct {
	Name  string
	Type  int
	Redis struct {
		Uri  string
		Key  string
		Type string
	}
	Process func([]byte) (err error)
	Control chan int
	Done    chan bool
}

func NewLogShipper(name string, otype int) (s *LogShipper, err error) {
	s = &LogShipper{
		Name:    name,
		Type:    otype,
		Control: make(chan int, 1),
		Done:    make(chan bool, 1),
	}

	switch otype {
	case config.T_STDOUT:
		{
			s.Process = ShipStdout
		}
	case config.T_REDIS:
		{

			s.Process = ShipRedis
		}
	}
	return
}

func (ls *LogShipper) Ship(logdata chan []byte) (err error) {
	var stop_loop bool

	if ls.Type == config.T_REDIS {
		Log.Debug("Setting up redis client")
		uri := ls.Redis.Uri
		key := ls.Redis.Key
		if err = SetupRedisClient(uri, key); err != nil {
			return
		}
	}

	stop_loop = false
	for {
		if stop_loop {
			break
		}

		select {
		case event := <-logdata:
			{
				ls.Process(event)
			}
		case cmd := <-ls.Control:
			{
				switch cmd {
				case config.CMD_CLEANUP:
					{
						Log.Debug("Performing cleanup")
						stop_loop = true
						continue
					}
				default:
					{
						Log.Warning("Invalid command received")
					}
				}
			}
		}
	}

	ls.Done <- true
	return
}
