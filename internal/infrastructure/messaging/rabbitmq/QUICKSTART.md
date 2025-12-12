# RabbitMQ Quick Start Guide

This guide will help you get started with the RabbitMQ infrastructure in under 5 minutes.

## 1. Start RabbitMQ

Using Docker Compose (recommended):

```bash
docker-compose up -d rabbitmq
```

Or using Docker directly:

```bash
docker run -d --name rabbitmq \
  -p 5672:5672 \
  -p 15672:15672 \
  rabbitmq:3.13-management-alpine
```

Access RabbitMQ Management UI at http://localhost:15672 (guest/guest)

## 2. Configure Environment

Add to your `.env` file:

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

## 3. Initialize Client

```go
import "github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq"

// Create client
client, err := rabbitmq.NewClientFromViper()
if err != nil {
    log.Fatal(err)
}
defer client.Close()
```

## 4. Choose Your Pattern

### Option A: Direct Exchange (Simple Routing)

**Use Case**: Send messages to specific queues (e.g., order processing, email notifications)

**Publisher:**
```go
import "github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq/direct"

publisher, _ := direct.NewPublisher(client, "orders.direct")
publisher.Publish(ctx, "order.created", orderData)
```

**Consumer:**
```go
consumer, _ := direct.NewConsumer(client, direct.ConsumerConfig{
    ExchangeName: "orders.direct",
    QueueName:    "order.processing.queue",
    RoutingKey:   "order.created",
    DLXEnabled:   true,
})

consumer.Consume(ctx, func(ctx context.Context, body []byte, headers map[string]interface{}) error {
    // Process message
    return nil
})
```

### Option B: Topic Exchange (Pattern Matching)

**Use Case**: Route messages based on patterns (e.g., logging, event broadcasting)

**Publisher:**
```go
import "github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq/topic"

publisher, _ := topic.NewPublisher(client, "events.topic")
publisher.Publish(ctx, "customer.created", customerData)
publisher.Publish(ctx, "customer.updated", customerData)
```

**Consumer:**
```go
// Listen to all customer events
consumer, _ := topic.NewConsumer(client, topic.ConsumerConfig{
    ExchangeName:   "events.topic",
    QueueName:      "customer.events.queue",
    RoutingPattern: "customer.*", // Matches customer.created, customer.updated, etc.
    DLXEnabled:     true,
})

consumer.Consume(ctx, func(ctx context.Context, routingKey string, body []byte, headers map[string]interface{}) error {
    log.Printf("Received %s event", routingKey)
    return nil
})
```

### Option C: RPC (Request-Response)

**Use Case**: Synchronous operations (e.g., validation, data retrieval)

**Server:**
```go
import "github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq/rpc"

server, _ := rpc.NewServer(client, rpc.ServerConfig{
    QueueName: "customer.get.rpc",
})

server.ServeJSON(ctx, func(ctx context.Context, request interface{}, headers map[string]interface{}) (interface{}, error) {
    req := request.(*GetCustomerRequest)
    // Process request
    return GetCustomerResponse{...}, nil
}, &GetCustomerRequest{})
```

**Client:**
```go
rpcClient, _ := rpc.NewClient(client, rpc.ClientConfig{
    Timeout: 10 * time.Second,
})

var response GetCustomerResponse
err := rpcClient.CallJSON(ctx, "customer.get.rpc", request, &response)
```

## 5. Run Examples

```go
import "github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/rabbitmq/examples"

// Run all examples
examples.RunAllExamples()

// Or run specific examples
examples.RunDirectExchangeOnly()
examples.RunTopicExchangeOnly()
examples.RunRPCOnly()
```

## Key Features

✅ **Automatic Retry**: Failed messages are retried up to `RABBITMQ_MAX_RETRY` times
✅ **Dead Letter Queue**: Messages that fail after max retries go to DLQ for inspection
✅ **Connection Pooling**: Efficient connection management with configurable pool size
✅ **OpenTelemetry Tracing**: Full distributed tracing support out of the box
✅ **Auto Reconnection**: Automatically reconnects on connection failures
✅ **Message Durability**: Messages survive broker restarts

## Error Handling

```go
consumer.Consume(ctx, func(ctx context.Context, body []byte, headers map[string]interface{}) error {
    // Return error to retry
    if err := process(body); err != nil {
        return err // Message will be retried
    }

    // Return nil to acknowledge
    return nil // Message is acknowledged and removed from queue
})
```

## Monitoring

- **RabbitMQ Management UI**: http://localhost:15672
- **Jaeger Tracing**: http://localhost:16686
- **Check Queues**: Monitor queue depths and consumer counts
- **Dead Letter Queues**: Check `*.dlq` queues for failed messages

## Troubleshooting

**Connection refused:**
```bash
# Check if RabbitMQ is running
docker ps | grep rabbitmq

# Check logs
docker logs rabbitmq
```

**Messages not being consumed:**
- Check consumer is running
- Verify queue bindings in Management UI
- Check for errors in application logs

**Messages going to DLQ:**
- Check application logs for error messages
- Inspect DLQ messages in Management UI
- Verify message format matches expected schema

## Next Steps

1. Read the full [README.md](./README.md) for detailed documentation
2. Check [examples/](./examples/) for complete working examples
3. Implement your own publishers and consumers
4. Monitor your queues in RabbitMQ Management UI
5. Set up alerts for DLQ message counts

## Production Checklist

- [ ] Configure appropriate `RABBITMQ_MAX_RETRY` value
- [ ] Enable DLX for all critical queues
- [ ] Set up monitoring for queue depths
- [ ] Configure alerts for DLQ messages
- [ ] Use strong credentials (not guest/guest)
- [ ] Enable TLS for production
- [ ] Set up RabbitMQ cluster for high availability
- [ ] Configure resource limits (memory, disk)
- [ ] Implement graceful shutdown handling
- [ ] Test failure scenarios and recovery
