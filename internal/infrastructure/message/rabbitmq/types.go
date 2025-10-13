package rabbitmq

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Config holds RabbitMQ configuration
type Config struct {
	Host       string
	Port       string
	Username   string
	Password   string
	Vhost      string
	PoolSize   int
	MaxRetry   int
	RetryDelay time.Duration
}

// ExchangeType defines the type of exchange
type ExchangeType string

const (
	ExchangeDirect ExchangeType = "direct"
	ExchangeTopic  ExchangeType = "topic"
	ExchangeFanout ExchangeType = "fanout"
)

// ExchangeConfig holds exchange configuration
type ExchangeConfig struct {
	Name       string
	Type       ExchangeType
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

// QueueConfig holds queue configuration
type QueueConfig struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

// PublishConfig holds publishing configuration
type PublishConfig struct {
	Exchange   string
	RoutingKey string
	Mandatory  bool
	Immediate  bool
	Headers    map[string]interface{}
	Priority   uint8
	Expiration string
}

// ConsumeConfig holds consumer configuration
type ConsumeConfig struct {
	Queue     string
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
}

// Message represents a message to be published
type Message struct {
	Body        []byte
	ContentType string
	Headers     map[string]interface{}
	Priority    uint8
	Expiration  string
	MessageID   string
	Timestamp   time.Time
	Type        string
	ReplyTo     string
	CorrelationID string
}

// DeliveryHandler is a function that processes delivered messages
type DeliveryHandler func(ctx context.Context, delivery amqp.Delivery) error

// Publisher interface for publishing messages
type Publisher interface {
	Publish(ctx context.Context, config PublishConfig, msg Message) error
	Close() error
}

// Consumer interface for consuming messages
type Consumer interface {
	Consume(ctx context.Context, config ConsumeConfig, handler DeliveryHandler) error
	Close() error
}

// Client is the main RabbitMQ client interface
type Client interface {
	Publisher
	Consumer
	DeclareExchange(config ExchangeConfig) error
	DeclareQueue(config QueueConfig) (amqp.Queue, error)
	BindQueue(queueName, routingKey, exchangeName string, args amqp.Table) error
	GetChannel() (*amqp.Channel, error)
	Close() error
}
