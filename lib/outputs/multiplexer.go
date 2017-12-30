package outputs

import (
	"github.com/r3boot/logshipper/lib/config"
)

type OutputMultiplexer struct {
	Control chan int
	Done    chan bool
}

func NewOutputMultiplexer() *OutputMultiplexer {
	om := &OutputMultiplexer{
		Control: make(chan int, 1),
		Done:    make(chan bool, 1),
	}

	return om
}

func (om *OutputMultiplexer) Run(logdata chan []byte) error {
	var redis_chan chan []byte
	var amqp_chan chan []byte

	redis_enabled := false
	if cfg.Redis.Name != "" {
		redis_chan = make(chan []byte, 1)
		go Redis.Ship(redis_chan)
		redis_enabled = true
		log.Debugf("OutputMultiplexer.Run: Started Redis output shipper")
	}

	amqp_enabled := false
	if cfg.Amqp.Name != "" {
		amqp_chan = make(chan []byte, 1)
		go Amqp.Ship(amqp_chan)
		amqp_enabled = true
		log.Debugf("OutputMultiplexer.Run: Started AMQP output shipper")
	}

	stop_loop := false
	for {
		if stop_loop {
			break
		}

		select {
		case event := <-logdata:
			{
				if redis_enabled {
					log.Debugf("OutputMultiplexer.Run: Multiplexing event to Redis")
					redis_chan <- event
				}
				if amqp_enabled {
					log.Debugf("OutputMultiplexer.Run: Multiplexing event to AMQP")
					amqp_chan <- event
				}
			}
		case cmd := <-om.Control:
			{
				switch cmd {
				case config.CMD_CLEANUP:
					{
						log.Debugf("OutputMultiplexer.Run: Shutting down")
						stop_loop = true
						continue
					}
				default:
					{
						log.Warningf("OutputMultiplexer.Run: Invalid command received")
					}
				}
			}
		}
	}

	om.Done <- true

	return nil
}
