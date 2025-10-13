package rabbitmq

import "errors"

var (
	// ErrClientClosed is returned when attempting to use a closed client
	ErrClientClosed = errors.New("rabbitmq: client is closed")

	// ErrChannelTimeout is returned when waiting for a channel times out
	ErrChannelTimeout = errors.New("rabbitmq: timeout waiting for channel")

	// ErrInvalidConfig is returned when the configuration is invalid
	ErrInvalidConfig = errors.New("rabbitmq: invalid configuration")

	// ErrConnectionFailed is returned when connection to RabbitMQ fails
	ErrConnectionFailed = errors.New("rabbitmq: connection failed")

	// ErrPublishFailed is returned when publishing a message fails
	ErrPublishFailed = errors.New("rabbitmq: publish failed")

	// ErrConsumeFailed is returned when consuming messages fails
	ErrConsumeFailed = errors.New("rabbitmq: consume failed")

	// ErrRPCTimeout is returned when an RPC call times out
	ErrRPCTimeout = errors.New("rabbitmq: RPC call timeout")

	// ErrRPCFailed is returned when an RPC call fails
	ErrRPCFailed = errors.New("rabbitmq: RPC call failed")
)
