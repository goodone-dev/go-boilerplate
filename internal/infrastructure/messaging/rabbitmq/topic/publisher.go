package topic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/messaging/rabbitmq"
	"github.com/google/uuid"
)

// Publisher handles topic exchange publishing
type Publisher struct {
	client       rabbitmq.Client
	exchangeName string
}

// NewPublisher creates a new topic exchange publisher
func NewPublisher(client rabbitmq.Client, exchangeName string) (*Publisher, error) {
	// Declare the topic exchange
	err := client.DeclareExchange(rabbitmq.ExchangeConfig{
		Name:       exchangeName,
		Type:       rabbitmq.ExchangeTopic,
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
		Args:       nil,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	return &Publisher{
		client:       client,
		exchangeName: exchangeName,
	}, nil
}

// Publish publishes a message to the topic exchange with a routing pattern
// Routing key examples: "logs.error", "events.customer.created", "notifications.email.sent"
func (p *Publisher) Publish(ctx context.Context, routingKey string, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	msg := rabbitmq.Message{
		Body:        body,
		ContentType: "application/json",
		MessageID:   uuid.New().String(),
		Timestamp:   time.Now(),
	}

	config := rabbitmq.PublishConfig{
		Exchange:   p.exchangeName,
		RoutingKey: routingKey,
		Mandatory:  false,
		Immediate:  false,
	}

	return p.client.Publish(ctx, config, msg)
}

// PublishWithPriority publishes a message with priority
func (p *Publisher) PublishWithPriority(ctx context.Context, routingKey string, payload interface{}, priority uint8) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	msg := rabbitmq.Message{
		Body:        body,
		ContentType: "application/json",
		MessageID:   uuid.New().String(),
		Timestamp:   time.Now(),
		Priority:    priority,
	}

	config := rabbitmq.PublishConfig{
		Exchange:   p.exchangeName,
		RoutingKey: routingKey,
		Mandatory:  false,
		Immediate:  false,
	}

	return p.client.Publish(ctx, config, msg)
}

// PublishWithHeaders publishes a message with custom headers
func (p *Publisher) PublishWithHeaders(ctx context.Context, routingKey string, payload interface{}, headers map[string]interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	msg := rabbitmq.Message{
		Body:        body,
		ContentType: "application/json",
		MessageID:   uuid.New().String(),
		Timestamp:   time.Now(),
		Headers:     headers,
	}

	config := rabbitmq.PublishConfig{
		Exchange:   p.exchangeName,
		RoutingKey: routingKey,
		Mandatory:  false,
		Immediate:  false,
	}

	return p.client.Publish(ctx, config, msg)
}

// Close closes the publisher
func (p *Publisher) Close() error {
	return nil // Client is shared, don't close it
}
