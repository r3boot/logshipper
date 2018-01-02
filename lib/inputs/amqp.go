package inputs

import (
	"strconv"

	"fmt"

	"github.com/r3boot/logshipper/lib/config"
	"github.com/streadway/amqp"
)

type AmqpSlurper struct {
	Name       string
	Type       string
	Url        string
	Queue      amqp.Queue
	Channel    *amqp.Channel
	Connection *amqp.Connection
	Deliveries <-chan amqp.Delivery
	Control    chan int
	Done       chan bool
}

func NewAmqpSlurper() (*AmqpSlurper, error) {
	var err error

	user := cfg.Amqp.Username
	pass := cfg.Amqp.Password
	host := cfg.Amqp.Host
	port := strconv.Itoa(cfg.Amqp.Port)

	url := "amqp://" + user + ":" + pass + "@" + host + ":" + port
	as := &AmqpSlurper{
		Name:    cfg.Amqp.Name,
		Type:    "amqp-input",
		Url:     url,
		Control: make(chan int, 1),
		Done:    make(chan bool, 1),
	}

	as.Connection, err = amqp.Dial(as.Url)
	if err != nil {
		return nil, fmt.Errorf("NewAmqpSlurper amqp.Dial: %v", err)
	}

	as.Channel, err = as.Connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("NewAmqpSlurper Connection.Channel: %v", err)
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
		return nil, fmt.Errorf("NewAmqpSlurper Channel.ExchangeDeclare: %v", err)
	}

	as.Queue, err = as.Channel.QueueDeclare(
		cfg.Amqp.Exchange, // name of the queue
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // noWait
		nil,               // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Declare: %s", err)
	}

	if err = as.Channel.QueueBind(
		as.Queue.Name,     // name of the queue
		"",                // bindingKey
		cfg.Amqp.Exchange, // sourceExchange
		false,             // noWait
		nil,               // arguments
	); err != nil {
		return nil, fmt.Errorf("Queue Bind: %s", err)
	}

	as.Deliveries, err = as.Channel.Consume(
		as.Queue.Name, // name
		"",            // consumerTag,
		false,         // noAck
		false,         // exclusive
		false,         // noLocal
		false,         // noWait
		nil,           // arguments
	)
	if err != nil {
		log.Debugf("error: %v", err)
		return nil, fmt.Errorf("NewAmqpSlurper Channel.Consume: %v", err)
	}

	return as, nil
}

func (as *AmqpSlurper) Parse(output chan []byte) error {
	defer as.Channel.Close()
	defer as.Connection.Close()

	log.Debugf("AmqpSlurper.Parse: Starting Parse loop")
	stop_loop := false
	for {
		if stop_loop {
			break
		}

		select {
		case delivery := <-as.Deliveries:
			{
				output <- delivery.Body
			}
		case cmd := <-as.Control:
			{
				switch cmd {
				case config.CMD_CLEANUP:
					{
						log.Debugf("AmqpSlurper.Parse: Shutting down")
						stop_loop = true
						continue
					}
				default:
					{
						log.Warningf("AmqpSlurper.Parse: Invalid command received")
					}
				}
			}
		}
	}

	as.Done <- true

	return nil
}
