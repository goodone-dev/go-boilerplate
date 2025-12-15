package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/messaging/rabbitmq"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

// RequestHandler is a function that processes RPC requests and returns a response
type RequestHandler func(ctx context.Context, body []byte, headers map[string]any) (any, error)

// Server handles RPC server operations
type Server struct {
	client    rabbitmq.Client
	queueName string
}

// ServerConfig holds RPC server configuration
type ServerConfig struct {
	QueueName string
}

// NewServer creates a new RPC server
func NewServer(client rabbitmq.Client, config ServerConfig) (*Server, error) {
	// Declare RPC queue
	_, err := client.DeclareQueue(rabbitmq.QueueConfig{
		Name:       config.QueueName,
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Args:       nil,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to declare RPC queue: %w", err)
	}

	return &Server{
		client:    client,
		queueName: config.QueueName,
	}, nil
}

// Serve starts serving RPC requests
func (s *Server) Serve(ctx context.Context, handler RequestHandler) error {
	deliveryHandler := func(ctx context.Context, delivery amqp.Delivery) error {
		log.Printf("RPC Server: Received request with correlation ID %s", delivery.CorrelationId)

		// Process the request
		response, err := handler(ctx, delivery.Body, delivery.Headers)
		if err != nil {
			log.Printf("RPC Server: Error processing request: %v", err)
			// Send error response
			return s.sendResponse(ctx, delivery.ReplyTo, delivery.CorrelationId, nil, err)
		}

		// Send success response
		return s.sendResponse(ctx, delivery.ReplyTo, delivery.CorrelationId, response, nil)
	}

	consumeConfig := rabbitmq.ConsumeConfig{
		Queue:     s.queueName,
		Consumer:  "",
		AutoAck:   false,
		Exclusive: false,
		NoLocal:   false,
		NoWait:    false,
		Args:      nil,
	}

	return s.client.Consume(ctx, consumeConfig, deliveryHandler)
}

// ServeJSON serves RPC requests with JSON marshaling/unmarshaling
func (s *Server) ServeJSON(ctx context.Context, handler func(ctx context.Context, request any, headers map[string]any) (any, error), requestType any) error {
	requestHandler := func(ctx context.Context, body []byte, headers map[string]any) (any, error) {
		// Create a new instance of the request type
		request := requestType

		if err := json.Unmarshal(body, &request); err != nil {
			return nil, fmt.Errorf("failed to unmarshal request: %w", err)
		}

		return handler(ctx, request, headers)
	}

	return s.Serve(ctx, requestHandler)
}

func (s *Server) sendResponse(ctx context.Context, replyTo, correlationID string, response any, err error) error {
	var body []byte
	var responseErr error

	if err != nil {
		// Create error response
		errorResponse := map[string]any{
			"error": err.Error(),
		}
		body, responseErr = json.Marshal(errorResponse)
	} else {
		body, responseErr = json.Marshal(response)
	}

	if responseErr != nil {
		return fmt.Errorf("failed to marshal response: %w", responseErr)
	}

	msg := rabbitmq.Message{
		Body:          body,
		ContentType:   "application/json",
		CorrelationID: correlationID,
		MessageID:     uuid.New().String(),
		Timestamp:     time.Now(),
	}

	config := rabbitmq.PublishConfig{
		Exchange:   "", // Default exchange for direct queue publishing
		RoutingKey: replyTo,
		Mandatory:  false,
		Immediate:  false,
	}

	return s.client.Publish(ctx, config, msg)
}

// Close closes the server
func (s *Server) Close() error {
	return nil // Client is shared, don't close it
}
