package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/BagusAK95/go-boilerplate/internal/config"
	"github.com/golang-migrate/migrate/v4"
	migratepostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

func Open() postgresConnection {
	pgConfig := setConfig()

	return postgresConnection{
		Master: open(pgConfig.Master),
		Slave:  open(pgConfig.Slave),
	}
}

func open(pgConfig postgres.Config) *gorm.DB {
	db, err := gorm.Open(postgres.New(pgConfig), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("❌ Could not to open PostgresSQL connection: %v", err)
	}

	if err := db.Use(tracing.NewPlugin(tracing.WithAttributes())); err != nil {
		log.Fatalf("❌ Could not to use tracing plugin for PostgresSQL: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ Could not to get PostgresSQL connection: %v", err)
	}

	sqlDB.SetMaxOpenConns(config.PostgresConfig.MaxOpenConnections)
	sqlDB.SetMaxIdleConns(config.PostgresConfig.MaxIdleConnections)
	sqlDB.SetConnMaxLifetime(config.PostgresConfig.ConnMaxLifetime)

	if err = sqlDB.Ping(); err != nil {
		log.Fatalf("❌ Could not to ping PostgresSQL database: %v", err)
	}

	if !config.PostgresConfig.AutoMigrate {
		return db
	}

	migrateDriver, err := migratepostgres.WithInstance(sqlDB, &migratepostgres.Config{})
	if err != nil {
		log.Fatalf("❌ Could not to create migrate instance for PostgresSQL:%v", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations/postgres", "postgres", migrateDriver)
	if err != nil {
		log.Fatalf("❌ Could not to create migrate instance for PostgresSQL:%v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("❌ Could not to migrate PostgresSQL:%v", err)
	}

	return db
}

func (c postgresConnection) Shutdown(ctx context.Context) error {
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
