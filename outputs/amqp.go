package outputs

import (
	"github.com/r3boot/logshipper/config"
	"github.com/streadway/amqp"
	"strconv"
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

func NewAmqpShipper() (as *AmqpShipper, err error) {
	var url string

	user := Config.Amqp.Username
	pass := Config.Amqp.Password
	host := Config.Amqp.Host
	port := strconv.Itoa(Config.Amqp.Port)

	url = "amqp://" + user + ":" + pass + "@" + host + ":" + port
	as = &AmqpShipper{
		Name:    Config.Amqp.Name,
		Type:    "amqp-input",
		Url:     url,
		Control: make(chan int, 1),
		Done:    make(chan bool, 1),
	}

	if as.Connection, err = amqp.Dial(as.Url); err != nil {
		//as = nil
		return
	}

	if as.Channel, err = as.Connection.Channel(); err != nil {
		//as = nil
		return
	}

	err = as.Channel.ExchangeDeclare(
		Config.Amqp.Exchange, // Name of exchange
		"fanout",             // Type of exchange
		true,                 // Durable
		false,                // Auto-deleted
		false,                // Internal queue
		false,                // no-wait
		nil,                  // Arguments
	)
	if err != nil {
		return
	}

	return
}

func (as *AmqpShipper) Ship(logdata chan []byte) (err error) {
	var stop_loop bool

	defer as.Channel.Close()
	defer as.Connection.Close()

	stop_loop = false
	for {
		if stop_loop {
			break
		}

		select {
		case event := <-logdata:
			{
				event_s := string(event)
				Log.Debug("Sending event to AMQP: " + event_s)
				err = as.Channel.Publish(
					Config.Amqp.Exchange, // exchange to use
					"",                   // key to use for routing
					false,                // mandatory
					false,                // immediate
					amqp.Publishing{
						ContentType: "application/json",
						Body:        event,
					},
				)
				if err != nil {
					Log.Warning(err)
					continue
				}
			}
		case cmd := <-as.Control:
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

	as.Done <- true
	return
}
