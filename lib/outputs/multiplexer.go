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
	var redisChan chan []byte
	var amqpChan chan []byte
	var esChan chan []byte
	var stdoutChan chan []byte

	redisEnabled := false
	if Redis != nil && cfg.Redis.Name != "" {
		redisChan = make(chan []byte, 1)
		go Redis.Ship(redisChan)
		redisEnabled = true
		log.Debugf("OutputMultiplexer.Run: Started Redis log shipper")
	}

	amqpEnabled := false
	if Amqp != nil && cfg.Amqp.Name != "" {
		amqpChan = make(chan []byte, 1)
		go Amqp.Ship(amqpChan)
		amqpEnabled = true
		log.Debugf("OutputMultiplexer.Run: Started AMQP log shipper")
	}

	esEnabled := false
	if ES != nil && cfg.ES.Name != "" {
		esChan = make(chan []byte, 1)
		go ES.Ship(esChan)
		esEnabled = true
		log.Debugf("OutputMultiplexer.Run: Started Elasticsearch log shipper")
	}

	stdoutEnabled := false
	if Stdout != nil {
		stdoutChan = make(chan []byte, 1)
		go Stdout.Ship(stdoutChan)
		stdoutEnabled = true
		log.Debugf("OutputMultiplexer.Run: Started Stdout log shipper")
	}

	stop_loop := false
	for {
		if stop_loop {
			break
		}

		select {
		case event := <-logdata:
			{
				if redisEnabled {
					redisChan <- event
				}
				if amqpEnabled {
					amqpChan <- event
				}
				if esEnabled {
					esChan <- event
				}
				if stdoutEnabled {
					stdoutChan <- event
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
