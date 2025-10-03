package config

import (
	"time"

	"github.com/spf13/viper"
)

var ContextTimeout time.Duration

type ApplicationConfigMap struct {
	Name string `mapstructure:"APP_NAME"`
	Env  string `mapstructure:"APP_ENV"`
	Port int    `mapstructure:"APP_PORT"`
	URL  string `mapstructure:"APP_URL"`
}

var ApplicationConfig ApplicationConfigMap

type RedisConfigMap struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     int    `mapstructure:"REDIS_PORT"`
	TLS      bool   `mapstructure:"REDIS_TLS"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}

var RedisConfig RedisConfigMap

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
	MaxOpenConnections int           `mapstructure:"POSTGRES_MAX_OPEN_CONNECTIONS"`
	MaxIdleConnections int           `mapstructure:"POSTGRES_MAX_IDLE_CONNECTIONS"`
	ConnMaxLifetime    time.Duration `mapstructure:"POSTGRES_CONN_MAX_LIFETIME"`
}

var PostgresConfig PostgresConfigMap

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
	MaxOpenConnections int           `mapstructure:"MYSQL_MAX_OPEN_CONNECTIONS"`
	MaxIdleConnections int           `mapstructure:"MYSQL_MAX_IDLE_CONNECTIONS"`
	ConnMaxLifetime    time.Duration `mapstructure:"MYSQL_CONN_MAX_LIFETIME"`
}

var MySQLConfig MySQLConfigMap

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
	MaxConnPoolSize   int    `mapstructure:"MONGO_MAX_CONN_POOL_SIZE"`
	MinConnPoolSize   int    `mapstructure:"MONGO_MIN_CONN_POOL_SIZE"`
	ConnIdleTimeoutMS int    `mapstructure:"MONGO_CONN_IDLE_TIMEOUT_MS"`
}

var MongoConfig MongoConfigMap

type JaegerConfigMap struct {
	Host        string `mapstructure:"JAEGER_HOST"`
	Port        int    `mapstructure:"JAEGER_PORT"`
	ServiceName string `mapstructure:"JAEGER_SERVICE_NAME"`
}

var JaegerConfig JaegerConfigMap

type MailConfigMap struct {
	Host     string `mapstructure:"MAIL_HOST"`
	Port     int    `mapstructure:"MAIL_PORT"`
	Username string `mapstructure:"MAIL_USERNAME"`
	Password string `mapstructure:"MAIL_PASSWORD"`
	TLS      bool   `mapstructure:"MAIL_TLS"`
}

var MailConfig MailConfigMap

func Load() (err error) {
	viper.AddConfigPath("./")
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
	if err = viper.Unmarshal(&JaegerConfig); err != nil {
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

	ContextTimeout, err = time.ParseDuration(viper.GetString("CONTEXT_TIMEOUT") + "s")

	return
}

func setDefaultConfig() {
	viper.SetDefault("CONTEXT_TIMEOUT", 5)

	viper.SetDefault("POSTGRES_TIMEZONE", "Asia/Jakarta")
	viper.SetDefault("POSTGRES_MAX_OPEN_CONNECTIONS", 10)
	viper.SetDefault("POSTGRES_MAX_IDLE_CONNECTIONS", 10)
	viper.SetDefault("POSTGRES_CONN_MAX_LIFETIME", 300)

	viper.SetDefault("MYSQL_MAX_OPEN_CONNECTIONS", 10)
	viper.SetDefault("MYSQL_MAX_IDLE_CONNECTIONS", 10)
	viper.SetDefault("MYSQL_CONN_MAX_LIFETIME", 300)

	viper.SetDefault("MONGO_MIN_CONN_POOL_SIZE", 2)
	viper.SetDefault("MONGO_MAX_CONN_POOL_SIZE", 100)
	viper.SetDefault("MONGO_CONN_IDLE_TIMEOUT_MS", 60000)
}
