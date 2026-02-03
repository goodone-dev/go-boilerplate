package rabbitmq

import (
	"context"
	"fmt"
	"maps"
	"sync"
	"time"

	"github.com/goodone-dev/go-boilerplate/internal/config"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/logger"
	"github.com/goodone-dev/go-boilerplate/internal/utils/retry"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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
func NewClient(ctx context.Context) Client {
	config := Config{
		Host:       config.RabbitMQ.Host,
		Port:       config.RabbitMQ.Port,
		Username:   config.RabbitMQ.Username,
		Password:   config.RabbitMQ.Password,
		Vhost:      config.RabbitMQ.Vhost,
		PoolSize:   config.RabbitMQ.PoolSize,
		MaxRetry:   config.RabbitMQ.MaxRetry,
		RetryDelay: config.RabbitMQ.RetryDelay,
	}

	c := &client{
		config:   config,
		channels: make(chan *amqp.Channel, config.PoolSize),
		tracer:   otel.Tracer("rabbitmq"),
	}

	_, err := retry.RetryWithBackoff(ctx, "RabbitMQ connection", func() (any, error) {
		return nil, c.connect()
	})
	if err != nil {
		logger.Fatal(ctx, err, "❌ RabbitMQ failed to establish connection after retries").Write()
	}

	// Initialize channel pool
	for i := 0; i < config.PoolSize; i++ {
		ch, err := c.conn.Channel()
		if err != nil {
			_ = c.Shutdown(ctx)

			logger.Fatal(ctx, err, "❌ RabbitMQ failed to create channel").Write()
			return nil
		}
		c.channels <- ch
	}

	// Handle reconnection
	go c.Monitor(ctx)

	return c
}

func (c *client) connect() error {
	dsn := fmt.Sprintf(
		"amqp://%s:%s@%s:%d%s",
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

	return nil
}

func (c *client) Monitor(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-c.notifyClose:
			if c.closed {
				return
			}

			logger.Error(ctx, err, "🛑 RabbitMQ connection lost").Write()
			c.reconnect(ctx)
		}
	}
}

func (c *client) reconnect(ctx context.Context) {
	_, err := retry.RetryWithBackoff(ctx, "RabbitMQ reconnection", func() (any, error) {
		return nil, c.connect()
	})
	if err != nil {
		logger.Fatal(ctx, err, "❌ RabbitMQ failed to establish connection after retries").Write()
	}

	// Recreate channel pool
	c.mu.Lock()
	close(c.channels)
	c.channels = make(chan *amqp.Channel, c.config.PoolSize)
	for i := 0; i < c.config.PoolSize; i++ {
		ch, err := c.conn.Channel()
		if err != nil {
			logger.Error(ctx, err, "❌ RabbitMQ channel creation failed").Write()
			c.mu.Unlock()
			continue
		}
		c.channels <- ch
	}
	c.mu.Unlock()
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
		_ = ch.Close()
		return
	}

	select {
	case c.channels <- ch:
	default:
		_ = ch.Close()
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
	ctx, span := c.tracer.Start(ctx, "RabbitMQ.Publish",
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
		msg.Headers = make(map[string]any)
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
				return
			case delivery, ok := <-deliveries:
				if !ok {
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

	ctx, span := c.tracer.Start(ctx, fmt.Sprintf("RabbitMQ Consume %s", delivery.RoutingKey),
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
		logger.Error(ctx, err, "❌ RabbitMQ error handling message").Write()

		if retryCount < c.config.MaxRetry {
			// Republish with incremented retry count
			retryCount++

			logger.Infof(ctx, "🔁 RabbitMQ retrying message (attempt %d/%d) after %v", retryCount, c.config.MaxRetry, c.config.RetryDelay).Write()
			_ = c.republish(ctx, delivery, retryCount)
		} else {
			// Max retries reached, reject without requeue (goes to DLX if configured)
			logger.Info(ctx, "🚫 RabbitMQ max retries reached, rejecting message").Write()
			_ = delivery.Nack(false, false)
		}
		return
	}

	span.SetStatus(codes.Ok, "message processed successfully")
	_ = delivery.Ack(false)
}

func (c *client) republish(ctx context.Context, delivery amqp.Delivery, retryCount int) error {
	time.Sleep(c.config.RetryDelay)

	// Clone headers and increment retry count
	newHeaders := make(map[string]any)
	maps.Copy(newHeaders, delivery.Headers)
	newHeaders["x-retry-count"] = retryCount

	// Republish message with updated headers
	err := c.Publish(ctx,
		PublishConfig{
			Exchange:   delivery.Exchange,
			RoutingKey: delivery.RoutingKey,
			Mandatory:  false,
			Immediate:  false,
		}, Message{
			Body:          delivery.Body,
			ContentType:   delivery.ContentType,
			Headers:       newHeaders,
			Priority:      delivery.Priority,
			Expiration:    delivery.Expiration,
			MessageID:     delivery.MessageId,
			Timestamp:     delivery.Timestamp,
			Type:          delivery.Type,
			ReplyTo:       delivery.ReplyTo,
			CorrelationID: delivery.CorrelationId,
		},
	)

	if err != nil {
		logger.Error(ctx, err, "❌ RabbitMQ failed to republish retry message, requeuing original").Write()
		_ = delivery.Nack(false, true) // Fallback to simple requeue if publish fails
		return err
	}

	// Ack the original message since we successfully republished it
	return delivery.Ack(false)
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

func (c *client) Ping(ctx context.Context) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return fmt.Errorf("connection is closed")
	}

	if c.conn == nil || c.conn.IsClosed() {
		return fmt.Errorf("connection is not established or closed")
	}

	return nil
}

func (c *client) Shutdown(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true

	// Close all channels in the pool
	close(c.channels)
	for ch := range c.channels {
		_ = ch.Close()
	}

	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}
