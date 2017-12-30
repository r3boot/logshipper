package outputs

import (
	"strconv"

	"fmt"

	"github.com/r3boot/logshipper/lib/config"
	"github.com/streadway/amqp"
)

type AmqpShipper struct {
	Name       string
	Type       string
	Url        string
	Queue      amqp.Queue
	Channel    *amqp.Channel
	Connection *amqp.Connection
	Control    chan int
	Done       chan bool
}

func NewAmqpShipper() (*AmqpShipper, error) {
	var err error

	user := cfg.Amqp.Username
	pass := cfg.Amqp.Password
	host := cfg.Amqp.Host
	port := strconv.Itoa(cfg.Amqp.Port)

	url := "amqp://" + user + ":" + pass + "@" + host + ":" + port
	as := &AmqpShipper{
		Name:    cfg.Amqp.Name,
		Type:    "amqp-input",
		Url:     url,
		Control: make(chan int, 1),
		Done:    make(chan bool, 1),
	}

	as.Connection, err = amqp.Dial(as.Url)
	if err != nil {
		return nil, fmt.Errorf("NewAmqpShipper amqp.Dial: %v", err)
	}

	as.Channel, err = as.Connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("NewAmqpShipper Connection.Channel: %v", err)
	}

	err = as.Channel.ExchangeDeclare(
		cfg.Amqp.Exchange, // Name of exchange
		"fanout",          // Type of exchange
		true,              // Durable
		false,             // Auto-deleted
		false,             // Internal queue
		false,             // no-wait
		nil,               // Arguments
	)
	if err != nil {
		return nil, fmt.Errorf("NewAmqpShipper: Channel.ExchangeDeclare: %v", err)
	}

	return as, nil
}

func (as *AmqpShipper) Ship(logdata chan []byte) error {
	defer as.Channel.Close()
	defer as.Connection.Close()

	stop_loop := false
	for {
		if stop_loop {
			break
		}

		select {
		case event := <-logdata:
			{
				event_s := string(event)
				log.Debugf("AmqpShipper.Ship: Sending event to AMQP: %s", event_s)
				err := as.Channel.Publish(
					cfg.Amqp.Exchange, // exchange to use
					"",                // key to use for routing
					false,             // mandatory
					false,             // immediate
					amqp.Publishing{
						ContentType: "application/json",
						Body:        event,
					},
				)
				if err != nil {
					log.Warningf("AmqpShipper.Ship Channel.Publish: %v", err)
					continue
				}
			}
		case cmd := <-as.Control:
			{
				switch cmd {
				case config.CMD_CLEANUP:
					{
						log.Debugf("AmqpShipper.Ship: Shutting down")
						stop_loop = true
						continue
					}
				default:
					{
						log.Warningf("AmqpShipper.Ship: Invalid command received")
					}
				}
			}
		}
	}

	as.Done <- true

	return nil
}
