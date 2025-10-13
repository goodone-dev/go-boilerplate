package rabbitmq

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

// LoadConfigFromViper loads RabbitMQ configuration from Viper
func LoadConfigFromViper() (Config, error) {
	poolSize, err := strconv.Atoi(viper.GetString("RABBITMQ_POOL_SIZE"))
	if err != nil {
		poolSize = 10 // default
	}

	maxRetry, err := strconv.Atoi(viper.GetString("RABBITMQ_MAX_RETRY"))
	if err != nil {
		maxRetry = 3 // default
	}

	retryDelay, err := strconv.Atoi(viper.GetString("RABBITMQ_RETRY_DELAY"))
	if err != nil {
		retryDelay = 5 // default
	}

	return Config{
		Host:       viper.GetString("RABBITMQ_HOST"),
		Port:       viper.GetString("RABBITMQ_PORT"),
		Username:   viper.GetString("RABBITMQ_USERNAME"),
		Password:   viper.GetString("RABBITMQ_PASSWORD"),
		Vhost:      viper.GetString("RABBITMQ_VHOST"),
		PoolSize:   poolSize,
		MaxRetry:   maxRetry,
		RetryDelay: time.Duration(retryDelay) * time.Second,
	}, nil
}

// NewClientFromViper creates a new RabbitMQ client from Viper configuration
func NewClientFromViper() (Client, error) {
	config, err := LoadConfigFromViper()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return NewClient(config)
}
