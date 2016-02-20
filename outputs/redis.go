package outputs

import (
	"github.com/r3boot/logshipper/config"
	"gopkg.in/redis.v3"
)

var RedisClient *redis.Client
var LogKey string

type RedisShipper struct {
	Name    string
	Uri     string
	Key     string
	Type    string
	Client  *redis.Client
	Control chan int
	Done    chan bool
}

func NewRedisShipper() (rs *RedisShipper, err error) {
	rs = &RedisShipper{
		Name: Config.Redis.Name,
		Uri:  Config.Redis.Uri,
		Key:  Config.Redis.Key,
		Type: Config.Redis.Type,
		Client: redis.NewClient(&redis.Options{
			Addr:     Config.Redis.Uri,
			Password: Config.Redis.Password,
			DB:       Config.Redis.Database,
		}),
		Control: make(chan int, 1),
		Done:    make(chan bool, 1),
	}

	// Test connectivity
	_, err = rs.Client.Ping().Result()
	if err != nil {
		rs = nil
	}

	return
}

func (rs *RedisShipper) Ship(logdata chan []byte) (err error) {
	var stop_loop bool

	stop_loop = false
	for {
		if stop_loop {
			break
		}

		select {
		case event := <-logdata:
			{
				event_s := string(event)
				Log.Debug("Sending event to redis: " + event_s)
				rs.Client.RPush(LogKey, event_s)
			}
		case cmd := <-rs.Control:
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

	rs.Done <- true
	return
}
