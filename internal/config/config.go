package config

import (
	"time"

	"github.com/spf13/viper"
)

var ContextTimeout time.Duration
var CorsConfig CorsConfigMap
var ApplicationConfig ApplicationConfigMap
var RedisConfig RedisConfigMap
var PostgresConfig PostgresConfigMap
var MySQLConfig MySQLConfigMap
var MongoConfig MongoConfigMap
var RabbitMQConfig RabbitMQConfigMap
var TracerConfig TracerConfigMap
var LoggerConfig LoggerConfigMap
var MailConfig MailConfigMap
var HttpServerConfig HttpServerConfigMap
var HttpClientConfig HttpClientConfigMap
var CircuitBreakerConfig CircuitBreakerConfigMap
var RateLimiterConfig RateLimiterConfigMap
var IdempotencyDuration time.Duration
var RetryConfig RetryConfigMap

type Environment string

const (
	EnvLocal Environment = "local"
	EnvDev   Environment = "development"
	EnvStag  Environment = "staging"
	EnvProd  Environment = "production"
)

type ApplicationConfigMap struct {
	Name string      `mapstructure:"APP_NAME"`
	Env  Environment `mapstructure:"APP_ENV"`
	Port int         `mapstructure:"APP_PORT"`
	URL  string      `mapstructure:"APP_URL"`
}

type RedisConfigMap struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     int    `mapstructure:"REDIS_PORT"`
	TLS      bool   `mapstructure:"REDIS_TLS"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}

type PostgresConfigMap struct {
	MasterHost         string        `mapstructure:"POSTGRES_MASTER_HOST"`
	MasterUsername     string        `mapstructure:"POSTGRES_MASTER_USERNAME"`
	MasterPassword     string        `mapstructure:"POSTGRES_MASTER_PASSWORD"`
	MasterPort         int           `mapstructure:"POSTGRES_MASTER_PORT"`
	MasterSSLMode      string        `mapstructure:"POSTGRES_MASTER_SSL_MODE"`
	SlaveHost          string        `mapstructure:"POSTGRES_SLAVE_HOST"`
	SlaveUsername      string        `mapstructure:"POSTGRES_SLAVE_USERNAME"`
	SlavePassword      string        `mapstructure:"POSTGRES_SLAVE_PASSWORD"`
	SlavePort          int           `mapstructure:"POSTGRES_SLAVE_PORT"`
	SlaveSSLMode       string        `mapstructure:"POSTGRES_SLAVE_SSL_MODE"`
	Host               string        `mapstructure:"POSTGRES_HOST"`
	Username           string        `mapstructure:"POSTGRES_USERNAME"`
	Password           string        `mapstructure:"POSTGRES_PASSWORD"`
	Port               int           `mapstructure:"POSTGRES_PORT"`
	SSLMode            string        `mapstructure:"POSTGRES_SSL_MODE"`
	Database           string        `mapstructure:"POSTGRES_DATABASE"`
	Timezone           string        `mapstructure:"POSTGRES_TIMEZONE"`
	AutoMigrate        bool          `mapstructure:"POSTGRES_AUTO_MIGRATE"`
	MaxOpenConnections int           `mapstructure:"POSTGRES_MAX_OPEN_CONNECTIONS"`
	MaxIdleConnections int           `mapstructure:"POSTGRES_MAX_IDLE_CONNECTIONS"`
	ConnMaxLifetime    time.Duration `mapstructure:"POSTGRES_CONN_MAX_LIFETIME"`
	InsertBatchSize    int           `mapstructure:"POSTGRES_INSERT_BATCH_SIZE"`
}

type MySQLConfigMap struct {
	MasterHost         string        `mapstructure:"MYSQL_MASTER_HOST"`
	MasterUsername     string        `mapstructure:"MYSQL_MASTER_USERNAME"`
	MasterPassword     string        `mapstructure:"MYSQL_MASTER_PASSWORD"`
	MasterPort         int           `mapstructure:"MYSQL_MASTER_PORT"`
	SlaveHost          string        `mapstructure:"MYSQL_SLAVE_HOST"`
	SlaveUsername      string        `mapstructure:"MYSQL_SLAVE_USERNAME"`
	SlavePassword      string        `mapstructure:"MYSQL_SLAVE_PASSWORD"`
	SlavePort          int           `mapstructure:"MYSQL_SLAVE_PORT"`
	Host               string        `mapstructure:"MYSQL_HOST"`
	Username           string        `mapstructure:"MYSQL_USERNAME"`
	Password           string        `mapstructure:"MYSQL_PASSWORD"`
	Port               int           `mapstructure:"MYSQL_PORT"`
	Database           string        `mapstructure:"MYSQL_DATABASE"`
	AutoMigrate        bool          `mapstructure:"MYSQL_AUTO_MIGRATE"`
	MaxOpenConnections int           `mapstructure:"MYSQL_MAX_OPEN_CONNECTIONS"`
	MaxIdleConnections int           `mapstructure:"MYSQL_MAX_IDLE_CONNECTIONS"`
	ConnMaxLifetime    time.Duration `mapstructure:"MYSQL_CONN_MAX_LIFETIME"`
	InsertBatchSize    int           `mapstructure:"MYSQL_INSERT_BATCH_SIZE"`
}

type MongoConfigMap struct {
	MasterHost        string `mapstructure:"MONGO_MASTER_HOST"`
	MasterPort        int    `mapstructure:"MONGO_MASTER_PORT"`
	MasterUsername    string `mapstructure:"MONGO_MASTER_USERNAME"`
	MasterPassword    string `mapstructure:"MONGO_MASTER_PASSWORD"`
	SlaveHost         string `mapstructure:"MONGO_SLAVE_HOST"`
	SlavePort         int    `mapstructure:"MONGO_SLAVE_PORT"`
	SlaveUsername     string `mapstructure:"MONGO_SLAVE_USERNAME"`
	SlavePassword     string `mapstructure:"MONGO_SLAVE_PASSWORD"`
	Host              string `mapstructure:"MONGO_HOST"`
	Port              int    `mapstructure:"MONGO_PORT"`
	Username          string `mapstructure:"MONGO_USERNAME"`
	Password          string `mapstructure:"MONGO_PASSWORD"`
	Database          string `mapstructure:"MONGO_DATABASE"`
	AutoMigrate       bool   `mapstructure:"MONGO_AUTO_MIGRATE"`
	MaxConnPoolSize   int    `mapstructure:"MONGO_MAX_CONN_POOL_SIZE"`
	MinConnPoolSize   int    `mapstructure:"MONGO_MIN_CONN_POOL_SIZE"`
	ConnIdleTimeoutMS int    `mapstructure:"MONGO_CONN_IDLE_TIMEOUT_MS"`
	InsertBatchSize   int    `mapstructure:"MONGO_INSERT_BATCH_SIZE"`
}

type RabbitMQConfigMap struct {
	Host               string        `mapstructure:"RABBITMQ_HOST"`
	Port               int           `mapstructure:"RABBITMQ_PORT"`
	Username           string        `mapstructure:"RABBITMQ_USERNAME"`
	Password           string        `mapstructure:"RABBITMQ_PASSWORD"`
	Vhost              string        `mapstructure:"RABBITMQ_VHOST"`
	PoolSize           int           `mapstructure:"RABBITMQ_POOL_SIZE"`
	MaxRetry           int           `mapstructure:"RABBITMQ_MAX_RETRY"`
	RetryDelay         time.Duration `mapstructure:"RABBITMQ_RETRY_DELAY"`
	DirectExchangeName string        `mapstructure:"RABBITMQ_DIRECT_EXCHANGE_NAME"`
	TopicExchangeName  string        `mapstructure:"RABBITMQ_TOPIC_EXCHANGE_NAME"`
}

type TracerConfigMap struct {
	Enabled bool   `mapstructure:"TRACER_ENABLED"`
	Host    string `mapstructure:"TRACER_EXPORTER_HOST"`
	Port    int    `mapstructure:"TRACER_EXPORTER_PORT"`
}

type LoggerConfigMap struct {
	Host  string `mapstructure:"LOGGER_EXPORTER_HOST"`
	Port  int    `mapstructure:"LOGGER_EXPORTER_PORT"`
	Level int    `mapstructure:"LOGGER_LEVEL"`
}

type MailConfigMap struct {
	Host     string `mapstructure:"MAIL_HOST"`
	Port     int    `mapstructure:"MAIL_PORT"`
	Username string `mapstructure:"MAIL_USERNAME"`
	Password string `mapstructure:"MAIL_PASSWORD"`
	TLS      bool   `mapstructure:"MAIL_TLS"`
}

type CorsConfigMap struct {
	AllowOrigins []string `mapstructure:"CORS_ALLOW_ORIGINS"`
	AllowMethods []string `mapstructure:"CORS_ALLOW_METHODS"`
}

type HttpServerConfigMap struct {
	ReadTimeout       time.Duration `mapstructure:"HTTP_SERVER_READ_TIMEOUT"`
	ReadHeaderTimeout time.Duration `mapstructure:"HTTP_SERVER_READ_HEADER_TIMEOUT"`
	WriteTimeout      time.Duration `mapstructure:"HTTP_SERVER_WRITE_TIMEOUT"`
	IdleTimeout       time.Duration `mapstructure:"HTTP_SERVER_IDLE_TIMEOUT"`
}

type HttpClientConfigMap struct {
	RetryCount    int           `mapstructure:"HTTP_CLIENT_RETRY_COUNT"`
	RetryWaitTime time.Duration `mapstructure:"HTTP_CLIENT_RETRY_WAIT_TIME"`
}

type CircuitBreakerConfigMap struct {
	MinRequests  int           `mapstructure:"CIRCUIT_BREAKER_MIN_REQUESTS"`
	FailureRatio float64       `mapstructure:"CIRCUIT_BREAKER_FAILURE_RATIO"`
	Timeout      time.Duration `mapstructure:"CIRCUIT_BREAKER_TIMEOUT"`
	MaxRequests  int           `mapstructure:"CIRCUIT_BREAKER_MAX_REQUESTS"`
}

type RateLimiterConfigMap struct {
	SingleLimit    int           `mapstructure:"RATE_LIMITER_SINGLE_LIMIT"`
	SingleDuration time.Duration `mapstructure:"RATE_LIMITER_SINGLE_DURATION"`
	GlobalLimit    int           `mapstructure:"RATE_LIMITER_GLOBAL_LIMIT"`
	GlobalDuration time.Duration `mapstructure:"RATE_LIMITER_GLOBAL_DURATION"`
}

type RetryConfigMap struct {
	MaxRetries     int           `mapstructure:"RETRY_MAX_RETRIES"`
	InitialBackoff time.Duration `mapstructure:"RETRY_INITIAL_BACKOFF"`
	MaxBackoff     time.Duration `mapstructure:"RETRY_MAX_BACKOFF"`
}

func Load() (err error) {
	viper.AddConfigPath("./")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	setDefaultConfig()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	// Unmarshal each section explicitly
	if err = viper.Unmarshal(&ApplicationConfig); err != nil {
		return
	}
	if err = viper.Unmarshal(&PostgresConfig); err != nil {
		return
	}
	if err = viper.Unmarshal(&MailConfig); err != nil {
		return
	}
	if err = viper.Unmarshal(&TracerConfig); err != nil {
		return
	}
	if err = viper.Unmarshal(&RedisConfig); err != nil {
		return
	}
	if err = viper.Unmarshal(&MySQLConfig); err != nil {
		return
	}
	if err = viper.Unmarshal(&MongoConfig); err != nil {
		return
	}
	if err = viper.Unmarshal(&RabbitMQConfig); err != nil {
		return
	}
	if err = viper.Unmarshal(&LoggerConfig); err != nil {
		return
	}
	if err = viper.Unmarshal(&CorsConfig); err != nil {
		return
	}
	if err = viper.Unmarshal(&HttpServerConfig); err != nil {
		return
	}
	if err = viper.Unmarshal(&HttpClientConfig); err != nil {
		return
	}
	if err = viper.Unmarshal(&CircuitBreakerConfig); err != nil {
		return
	}
	if err = viper.Unmarshal(&RateLimiterConfig); err != nil {
		return
	}
	if err = viper.Unmarshal(&RetryConfig); err != nil {
		return
	}

	ContextTimeout = viper.GetDuration("CONTEXT_TIMEOUT")
	IdempotencyDuration = viper.GetDuration("IDEMPOTENCY_DURATION")

	return
}

func setDefaultConfig() {
	// Application defaults
	viper.SetDefault("APP_PORT", 8080)
	viper.SetDefault("APP_ENV", "local")
	viper.SetDefault("CONTEXT_TIMEOUT", "5s")

	// CORS defaults
	viper.SetDefault("CORS_ALLOW_ORIGINS", "*")
	viper.SetDefault("CORS_ALLOW_METHODS", "GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS")

	// PostgreSQL defaults
	viper.SetDefault("POSTGRES_TIMEZONE", "Asia/Jakarta")
	viper.SetDefault("POSTGRES_MAX_OPEN_CONNECTIONS", 10)
	viper.SetDefault("POSTGRES_MAX_IDLE_CONNECTIONS", 10)
	viper.SetDefault("POSTGRES_CONN_MAX_LIFETIME", "300s")
	viper.SetDefault("POSTGRES_INSERT_BATCH_SIZE", 100)

	// MySQL defaults
	viper.SetDefault("MYSQL_MAX_OPEN_CONNECTIONS", 10)
	viper.SetDefault("MYSQL_MAX_IDLE_CONNECTIONS", 10)
	viper.SetDefault("MYSQL_CONN_MAX_LIFETIME", "300s")
	viper.SetDefault("MYSQL_INSERT_BATCH_SIZE", 100)

	// MongoDB defaults
	viper.SetDefault("MONGO_MIN_CONN_POOL_SIZE", 2)
	viper.SetDefault("MONGO_MAX_CONN_POOL_SIZE", 100)
	viper.SetDefault("MONGO_CONN_IDLE_TIMEOUT_MS", 60000)
	viper.SetDefault("MONGO_INSERT_BATCH_SIZE", 100)

	// RabbitMQ defaults
	viper.SetDefault("RABBITMQ_POOL_SIZE", 10)
	viper.SetDefault("RABBITMQ_MAX_RETRY", 3)
	viper.SetDefault("RABBITMQ_RETRY_DELAY", "5s")
	viper.SetDefault("RABBITMQ_DIRECT_EXCHANGE_NAME", "direct.exchange")
	viper.SetDefault("RABBITMQ_TOPIC_EXCHANGE_NAME", "topic.exchange")

	// HTTP Server defaults (in seconds)
	viper.SetDefault("HTTP_SERVER_READ_TIMEOUT", "5s")
	viper.SetDefault("HTTP_SERVER_READ_HEADER_TIMEOUT", "2s")
	viper.SetDefault("HTTP_SERVER_WRITE_TIMEOUT", "10s")
	viper.SetDefault("HTTP_SERVER_IDLE_TIMEOUT", "120s")

	// HTTP Client defaults
	viper.SetDefault("HTTP_CLIENT_RETRY_COUNT", 1)
	viper.SetDefault("HTTP_CLIENT_RETRY_WAIT_TIME", "1s")

	// Circuit Breaker defaults
	viper.SetDefault("CIRCUIT_BREAKER_MIN_REQUESTS", 3)
	viper.SetDefault("CIRCUIT_BREAKER_FAILURE_RATIO", 0.5)
	viper.SetDefault("CIRCUIT_BREAKER_TIMEOUT", "60s")
	viper.SetDefault("CIRCUIT_BREAKER_MAX_REQUESTS", 1)

	// Rate Limiter defaults
	viper.SetDefault("RATE_LIMITER_SINGLE_LIMIT", 60)
	viper.SetDefault("RATE_LIMITER_SINGLE_DURATION", "60s")
	viper.SetDefault("RATE_LIMITER_GLOBAL_LIMIT", 1000)
	viper.SetDefault("RATE_LIMITER_GLOBAL_DURATION", "60s")

	// Idempotency defaults
	viper.SetDefault("IDEMPOTENCY_DURATION", "300s")

	// Retry defaults
	viper.SetDefault("RETRY_MAX_RETRIES", 5)
	viper.SetDefault("RETRY_INITIAL_BACKOFF", "1s")
	viper.SetDefault("RETRY_MAX_BACKOFF", "30s")
}
