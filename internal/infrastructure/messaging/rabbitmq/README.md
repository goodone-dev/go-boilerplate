# RabbitMQ Infrastructure

Advanced RabbitMQ implementation with support for Direct, Topic, and RPC exchange patterns. This implementation includes connection pooling, message durability, Dead Letter Exchange (DLX), automatic retry logic, and OpenTelemetry tracing integration.

## Features

- ✅ **Connection Pooling**: Efficient connection management with configurable pool size
- ✅ **Message Durability**: Messages persist across broker restarts
- ✅ **Dead Letter Exchange (DLX)**: Failed messages are routed to DLQ after max retries
- ✅ **Automatic Retry**: Configurable retry attempts with message requeuing
- ✅ **OpenTelemetry Tracing**: Full distributed tracing support
- ✅ **Auto Reconnection**: Automatic reconnection on connection failures
- ✅ **Multiple Exchange Types**: Direct, Topic, and RPC patterns

## Configuration

Add the following to your `.env` file:

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

## Usage

### Initialize Client

```go
import "github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq"

// Create client from Viper configuration
client, err := rabbitmq.NewClientFromViper()
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// Or create with custom config
config := rabbitmq.Config{
    Host:       "localhost",
    Port:       "5672",
    Username:   "guest",
    Password:   "guest",
    Vhost:      "/",
    PoolSize:   10,
    MaxRetry:   3,
    RetryDelay: 5 * time.Second,
}
client, err := rabbitmq.NewClient(config)
```

## Direct Exchange

Direct exchange routes messages to queues based on exact routing key matches.

### Publisher

```go
import "github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq/direct"

// Create publisher
publisher, err := direct.NewPublisher(client, "customer.direct")
if err != nil {
    log.Fatal(err)
}
defer publisher.Close()

// Publish message
event := CustomerCreatedEvent{
    CustomerID: "cust-123",
    Email:      "john@example.com",
    Name:       "John Doe",
}

err = publisher.Publish(ctx, "customer.created", event)
```

### Consumer

```go
import "github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq/direct"

// Create consumer with DLX enabled
consumer, err := direct.NewConsumer(client, direct.ConsumerConfig{
    ExchangeName: "customer.direct",
    QueueName:    "customer.created.queue",
    RoutingKey:   "customer.created",
    DLXEnabled:   true, // Enable Dead Letter Exchange
})
if err != nil {
    log.Fatal(err)
}
defer consumer.Close()

// Start consuming
err = consumer.Consume(ctx, func(ctx context.Context, body []byte, headers map[string]interface{}) error {
    log.Printf("Received: %s", string(body))

    // Process message
    // If error is returned, message will be retried or sent to DLQ
    return nil
})
```

## Topic Exchange

Topic exchange routes messages based on routing key patterns using wildcards.

- `*` (star) matches exactly one word
- `#` (hash) matches zero or more words

### Publisher

```go
import "github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq/topic"

// Create publisher
publisher, err := topic.NewPublisher(client, "events.topic")
if err != nil {
    log.Fatal(err)
}
defer publisher.Close()

// Publish with routing pattern
err = publisher.Publish(ctx, "customer.created", event)
err = publisher.Publish(ctx, "customer.updated", event)
err = publisher.Publish(ctx, "order.created", event)
```

### Consumer

```go
import "github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq/topic"

// Consumer for all customer events (customer.*)
consumer, err := topic.NewConsumer(client, topic.ConsumerConfig{
    ExchangeName:   "events.topic",
    QueueName:      "customer.all.queue",
    RoutingPattern: "customer.*", // Matches customer.created, customer.updated, etc.
    DLXEnabled:     true,
})

// Consumer for all events (customer.#)
consumerAll, err := topic.NewConsumer(client, topic.ConsumerConfig{
    ExchangeName:   "events.topic",
    QueueName:      "events.all.queue",
    RoutingPattern: "customer.#", // Matches customer.created, customer.updated.profile, etc.
    DLXEnabled:     true,
})

// Start consuming
err = consumer.Consume(ctx, func(ctx context.Context, routingKey string, body []byte, headers map[string]interface{}) error {
    log.Printf("Received from %s: %s", routingKey, string(body))
    return nil
})
```

## RPC Pattern

RPC (Remote Procedure Call) pattern for synchronous request-response communication.

### RPC Server

```go
import "github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq/rpc"

// Create RPC server
server, err := rpc.NewServer(client, rpc.ServerConfig{
    QueueName: "customer.get.rpc",
})
if err != nil {
    log.Fatal(err)
}
defer server.Close()

// Serve requests
err = server.ServeJSON(ctx, func(ctx context.Context, request interface{}, headers map[string]interface{}) (interface{}, error) {
    req := request.(*GetCustomerRequest)

    // Process request
    response := GetCustomerResponse{
        CustomerID: req.CustomerID,
        Email:      "customer@example.com",
        Name:       "Customer Name",
    }

    return response, nil
}, &GetCustomerRequest{})
```

### RPC Client

```go
import "github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq/rpc"

// Create RPC client
rpcClient, err := rpc.NewClient(client, rpc.ClientConfig{
    Timeout: 10 * time.Second,
})
if err != nil {
    log.Fatal(err)
}
defer rpcClient.Close()

// Make RPC call
request := GetCustomerRequest{
    CustomerID: "cust-789",
}

var response GetCustomerResponse
err = rpcClient.CallJSON(ctx, "customer.get.rpc", request, &response)
if err != nil {
    log.Fatal(err)
}

log.Printf("Response: %+v", response)
```

## Error Handling & Retry Logic

The infrastructure automatically handles errors with the following behavior:

1. **On Error**: Message is Nack'd and requeued
2. **Retry Count**: Tracked in message headers (`x-retry-count`)
3. **Max Retries**: After reaching `RABBITMQ_MAX_RETRY`, message is sent to DLQ
4. **Dead Letter Queue**: Failed messages are preserved for manual inspection

### Example Error Handling

```go
err = consumer.Consume(ctx, func(ctx context.Context, body []byte, headers map[string]interface{}) error {
    // If this returns an error, the message will be retried
    if err := processMessage(body); err != nil {
        return fmt.Errorf("processing failed: %w", err)
    }

    // Success - message will be acknowledged
    return nil
})
```

## OpenTelemetry Tracing

All publish and consume operations are automatically traced with OpenTelemetry:

- **Trace Context Propagation**: Trace context is injected into message headers
- **Span Attributes**: Exchange, routing key, and message ID are recorded
- **Error Recording**: Errors are automatically recorded in spans

```go
// Tracing is automatic, just use context
ctx, span := tracer.Start(context.Background(), "my-operation")
defer span.End()

// Publish with tracing
err := publisher.Publish(ctx, "customer.created", event)

// Consume with tracing
err = consumer.Consume(ctx, handler)
```

## Best Practices

1. **Always Enable DLX**: Set `DLXEnabled: true` to prevent message loss
2. **Set Appropriate Timeouts**: Configure RPC timeout based on operation complexity
3. **Use Durable Queues**: All queues are durable by default
4. **Monitor DLQ**: Regularly check dead letter queues for failed messages
5. **Connection Pooling**: Use shared client instance across your application
6. **Graceful Shutdown**: Always call `Close()` on publishers/consumers

## Folder Structure

```
internal/infrastructure/message/rabbitmq/
├── client.go           # Main client with connection pooling
├── types.go            # Type definitions and interfaces
├── tracing.go          # OpenTelemetry integration
├── factory.go          # Factory functions for easy initialization
├── README.md           # This file
├── direct/             # Direct exchange implementation
│   ├── publisher.go
│   └── consumer.go
├── topic/              # Topic exchange implementation
│   ├── publisher.go
│   └── consumer.go
├── rpc/                # RPC pattern implementation
│   ├── server.go
│   └── client.go
└── examples/           # Usage examples
    ├── customer_events.go
    └── customer_rpc.go
```

## Running Examples

See `examples/` directory for complete working examples:

- `customer_events.go`: Direct and Topic exchange examples
- `customer_rpc.go`: RPC server and client examples

## Docker Compose Setup

Add RabbitMQ to your `docker-compose.yml`:

```yaml
rabbitmq:
  image: rabbitmq:3-management
  ports:
    - "5672:5672"
    - "15672:15672"
  environment:
    RABBITMQ_DEFAULT_USER: guest
    RABBITMQ_DEFAULT_PASS: guest
  volumes:
    - rabbitmq_data:/var/lib/rabbitmq

volumes:
  rabbitmq_data:
```

Access RabbitMQ Management UI at http://localhost:15672 (guest/guest)
