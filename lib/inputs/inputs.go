package inputs

import (
	"fmt"

	"github.com/r3boot/logshipper/lib/config"
	"github.com/r3boot/logshipper/lib/logger"
)

func NewInputs(l *logger.Logger, c *config.Config) error {
	log = l
	cfg = c

	err := SetupSyslog()
	if err != nil {
		return fmt.Errorf("NewInputs: %v", err)
	}

	err = SetupCFL()
	if err != nil {
		return fmt.Errorf("NewInputs: %v", err)
	}

	err = SetupExim()
	if err != nil {
		return fmt.Errorf("NewInputs: %v", err)
	}

	for _, input := range cfg.Inputs {
		switch input.Type {
		case config.T_AMQPBRIDGE:
			{
				as, err := NewAmqpSlurper()
				if err != nil {
					return fmt.Errorf("NewInputs: %v", err)
				}
				InputAmqp = append(InputAmqp, as)
			}
		default:
			{
				mf, err := NewMonitoredFile(
					input.Name,
					input.Path,
					input.Type,
					input.TsFormat,
				)
				if err != nil {
					return fmt.Errorf("NewInputs: %v", err)
				}
				MonitoredFiles = append(MonitoredFiles, mf)
			}
		}
	}

	return nil
}
