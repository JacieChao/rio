package create

import (
	"time"

	client "github.com/rancher/rio/types/client/rio/v1"
)

func populateHealthCheck(c *Create, service *client.Service) error {
	var err error

	hc := &client.HealthConfig{
		HealthyThreshold:   int64(c.HealthRetries),
		UnhealthyThreshold: int64(c.UnhealthyRetries),
	}

	hc.InitialDelaySeconds, err = ParseDurationUnit(c.HealthStartPeriod, "--health-start-period", time.Second)
	if err != nil {
		return err
	}

	hc.IntervalSeconds, err = ParseDurationUnit(c.HealthInterval, "--health-interval", time.Second)
	if err != nil {
		return err
	}

	if len(c.HealthCmd) > 0 {
		hc.Test = []string{"CMD-SHELL", c.HealthCmd}
	}

	if len(c.HealthURL) > 0 {
		hc.Test = []string{c.HealthURL}
	}

	hc.TimeoutSeconds, err = ParseDurationUnit(c.HealthTimeout, "--health-timeout", time.Second)
	if err != nil {
		return err
	}

	if len(c.HealthCmd) > 0 || len(c.HealthURL) > 0 {
		service.Healthcheck = hc
	}

	return populateReadyCheck(c, service)
}

func populateReadyCheck(c *Create, service *client.Service) error {
	var err error

	hc := &client.HealthConfig{
		HealthyThreshold:   int64(c.ReadyRetries),
		UnhealthyThreshold: int64(c.UnreadyRetries),
	}

	hc.InitialDelaySeconds, err = ParseDurationUnit(c.ReadyStartPeriod, "--ready-start-period", time.Second)
	if err != nil {
		return err
	}

	hc.IntervalSeconds, err = ParseDurationUnit(c.ReadyInterval, "--ready-interval", time.Second)
	if err != nil {
		return err
	}

	if len(c.ReadyCmd) > 0 {
		hc.Test = []string{"CMD-SHELL", c.ReadyCmd}
	}

	if len(c.ReadyURL) > 0 {
		hc.Test = []string{c.ReadyURL}
	}

	hc.TimeoutSeconds, err = ParseDurationUnit(c.ReadyTimeout, "--ready-timeout", time.Second)

	if len(c.ReadyCmd) > 0 || len(c.ReadyURL) > 0 {
		service.Readycheck = hc
	}

	return err
}
