package examples

import (
	"context"
	"fmt"
	"log"

	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq/direct"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq/topic"
)

// CustomerCreatedEvent represents a customer created event
type CustomerCreatedEvent struct {
	CustomerID string `json:"customer_id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	CreatedAt  string `json:"created_at"`
}

// CustomerUpdatedEvent represents a customer updated event
type CustomerUpdatedEvent struct {
	CustomerID string `json:"customer_id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	UpdatedAt  string `json:"updated_at"`
}

// CustomerDeletedEvent represents a customer deleted event
type CustomerDeletedEvent struct {
	CustomerID string `json:"customer_id"`
	DeletedAt  string `json:"deleted_at"`
}

// DirectExchangeExample demonstrates direct exchange usage for customer events
func DirectExchangeExample(client rabbitmq.Client) error {
	ctx := context.Background()

	// Create publisher
	publisher, err := direct.NewPublisher(client, "customer.direct")
	if err != nil {
		return fmt.Errorf("failed to create publisher: %w", err)
	}
	defer publisher.Close()

	// Create consumer for customer.created events
	consumerCreated, err := direct.NewConsumer(client, direct.ConsumerConfig{
		ExchangeName: "customer.direct",
		QueueName:    "customer.created.queue",
		RoutingKey:   "customer.created",
		DLXEnabled:   true,
	})
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}
	defer consumerCreated.Close()

	// Start consuming customer.created events
	go func() {
		err := consumerCreated.Consume(ctx, func(ctx context.Context, body []byte, headers map[string]interface{}) error {
			log.Printf("Direct Exchange: Received customer.created event: %s", string(body))
			// Process the event
			// If processing fails, the message will be requeued or sent to DLQ
			return nil
		})
		if err != nil {
			log.Printf("Error consuming messages: %v", err)
		}
	}()

	// Publish customer.created event
	event := CustomerCreatedEvent{
		CustomerID: "cust-123",
		Email:      "john@example.com",
		Name:       "John Doe",
		CreatedAt:  "2025-10-13T18:00:00Z",
	}

	if err := publisher.Publish(ctx, "customer.created", event); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	log.Println("Direct Exchange: Published customer.created event")
	return nil
}

// TopicExchangeExample demonstrates topic exchange usage for customer events
func TopicExchangeExample(client rabbitmq.Client) error {
	ctx := context.Background()

	// Create publisher
	publisher, err := topic.NewPublisher(client, "events.topic")
	if err != nil {
		return fmt.Errorf("failed to create publisher: %w", err)
	}
	defer publisher.Close()

	// Create consumer for all customer events (customer.*)
	consumerAllCustomer, err := topic.NewConsumer(client, topic.ConsumerConfig{
		ExchangeName:   "events.topic",
		QueueName:      "customer.all.queue",
		RoutingPattern: "customer.*",
		DLXEnabled:     true,
	})
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}
	defer consumerAllCustomer.Close()

	// Create consumer for all events (events.#)
	consumerAllEvents, err := topic.NewConsumer(client, topic.ConsumerConfig{
		ExchangeName:   "events.topic",
		QueueName:      "events.all.queue",
		RoutingPattern: "customer.#",
		DLXEnabled:     true,
	})
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}
	defer consumerAllEvents.Close()

	// Start consuming all customer events
	go func() {
		err := consumerAllCustomer.Consume(ctx, func(ctx context.Context, routingKey string, body []byte, headers map[string]interface{}) error {
			log.Printf("Topic Exchange (customer.*): Received event with routing key %s: %s", routingKey, string(body))
			return nil
		})
		if err != nil {
			log.Printf("Error consuming messages: %v", err)
		}
	}()

	// Start consuming all events
	go func() {
		err := consumerAllEvents.Consume(ctx, func(ctx context.Context, routingKey string, body []byte, headers map[string]interface{}) error {
			log.Printf("Topic Exchange (customer.#): Received event with routing key %s: %s", routingKey, string(body))
			return nil
		})
		if err != nil {
			log.Printf("Error consuming messages: %v", err)
		}
	}()

	// Publish various customer events
	createdEvent := CustomerCreatedEvent{
		CustomerID: "cust-456",
		Email:      "jane@example.com",
		Name:       "Jane Smith",
		CreatedAt:  "2025-10-13T18:00:00Z",
	}

	if err := publisher.Publish(ctx, "customer.created", createdEvent); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	updatedEvent := CustomerUpdatedEvent{
		CustomerID: "cust-456",
		Email:      "jane.smith@example.com",
		Name:       "Jane Smith",
		UpdatedAt:  "2025-10-13T18:05:00Z",
	}

	if err := publisher.Publish(ctx, "customer.updated", updatedEvent); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	log.Println("Topic Exchange: Published customer events")
	return nil
}
