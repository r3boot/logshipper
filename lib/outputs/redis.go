package outputs

import (
	"fmt"

	"github.com/r3boot/logshipper/lib/config"
	"gopkg.in/redis.v3"
)

type RedisShipper struct {
	Name    string
	Uri     string
	Key     string
	Type    string
	Client  *redis.Client
	Control chan int
	Done    chan bool
}

func NewRedisShipper() (*RedisShipper, error) {
	rs := &RedisShipper{
		Name: cfg.Redis.Name,
		Uri:  cfg.Redis.Uri,
		Key:  cfg.Redis.Key,
		Type: cfg.Redis.Type,
		Client: redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Uri,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.Database,
		}),
		Control: make(chan int, 1),
		Done:    make(chan bool, 1),
	}

	// Test connectivity
	_, err := rs.Client.Ping().Result()
	if err != nil {
		return nil, fmt.Errorf("NewRedisShipper Client.Ping: %v", err)
	}

	return rs, nil
}

func (rs *RedisShipper) Ship(logdata chan []byte) error {
	stop_loop := false
	for {
		if stop_loop {
			break
		}

		select {
		case event := <-logdata:
			{
				event_s := string(event)
				log.Debugf("RedisShipper.Ship: Sending event to redis: " + event_s)
				rs.Client.RPush(rs.Key, event_s)
			}
		case cmd := <-rs.Control:
			{
				switch cmd {
				case config.CMD_CLEANUP:
					{
						log.Debugf("RedisShipper.Ship: Shutting down")
						stop_loop = true
						continue
					}
				default:
					{
						log.Warningf("RedisShipper.Ship: Invalid command received")
					}
				}
			}
		}
	}

	rs.Done <- true

	return nil
}
