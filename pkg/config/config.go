package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Config represents the application's configuration file.
type Config struct {
	Debug                       bool   `envconfig:"DEBUG" required:"true"`
	AwsRegion                   string `envconfig:"AWS_REGION" required:"true"`
	AwsSqsQueueURL              string `envconfig:"AWS_SQS_QUEUE_URL" required:"true"`
	ConsumerThreads             int    `envconfig:"CONSUMER_THREADS" default:"5"`
	ErrorVisibilityTimeoutSec   int64  `envconfig:"ERROR_VISIBILITY_TIMEOUT_SEC" default:"60"`
	DefaultVisibilityTimeoutSec int64  `envconfig:"DEFAULT_VISIBILITY_TIMEOUT_SEC" default:"300"`
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
