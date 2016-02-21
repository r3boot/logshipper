package outputs

import (
	"github.com/r3boot/logshipper/config"
)

type OutputMultiplexer struct {
	Control chan int
	Done    chan bool
}

func NewOutputMultiplexer() (om *OutputMultiplexer, err error) {
	om = &OutputMultiplexer{
		Control: make(chan int, 1),
		Done:    make(chan bool, 1),
	}

	return
}

func (om *OutputMultiplexer) Run(logdata chan []byte) (err error) {
	var stop_loop bool
	var redis_enabled bool
	var redis_chan chan []byte
	var amqp_enabled bool
	var amqp_chan chan []byte

	redis_enabled = false
	if Config.Redis.Name != "" {
		redis_chan = make(chan []byte, 1)
		go Redis.Ship(redis_chan)
		redis_enabled = true
		Log.Debug("Started Redis output shipper")
	}

	amqp_enabled = false
	if Config.Amqp.Name != "" {
		amqp_chan = make(chan []byte, 1)
		go Amqp.Ship(amqp_chan)
		amqp_enabled = true
		Log.Debug("Started AMQP output shipper")
	}

	stop_loop = false
	for {
		if stop_loop {
			break
		}

		select {
		case event := <-logdata:
			{
				if redis_enabled {
					Log.Debug("Multiplexing event to Redis")
					redis_chan <- event
				}
				if amqp_enabled {
					Log.Debug("Multiplexing event to AMQP")
					amqp_chan <- event
				}
			}
		case cmd := <-om.Control:
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

	om.Done <- true
	return
}
