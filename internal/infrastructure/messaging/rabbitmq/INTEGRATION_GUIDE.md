# RabbitMQ Integration Guide

This guide shows how to integrate the RabbitMQ infrastructure into your Go application.

## Step 1: Start RabbitMQ

```bash
# Using docker-compose (recommended)
docker-compose up -d rabbitmq

# Verify it's running
docker ps | grep rabbitmq

# Access Management UI
open http://localhost:15672
# Login: guest / guest
```

## Step 2: Update Configuration

Your `.env` file should already have the RabbitMQ configuration from `.env.example`:

```env
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USERNAME=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_VHOST=/
RABBITMQ_POOL_SIZE=10
RABBITMQ_MAX_RETRY=3
RABBITMQ_RETRY_DELAY=5
```

## Step 3: Initialize in Your Application

### Option A: In main.go

```go
package main

import (
    "context"
    "log"

    "github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq"
    "github.com/spf13/viper"
)

func main() {
    // Load config
    viper.AutomaticEnv()

    // Initialize RabbitMQ client
    rabbitClient, err := rabbitmq.NewClientFromViper()
    if err != nil {
        log.Fatalf("Failed to initialize RabbitMQ: %v", err)
    }
    defer rabbitClient.Close()

    // Pass client to your services
    // ...
}
```

### Option B: As a Dependency Injection

```go
package infrastructure

import (
    "github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq"
)

type Infrastructure struct {
    RabbitMQ rabbitmq.Client
    // ... other infrastructure
}

func NewInfrastructure() (*Infrastructure, error) {
    rabbitClient, err := rabbitmq.NewClientFromViper()
    if err != nil {
        return nil, err
    }

    return &Infrastructure{
        RabbitMQ: rabbitClient,
    }, nil
}

func (i *Infrastructure) Close() error {
    return i.RabbitMQ.Close()
}
```

## Step 4: Implement Your Use Cases

### Example 1: Customer Created Event (Direct Exchange)

**Publisher (in your customer service):**

```go
package customer

import (
    "context"
    "time"

    "github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq"
    "github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq/direct"
)

type Service struct {
    publisher *direct.Publisher
}

func NewService(rabbitClient rabbitmq.Client) (*Service, error) {
    publisher, err := direct.NewPublisher(rabbitClient, "customer.events")
    if err != nil {
        return nil, err
    }

    return &Service{
        publisher: publisher,
    }, nil
}

func (s *Service) CreateCustomer(ctx context.Context, req CreateCustomerRequest) error {
    // Create customer in database
    customer, err := s.repo.Create(ctx, req)
    if err != nil {
        return err
    }

    // Publish event
    event := CustomerCreatedEvent{
        CustomerID: customer.ID,
        Email:      customer.Email,
        Name:       customer.Name,
        CreatedAt:  time.Now(),
    }

    if err := s.publisher.Publish(ctx, "customer.created", event); err != nil {
        log.Printf("Failed to publish event: %v", err)
        // Don't fail the request, just log the error
    }

    return nil
}
```

**Consumer (in your email service):**

```go
package email

import (
    "context"
    "encoding/json"
    "log"

    "github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq"
    "github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq/direct"
)

type Service struct {
    consumer *direct.Consumer
}

func NewService(rabbitClient rabbitmq.Client) (*Service, error) {
    consumer, err := direct.NewConsumer(rabbitClient, direct.ConsumerConfig{
        ExchangeName: "customer.events",
        QueueName:    "email.customer.created",
        RoutingKey:   "customer.created",
        DLXEnabled:   true,
    })
    if err != nil {
        return nil, err
    }

    return &Service{
        consumer: consumer,
    }, nil
}

func (s *Service) Start(ctx context.Context) error {
    return s.consumer.Consume(ctx, func(ctx context.Context, body []byte, headers map[string]any) error {
        var event CustomerCreatedEvent
        if err := json.Unmarshal(body, &event); err != nil {
            return err // Will retry
        }

        // Send welcome email
        if err := s.sendWelcomeEmail(ctx, event); err != nil {
            log.Printf("Failed to send email: %v", err)
            return err // Will retry
        }

        log.Printf("Welcome email sent to %s", event.Email)
        return nil // Success - message will be acknowledged
    })
}
```

### Example 2: Logging System (Topic Exchange)

**Publisher:**

```go
package logger

import (
    "context"

    "github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq"
    "github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq/topic"
)

type MessageBrokerLogger struct {
    publisher *topic.Publisher
}

func NewMessageBrokerLogger(rabbitClient rabbitmq.Client) (*MessageBrokerLogger, error) {
    publisher, err := topic.NewPublisher(rabbitClient, "logs.topic")
    if err != nil {
        return nil, err
    }

    return &MessageBrokerLogger{
        publisher: publisher,
    }, nil
}

func (l *MessageBrokerLogger) LogError(ctx context.Context, service, message string) {
    l.publisher.Publish(ctx, "logs."+service+".error", map[string]any{
        "service": service,
        "level":   "error",
        "message": message,
    })
}

func (l *MessageBrokerLogger) LogInfo(ctx context.Context, service, message string) {
    l.publisher.Publish(ctx, "logs."+service+".info", map[string]any{
        "service": service,
        "level":   "info",
        "message": message,
    })
}
```

**Consumer (all error logs):**

```go
package monitoring

import (
    "context"

    "github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq"
    "github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq/topic"
)

func StartErrorLogConsumer(ctx context.Context, rabbitClient rabbitmq.Client) error {
    consumer, err := topic.NewConsumer(rabbitClient, topic.ConsumerConfig{
        ExchangeName:   "logs.topic",
        QueueName:      "monitoring.errors",
        RoutingPattern: "logs.*.error", // All error logs from any service
        DLXEnabled:     true,
    })
    if err != nil {
        return err
    }

    return consumer.Consume(ctx, func(ctx context.Context, routingKey string, body []byte, headers map[string]any) error {
        // Send to monitoring system (e.g., Sentry, DataDog)
        log.Printf("Error log from %s: %s", routingKey, string(body))
        return nil
    })
}
```

### Example 3: Customer Validation (RPC)

**Server:**

```go
package customer

import (
    "context"

    "github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq"
    "github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq/rpc"
)

type ValidationService struct {
    server *rpc.Server
    repo   CustomerRepository
}

func NewValidationService(rabbitClient rabbitmq.Client, repo CustomerRepository) (*ValidationService, error) {
    server, err := rpc.NewServer(rabbitClient, rpc.ServerConfig{
        QueueName: "customer.validate.rpc",
    })
    if err != nil {
        return nil, err
    }

    return &ValidationService{
        server: server,
        repo:   repo,
    }, nil
}

func (s *ValidationService) Start(ctx context.Context) error {
    return s.server.ServeJSON(ctx, func(ctx context.Context, request any, headers map[string]any) (any, error) {
        req := request.(*ValidateCustomerRequest)

        // Validate customer
        exists, err := s.repo.ExistsByEmail(ctx, req.Email)
        if err != nil {
            return nil, err
        }

        return ValidateCustomerResponse{
            Valid:   !exists,
            Message: "Email is available",
        }, nil
    }, &ValidateCustomerRequest{})
}
```

**Client:**

```go
package api

import (
    "context"
    "time"

    "github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq"
    "github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq/rpc"
)

type CustomerHandler struct {
    rpcClient *rpc.Client
}

func NewCustomerHandler(rabbitClient rabbitmq.Client) (*CustomerHandler, error) {
    rpcClient, err := rpc.NewClient(rabbitClient, rpc.ClientConfig{
        Timeout: 5 * time.Second,
    })
    if err != nil {
        return nil, err
    }

    return &CustomerHandler{
        rpcClient: rpcClient,
    }, nil
}

func (h *CustomerHandler) ValidateEmail(ctx context.Context, email string) (bool, error) {
    request := ValidateCustomerRequest{
        Email: email,
    }

    var response ValidateCustomerResponse
    if err := h.rpcClient.CallJSON(ctx, "customer.validate.rpc", request, &response); err != nil {
        return false, err
    }

    return response.Valid, nil
}
```

## Step 5: Start Consumers

### Option A: In main.go

```go
func main() {
    // ... initialize rabbitClient

    ctx := context.Background()

    // Start consumers
    go func() {
        emailService, _ := email.NewService(rabbitClient)
        if err := emailService.Start(ctx); err != nil {
            log.Printf("Email consumer error: %v", err)
        }
    }()

    go func() {
        if err := monitoring.StartErrorLogConsumer(ctx, rabbitClient); err != nil {
            log.Printf("Monitoring consumer error: %v", err)
        }
    }()

    // Start RPC servers
    go func() {
        validationService, _ := customer.NewValidationService(rabbitClient, customerRepo)
        if err := validationService.Start(ctx); err != nil {
            log.Printf("Validation RPC server error: %v", err)
        }
    }()

    // Start HTTP server
    // ...
}
```

### Option B: Separate Consumer Service

Create a separate binary for consumers:

```go
// cmd/consumer/main.go
package main

func main() {
    rabbitClient, _ := rabbitmq.NewClientFromViper()
    defer rabbitClient.Close()

    ctx := context.Background()

    // Start all consumers
    consumers := []func(context.Context, rabbitmq.Client) error{
        email.StartCustomerCreatedConsumer,
        email.StartOrderCreatedConsumer,
        monitoring.StartErrorLogConsumer,
    }

    for _, consumer := range consumers {
        go func(c func(context.Context, rabbitmq.Client) error) {
            if err := c(ctx, rabbitClient); err != nil {
                log.Printf("Consumer error: %v", err)
            }
        }(consumer)
    }

    // Wait for shutdown signal
    select {}
}
```

## Step 6: Testing

### Unit Tests

```go
package customer_test

import (
    "context"
    "testing"

    "github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq"
)

func TestCustomerService_CreateCustomer(t *testing.T) {
    // Use test RabbitMQ instance or mock
    client, _ := rabbitmq.NewClient(rabbitmq.Config{
        Host:     "localhost",
        Port:     "5672",
        Username: "guest",
        Password: "guest",
        Vhost:    "/",
        PoolSize: 1,
        MaxRetry: 3,
    })
    defer client.Close()

    service, _ := customer.NewService(client)

    // Test customer creation
    err := service.CreateCustomer(context.Background(), CreateCustomerRequest{
        Email: "test@example.com",
        Name:  "Test User",
    })

    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
}
```

## Step 7: Monitoring

### Check Queue Status

```bash
# Using RabbitMQ Management UI
open http://localhost:15672

# Or using CLI
docker exec rabbitmq rabbitmqctl list_queues name messages consumers
```

### Check Dead Letter Queues

```bash
# List all DLQs
docker exec rabbitmq rabbitmqctl list_queues | grep dlq
```

### View Traces

```bash
# Open Jaeger UI
open http://localhost:16686
```

## Troubleshooting

### Consumer Not Receiving Messages

1. Check consumer is running
2. Verify queue bindings in Management UI
3. Check routing key matches
4. Verify exchange type is correct

### Messages Going to DLQ

1. Check application logs for errors
2. Inspect DLQ messages in Management UI
3. Verify message format
4. Check consumer error handling

### High Memory Usage

1. Reduce connection pool size
2. Implement message size limits
3. Monitor queue depths
4. Add consumer instances

## Production Checklist

- [ ] Use strong credentials (not guest/guest)
- [ ] Enable TLS/SSL
- [ ] Set up RabbitMQ cluster
- [ ] Configure monitoring and alerts
- [ ] Implement graceful shutdown
- [ ] Test failure scenarios
- [ ] Document message schemas
- [ ] Set up log aggregation
- [ ] Configure resource limits
- [ ] Implement circuit breakers

## Next Steps

1. Implement your specific use cases
2. Add comprehensive tests
3. Set up monitoring and alerts
4. Document your message schemas
5. Plan for scaling and high availability

For more information, see:
- [README.md](./README.md) - Full documentation
- [QUICKSTART.md](./QUICKSTART.md) - Quick start guide
- [IMPLEMENTATION_SUMMARY.md](./IMPLEMENTATION_SUMMARY.md) - Implementation details
