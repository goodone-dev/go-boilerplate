package direct

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/logger"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/messaging/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

// MessageHandler is a function that processes messages
type MessageHandler func(ctx context.Context, body []byte, headers map[string]any) error

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
func NewConsumer(ctx context.Context, client rabbitmq.Client, config ConsumerConfig) *Consumer {
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
		logger.Fatal(ctx, err, "❌ RabbitMQ failed to declare exchange")
		return nil
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
			logger.Fatal(ctx, err, "❌ RabbitMQ failed to declare DLX")
			return nil
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
			logger.Fatal(ctx, err, "❌ RabbitMQ failed to declare DLQ")
			return nil
		}

		// Bind DLQ to DLX
		err = client.BindQueue(dlqName, config.RoutingKey, dlxName, nil)
		if err != nil {
			logger.Fatal(ctx, err, "❌ RabbitMQ failed to bind DLQ")
			return nil
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
		logger.Fatal(ctx, err, "❌ RabbitMQ failed to declare queue")
		return nil
	}

	// Bind queue to exchange
	err = client.BindQueue(config.QueueName, config.RoutingKey, config.ExchangeName, nil)
	if err != nil {
		logger.Fatal(ctx, err, "❌ RabbitMQ failed to bind queue")
		return nil
	}

	return consumer
}

// Consume starts consuming messages from the queue
func (c *Consumer) Consume(ctx context.Context, handler MessageHandler) error {
	deliveryHandler := func(ctx context.Context, delivery amqp.Delivery) error {
		logger.Infof(ctx, "📩 RabbitMQ received message from queue %s with routing key %s", c.queueName, delivery.RoutingKey)
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
func (c *Consumer) ConsumeJSON(ctx context.Context, handler func(ctx context.Context, payload any, headers map[string]any) error, payloadType any) error {
	messageHandler := func(ctx context.Context, body []byte, headers map[string]any) error {
		// Use reflection to create a new instance of the payload type
		t := reflect.TypeOf(payloadType)
		if t.Kind() == reflect.Pointer {
			t = t.Elem()
		}

		// Create a pointer to the type to allow json.Unmarshal to fill it
		ptr := reflect.New(t).Interface()

		if err := json.Unmarshal(body, ptr); err != nil {
			return fmt.Errorf("failed to unmarshal message: %w", err)
		}

		// Get the actual value to pass to the handler
		payload := reflect.ValueOf(ptr).Elem().Interface()

		return handler(ctx, payload, headers)
	}

	return c.Consume(ctx, messageHandler)
}

// Shutdown closes the consumer
func (c *Consumer) Shutdown(ctx context.Context) error {
	return nil // Client is shared, don't close it
}
