package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/logger"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/messaging/rabbitmq"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Client handles RPC client operations
type Client struct {
	client         rabbitmq.Client
	replyQueueName string
	pending        map[string]chan *amqp.Delivery
	mu             sync.RWMutex
	timeout        time.Duration
}

// ClientConfig holds RPC client configuration
type ClientConfig struct {
	Timeout time.Duration // Default timeout for RPC calls
}

// NewClient creates a new RPC client
func NewClient(ctx context.Context, client rabbitmq.Client, config ClientConfig) *Client {
	// Declare exclusive reply queue
	replyQueue, err := client.DeclareQueue(rabbitmq.QueueConfig{
		Name:       "",
		Durable:    false,
		AutoDelete: true,
		Exclusive:  true,
		NoWait:     false,
		Args:       nil,
	})
	if err != nil {
		logger.Fatalf(ctx, err, "❌ Failed to declare reply queue")
		return nil
	}

	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	rpcClient := &Client{
		client:         client,
		replyQueueName: replyQueue.Name,
		pending:        make(map[string]chan *amqp.Delivery),
		timeout:        config.Timeout,
	}

	// Start consuming replies
	go rpcClient.consumeReplies(ctx)

	return rpcClient
}

func (c *Client) consumeReplies(ctx context.Context) {
	deliveryHandler := func(ctx context.Context, delivery amqp.Delivery) error {
		c.mu.RLock()
		ch, ok := c.pending[delivery.CorrelationId]
		c.mu.RUnlock()

		if ok {
			ch <- &delivery
		}

		return nil
	}

	consumeConfig := rabbitmq.ConsumeConfig{
		Queue:     c.replyQueueName,
		Consumer:  "",
		AutoAck:   true,
		Exclusive: true,
		NoLocal:   false,
		NoWait:    false,
		Args:      nil,
	}

	if err := c.client.Consume(ctx, consumeConfig, deliveryHandler); err != nil {
		logger.Fatalf(ctx, err, "❌ Failed to consume replies")
	}
}

// Call makes an RPC call and waits for the response
func (c *Client) Call(ctx context.Context, queueName string, request any) ([]byte, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	correlationID := uuid.New().String()
	replyChan := make(chan *amqp.Delivery, 1)

	logger.Infof(ctx, "✉️ Making RPC call to queue %s with correlation ID %s", queueName, correlationID)

	c.mu.Lock()
	c.pending[correlationID] = replyChan
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		delete(c.pending, correlationID)
		c.mu.Unlock()
		close(replyChan)
	}()

	msg := rabbitmq.Message{
		Body:          body,
		ContentType:   "application/json",
		CorrelationID: correlationID,
		ReplyTo:       c.replyQueueName,
		MessageID:     uuid.New().String(),
		Timestamp:     time.Now(),
	}

	config := rabbitmq.PublishConfig{
		Exchange:   "", // Default exchange for direct queue publishing
		RoutingKey: queueName,
		Mandatory:  false,
		Immediate:  false,
	}

	if err := c.client.Publish(ctx, config, msg); err != nil {
		return nil, fmt.Errorf("failed to publish RPC request: %w", err)
	}

	// Wait for response with timeout
	select {
	case delivery := <-replyChan:
		// Check if response contains an error
		var errorResponse map[string]any
		if err := json.Unmarshal(delivery.Body, &errorResponse); err == nil {
			if errMsg, ok := errorResponse["error"].(string); ok {
				return nil, fmt.Errorf("RPC error: %s", errMsg)
			}
		}
		return delivery.Body, nil
	case <-time.After(c.timeout):
		return nil, fmt.Errorf("RPC call timeout after %v", c.timeout)
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// CallWithTimeout makes an RPC call with a custom timeout
func (c *Client) CallWithTimeout(ctx context.Context, queueName string, request any, timeout time.Duration) ([]byte, error) {
	oldTimeout := c.timeout
	c.timeout = timeout
	defer func() { c.timeout = oldTimeout }()

	return c.Call(ctx, queueName, request)
}

// CallJSON makes an RPC call and unmarshals the response into the provided type
func (c *Client) CallJSON(ctx context.Context, queueName string, request any, response any) error {
	body, err := c.Call(ctx, queueName, request)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, response); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

// Shutdown closes the client
func (c *Client) Shutdown(ctx context.Context) error {
	return nil // Client is shared, don't close it
}
