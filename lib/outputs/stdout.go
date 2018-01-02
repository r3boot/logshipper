package outputs

import "github.com/r3boot/logshipper/lib/config"

type StdoutShipper struct {
	Name    string
	Type    string
	Control chan int
	Done    chan bool
}

func NewStdoutShipper() *StdoutShipper {
	return &StdoutShipper{
		Control: make(chan int, 1),
		Done:    make(chan bool, 1),
	}
}

func (as *StdoutShipper) Ship(logdata chan []byte) error {
	stop_loop := false
	for {
		if stop_loop {
			break
		}

		select {
		case event := <-logdata:
			{
				log.Infof("Received: %s", event)
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
