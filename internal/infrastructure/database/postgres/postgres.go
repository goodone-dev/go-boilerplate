package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	migratepostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/goodone-dev/go-boilerplate/internal/config"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/logger"
	"github.com/goodone-dev/go-boilerplate/internal/utils/retry"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

type postgresConfig struct {
	Master postgres.Config
	Slave  postgres.Config
}

func setConfig() postgresConfig {
	cfg := config.PostgresConfig
	dsn := "host=%s user=%s password=%s dbname=%s port=%d sslmode=%s timezone=%s"

	masterConfig := postgres.Config{
		DSN:                  fmt.Sprintf(dsn, cfg.Host, cfg.Username, cfg.Password, cfg.Database, cfg.Port, cfg.SSLMode, cfg.Timezone),
		PreferSimpleProtocol: true,
	}

	if len(cfg.MasterHost) > 0 {
		masterConfig.DSN = fmt.Sprintf(dsn, cfg.MasterHost, cfg.MasterUsername, cfg.MasterPassword, cfg.Database, cfg.MasterPort, cfg.MasterSSLMode, cfg.Timezone)
	}

	slaveConfig := postgres.Config{
		DSN:                  fmt.Sprintf(dsn, cfg.Host, cfg.Username, cfg.Password, cfg.Database, cfg.Port, cfg.SSLMode, cfg.Timezone),
		PreferSimpleProtocol: true,
	}

	if len(cfg.SlaveHost) > 0 {
		slaveConfig.DSN = fmt.Sprintf(dsn, cfg.SlaveHost, cfg.SlaveUsername, cfg.SlavePassword, cfg.Database, cfg.SlavePort, cfg.SlaveSSLMode, cfg.Timezone)
	}

	return postgresConfig{
		Master: masterConfig,
		Slave:  slaveConfig,
	}
}

type postgresConnection struct {
	Master *gorm.DB
	Slave  *gorm.DB
}

func Open(ctx context.Context) *postgresConnection {
	pgConfig := setConfig()

	conn := &postgresConnection{
		Master: open(ctx, pgConfig.Master),
		Slave:  open(ctx, pgConfig.Slave),
	}

	go conn.Monitor(ctx)

	return conn
}

func open(ctx context.Context, pgConfig postgres.Config) *gorm.DB {
	gormConfig := &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	}

	db, err := retry.RetryWithBackoff(ctx, "PostgreSQL connection", func() (*gorm.DB, error) {
		return gorm.Open(postgres.New(pgConfig), gormConfig)
	})
	if err != nil {
		logger.With().Fatal(ctx, err, "❌ PostgreSQL failed to establish connection after retries")
	}

	if err := db.Use(tracing.NewPlugin(tracing.WithAttributes())); err != nil {
		logger.With().Fatal(ctx, err, "❌ PostgreSQL failed to initialize tracing plugin")
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.With().Fatal(ctx, err, "❌ PostgreSQL failed to access connection pool")
	}

	sqlDB.SetMaxOpenConns(config.PostgresConfig.MaxOpenConnections)
	sqlDB.SetMaxIdleConns(config.PostgresConfig.MaxIdleConnections)
	sqlDB.SetConnMaxLifetime(config.PostgresConfig.ConnMaxLifetime)

	_, err = retry.RetryWithBackoff(ctx, "PostgreSQL connection test", func() (any, error) {
		return nil, sqlDB.Ping()
	})
	if err != nil {
		logger.With().Fatal(ctx, err, "❌ PostgreSQL connection test failed")
	}

	if !config.PostgresConfig.AutoMigrate {
		return db
	}

	migrateDriver, err := migratepostgres.WithInstance(sqlDB, &migratepostgres.Config{})
	if err != nil {
		logger.With().Fatal(ctx, err, "❌ PostgreSQL failed to initialize migration driver")
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations/postgres", "postgres", migrateDriver)
	if err != nil {
		logger.With().Fatal(ctx, err, "❌ PostgreSQL failed to create migration instance")
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		logger.With().Fatal(ctx, err, "❌ PostgreSQL failed migration")
	}

	return db
}

func (c *postgresConnection) Shutdown(ctx context.Context) error {
	if err := close(c.Master); err != nil {
		return err
	}

	if err := close(c.Slave); err != nil {
		return err
	}

	return nil
}

func close(conn *gorm.DB) error {
	sqlDB, err := conn.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

func (c *postgresConnection) Ping(ctx context.Context) error {
	masterDB, err := c.Master.DB()
	if err != nil {
		return err
	} else if err := masterDB.Ping(); err != nil {
		return err
	}

	slaveDB, err := c.Slave.DB()
	if err != nil {
		return err
	} else if err := slaveDB.Ping(); err != nil {
		return err
	}

	return nil
}

func (c *postgresConnection) Monitor(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	var wasLost bool

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := c.Ping(ctx)
			if err != nil {
				if !wasLost {
					logger.With().Errorf(ctx, err, "🛑 PostgreSQL connection lost")
					wasLost = true
				}
			} else {
				if wasLost {
					logger.With().Info(ctx, "✅ PostgreSQL connection restored")
					wasLost = false
				}
			}
		}
	}
}
