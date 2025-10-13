package direct

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

// MessageHandler is a function that processes messages
type MessageHandler func(ctx context.Context, body []byte, headers map[string]interface{}) error

// Consumer handles direct exchange consumption
type Consumer struct {
	client       rabbitmq.Client
	exchangeName string
	queueName    string
	routingKey   string
	dlxName      string
}

// ConsumerConfig holds consumer configuration
type ConsumerConfig struct {
	ExchangeName string
	QueueName    string
	RoutingKey   string
	DLXEnabled   bool // Enable Dead Letter Exchange
}

// NewConsumer creates a new direct exchange consumer with DLX support
func NewConsumer(client rabbitmq.Client, config ConsumerConfig) (*Consumer, error) {
	// Declare the direct exchange
	err := client.DeclareExchange(rabbitmq.ExchangeConfig{
		Name:       config.ExchangeName,
		Type:       rabbitmq.ExchangeDirect,
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
		Args:       nil,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	consumer := &Consumer{
		client:       client,
		exchangeName: config.ExchangeName,
		queueName:    config.QueueName,
		routingKey:   config.RoutingKey,
	}

	// Setup Dead Letter Exchange if enabled
	if config.DLXEnabled {
		dlxName := config.ExchangeName + ".dlx"
		dlqName := config.QueueName + ".dlq"

		// Declare DLX
		err = client.DeclareExchange(rabbitmq.ExchangeConfig{
			Name:       dlxName,
			Type:       rabbitmq.ExchangeDirect,
			Durable:    true,
			AutoDelete: false,
			Internal:   false,
			NoWait:     false,
			Args:       nil,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to declare DLX: %w", err)
		}

		// Declare Dead Letter Queue
		_, err = client.DeclareQueue(rabbitmq.QueueConfig{
			Name:       dlqName,
			Durable:    true,
			AutoDelete: false,
			Exclusive:  false,
			NoWait:     false,
			Args:       nil,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to declare DLQ: %w", err)
		}

		// Bind DLQ to DLX
		err = client.BindQueue(dlqName, config.RoutingKey, dlxName, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to bind DLQ: %w", err)
		}

		consumer.dlxName = dlxName
	}

	// Declare main queue with DLX configuration
	queueArgs := amqp.Table{}
	if config.DLXEnabled {
		queueArgs["x-dead-letter-exchange"] = consumer.dlxName
		queueArgs["x-dead-letter-routing-key"] = config.RoutingKey
	}

	_, err = client.DeclareQueue(rabbitmq.QueueConfig{
		Name:       config.QueueName,
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Args:       queueArgs,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	// Bind queue to exchange
	err = client.BindQueue(config.QueueName, config.RoutingKey, config.ExchangeName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	return consumer, nil
}

// Consume starts consuming messages from the queue
func (c *Consumer) Consume(ctx context.Context, handler MessageHandler) error {
	deliveryHandler := func(ctx context.Context, delivery amqp.Delivery) error {
		log.Printf("Direct Consumer: Received message from queue %s", c.queueName)
		return handler(ctx, delivery.Body, delivery.Headers)
	}

	consumeConfig := rabbitmq.ConsumeConfig{
		Queue:     c.queueName,
		Consumer:  "",
		AutoAck:   false,
		Exclusive: false,
		NoLocal:   false,
		NoWait:    false,
		Args:      nil,
	}

	return c.client.Consume(ctx, consumeConfig, deliveryHandler)
}

// ConsumeJSON consumes messages and unmarshals them into the provided type
func (c *Consumer) ConsumeJSON(ctx context.Context, handler func(ctx context.Context, payload interface{}, headers map[string]interface{}) error, payloadType interface{}) error {
	messageHandler := func(ctx context.Context, body []byte, headers map[string]interface{}) error {
		// Create a new instance of the payload type
		payload := payloadType

		if err := json.Unmarshal(body, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal message: %w", err)
		}

		return handler(ctx, payload, headers)
	}

	return c.Consume(ctx, messageHandler)
}

// Close closes the consumer
func (c *Consumer) Close() error {
	return nil // Client is shared, don't close it
}
