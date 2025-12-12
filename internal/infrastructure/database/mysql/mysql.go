package mysql

import (
	"context"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	migratemysql "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/goodone-dev/go-boilerplate/internal/config"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/logger"
	"github.com/goodone-dev/go-boilerplate/internal/utils/retry"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

type mysqlConfig struct {
	Master mysql.Config
	Slave  mysql.Config
}

func setConfig() mysqlConfig {
	cfg := config.MySQLConfig
	dsn := "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local"

	masterConfig := mysql.Config{
		DSN: fmt.Sprintf(dsn, cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database),
	}

	if len(cfg.MasterHost) > 0 {
		masterConfig.DSN = fmt.Sprintf(dsn, cfg.MasterUsername, cfg.MasterPassword, cfg.MasterHost, cfg.MasterPort, cfg.Database)
	}

	slaveConfig := mysql.Config{
		DSN: fmt.Sprintf(dsn, cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database),
	}

	if len(cfg.SlaveHost) > 0 {
		slaveConfig.DSN = fmt.Sprintf(dsn, cfg.SlaveUsername, cfg.SlavePassword, cfg.SlaveHost, cfg.SlavePort, cfg.Database)
	}

	return mysqlConfig{
		Master: masterConfig,
		Slave:  slaveConfig,
	}
}

type mysqlConnection struct {
	Master *gorm.DB
	Slave  *gorm.DB
}

func Open(ctx context.Context) *mysqlConnection {
	mysqlConfig := setConfig()

	conn := &mysqlConnection{
		Master: open(ctx, mysqlConfig.Master),
		Slave:  open(ctx, mysqlConfig.Slave),
	}

	go conn.Monitor(ctx)

	return conn
}

func open(ctx context.Context, mysqlConfig mysql.Config) *gorm.DB {
	gormConfig := &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	}

	db, err := retry.RetryWithBackoff(ctx, "MySQL connection", func() (*gorm.DB, error) {
		return gorm.Open(mysql.New(mysqlConfig), gormConfig)
	})
	if err != nil {
		logger.Fatal(ctx, err, "❌ Failed to establish MySQL connection after retries")
	}

	if err := db.Use(tracing.NewPlugin(tracing.WithAttributes())); err != nil {
		logger.Fatal(ctx, err, "❌ Failed to initialize MySQL tracing plugin")
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal(ctx, err, "❌ Failed to access MySQL connection pool")
	}

	sqlDB.SetMaxOpenConns(config.MySQLConfig.MaxOpenConnections)
	sqlDB.SetMaxIdleConns(config.MySQLConfig.MaxIdleConnections)
	sqlDB.SetConnMaxLifetime(config.MySQLConfig.ConnMaxLifetime)

	_, err = retry.RetryWithBackoff(ctx, "MySQL connection test", func() (any, error) {
		return nil, sqlDB.Ping()
	})
	if err != nil {
		logger.Fatal(ctx, err, "❌ MySQL connection test failed")
	}

	if !config.MySQLConfig.AutoMigrate {
		return db
	}

	migrateDriver, err := migratemysql.WithInstance(sqlDB, &migratemysql.Config{})
	if err != nil {
		logger.Fatal(ctx, err, "❌ Failed to initialize MySQL migration driver")
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations/mysql", "mysql", migrateDriver)
	if err != nil {
		logger.Fatal(ctx, err, "❌ Failed to create migration instance from MySQL driver")
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		logger.Fatal(ctx, err, "❌ MySQL migration failed")
	}

	return db
}

func (c *mysqlConnection) Shutdown(ctx context.Context) error {
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

func (c *mysqlConnection) Ping(ctx context.Context) error {
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

func (c *mysqlConnection) Monitor(ctx context.Context) {
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
					logger.Errorf(ctx, err, "🛑 MySQL connection lost")
					wasLost = true
				}
			} else {
				if wasLost {
					logger.Info(ctx, "✅ MySQL connection restored")
					wasLost = false
				}
			}
		}
	}
}
