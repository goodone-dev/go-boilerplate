package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	reconnectDelay = 5 * time.Second
	resendDelay    = 5 * time.Second
)

// client implements the Client interface with connection pooling
type client struct {
	config      Config
	conn        *amqp.Connection
	channels    chan *amqp.Channel
	mu          sync.RWMutex
	closed      bool
	tracer      trace.Tracer
	notifyClose chan *amqp.Error
}

// NewClient creates a new RabbitMQ client with connection pooling
func NewClient(config Config) (Client, error) {
	c := &client{
		config:   config,
		channels: make(chan *amqp.Channel, config.PoolSize),
		tracer:   otel.Tracer("rabbitmq"),
	}

	if err := c.connect(); err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	// Initialize channel pool
	for i := 0; i < config.PoolSize; i++ {
		ch, err := c.conn.Channel()
		if err != nil {
			c.Close()
			return nil, fmt.Errorf("failed to create channel: %w", err)
		}
		c.channels <- ch
	}

	// Handle reconnection
	go c.handleReconnect()

	return c, nil
}

func (c *client) connect() error {
	dsn := fmt.Sprintf(
		"amqp://%s:%s@%s:%s%s",
		c.config.Username,
		c.config.Password,
		c.config.Host,
		c.config.Port,
		c.config.Vhost,
	)

	conn, err := amqp.Dial(dsn)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}

	c.conn = conn
	c.notifyClose = make(chan *amqp.Error)
	c.conn.NotifyClose(c.notifyClose)

	log.Println("RabbitMQ: Connected successfully")
	return nil
}

func (c *client) handleReconnect() {
	for {
		err := <-c.notifyClose
		if c.closed {
			return
		}

		log.Printf("RabbitMQ: Connection closed: %v. Reconnecting...", err)

		for {
			time.Sleep(reconnectDelay)

			if err := c.connect(); err != nil {
				log.Printf("RabbitMQ: Reconnection failed: %v", err)
				continue
			}

			// Recreate channel pool
			c.mu.Lock()
			close(c.channels)
			c.channels = make(chan *amqp.Channel, c.config.PoolSize)
			for i := 0; i < c.config.PoolSize; i++ {
				ch, err := c.conn.Channel()
				if err != nil {
					log.Printf("RabbitMQ: Failed to create channel: %v", err)
					c.mu.Unlock()
					continue
				}
				c.channels <- ch
			}
			c.mu.Unlock()

			log.Println("RabbitMQ: Reconnected successfully")
			break
		}
	}
}

func (c *client) GetChannel() (*amqp.Channel, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return nil, fmt.Errorf("client is closed")
	}

	select {
	case ch := <-c.channels:
		return ch, nil
	case <-time.After(5 * time.Second):
		return nil, fmt.Errorf("timeout waiting for channel")
	}
}

func (c *client) returnChannel(ch *amqp.Channel) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		ch.Close()
		return
	}

	select {
	case c.channels <- ch:
	default:
		ch.Close()
	}
}

func (c *client) DeclareExchange(config ExchangeConfig) error {
	ch, err := c.GetChannel()
	if err != nil {
		return err
	}
	defer c.returnChannel(ch)

	return ch.ExchangeDeclare(
		config.Name,
		string(config.Type),
		config.Durable,
		config.AutoDelete,
		config.Internal,
		config.NoWait,
		config.Args,
	)
}

func (c *client) DeclareQueue(config QueueConfig) (amqp.Queue, error) {
	ch, err := c.GetChannel()
	if err != nil {
		return amqp.Queue{}, err
	}
	defer c.returnChannel(ch)

	return ch.QueueDeclare(
		config.Name,
		config.Durable,
		config.AutoDelete,
		config.Exclusive,
		config.NoWait,
		config.Args,
	)
}

func (c *client) BindQueue(queueName, routingKey, exchangeName string, args amqp.Table) error {
	ch, err := c.GetChannel()
	if err != nil {
		return err
	}
	defer c.returnChannel(ch)

	return ch.QueueBind(queueName, routingKey, exchangeName, false, args)
}

func (c *client) Publish(ctx context.Context, config PublishConfig, msg Message) error {
	ctx, span := c.tracer.Start(ctx, "rabbitmq.publish",
		trace.WithAttributes(
			attribute.String("exchange", config.Exchange),
			attribute.String("routing_key", config.RoutingKey),
		),
	)
	defer span.End()

	ch, err := c.GetChannel()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	defer c.returnChannel(ch)

	// Inject trace context into headers
	if msg.Headers == nil {
		msg.Headers = make(map[string]interface{})
	}
	carrier := NewHeaderCarrier(msg.Headers)
	otel.GetTextMapPropagator().Inject(ctx, carrier)

	publishing := amqp.Publishing{
		ContentType:   msg.ContentType,
		Body:          msg.Body,
		Headers:       msg.Headers,
		Priority:      msg.Priority,
		Expiration:    msg.Expiration,
		MessageId:     msg.MessageID,
		Timestamp:     msg.Timestamp,
		Type:          msg.Type,
		ReplyTo:       msg.ReplyTo,
		CorrelationId: msg.CorrelationID,
		DeliveryMode:  amqp.Persistent, // Make messages persistent
	}

	err = ch.PublishWithContext(
		ctx,
		config.Exchange,
		config.RoutingKey,
		config.Mandatory,
		config.Immediate,
		publishing,
	)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("failed to publish message: %w", err)
	}

	span.SetStatus(codes.Ok, "message published successfully")
	return nil
}

func (c *client) Consume(ctx context.Context, config ConsumeConfig, handler DeliveryHandler) error {
	ch, err := c.GetChannel()
	if err != nil {
		return err
	}

	// Set QoS to limit unacknowledged messages
	if err := ch.Qos(1, 0, false); err != nil {
		c.returnChannel(ch)
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	deliveries, err := ch.Consume(
		config.Queue,
		config.Consumer,
		config.AutoAck,
		config.Exclusive,
		config.NoLocal,
		config.NoWait,
		config.Args,
	)
	if err != nil {
		c.returnChannel(ch)
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		defer c.returnChannel(ch)

		for {
			select {
			case <-ctx.Done():
				log.Printf("RabbitMQ: Consumer stopped for queue %s", config.Queue)
				return
			case delivery, ok := <-deliveries:
				if !ok {
					log.Printf("RabbitMQ: Delivery channel closed for queue %s", config.Queue)
					return
				}

				c.handleDelivery(ctx, delivery, handler)
			}
		}
	}()

	return nil
}

func (c *client) handleDelivery(ctx context.Context, delivery amqp.Delivery, handler DeliveryHandler) {
	// Extract trace context from headers
	carrier := NewHeaderCarrier(delivery.Headers)
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	ctx, span := c.tracer.Start(ctx, "rabbitmq.consume",
		trace.WithAttributes(
			attribute.String("exchange", delivery.Exchange),
			attribute.String("routing_key", delivery.RoutingKey),
			attribute.String("message_id", delivery.MessageId),
		),
	)
	defer span.End()

	retryCount := c.getRetryCount(delivery.Headers)

	err := handler(ctx, delivery)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.Printf("RabbitMQ: Error handling message: %v", err)

		if retryCount < c.config.MaxRetry {
			// Nack and requeue with incremented retry count
			log.Printf("RabbitMQ: Requeuing message (retry %d/%d)", retryCount+1, c.config.MaxRetry)
			delivery.Nack(false, true)
		} else {
			// Max retries reached, reject without requeue (goes to DLX if configured)
			log.Printf("RabbitMQ: Max retries reached, rejecting message")
			delivery.Nack(false, false)
		}
		return
	}

	span.SetStatus(codes.Ok, "message processed successfully")
	delivery.Ack(false)
}

func (c *client) getRetryCount(headers amqp.Table) int {
	if headers == nil {
		return 0
	}

	if count, ok := headers["x-retry-count"].(int32); ok {
		return int(count)
	}

	if count, ok := headers["x-retry-count"].(int); ok {
		return count
	}

	return 0
}

func (c *client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true

	// Close all channels in the pool
	close(c.channels)
	for ch := range c.channels {
		ch.Close()
	}

	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}
