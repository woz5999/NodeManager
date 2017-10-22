package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Config represents the application's configuration file.
type Config struct {
	Debug           bool   `envconfig:"DEBUG"`
	AwsRegion       string `envconfig:"AWS_REGION"`
	AwsSqsQueueURL  string `envconfig:"AWS_SQS_QUEUE_URL"`
	LeaveTimeoutSec int    `envconfig:"LEAVE_TIMEOUT_SEC":default:180`
	ConsumerThreads int    `envconfig:"CONSUMER_THREADS":default:5`
}

// GetConfig returns the application configuration specified by the config file.
func GetConfig() (*Config, error) {
	c := &Config{}

	err := envconfig.Process("nodeman", c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
