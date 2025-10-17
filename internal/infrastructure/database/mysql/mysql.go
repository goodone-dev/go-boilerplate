package mysql

import (
	"context"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	migratemysql "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/goodone-dev/go-boilerplate/internal/config"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/logger"
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

func Open(ctx context.Context) mysqlConnection {
	mysqlConfig := setConfig()

	return mysqlConnection{
		Master: open(ctx, mysqlConfig.Master),
		Slave:  open(ctx, mysqlConfig.Slave),
	}
}

func open(ctx context.Context, mysqlConfig mysql.Config) *gorm.DB {
	db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		logger.Fatal(ctx, err, "failed to establish mysql connection")
	}

	if err := db.Use(tracing.NewPlugin(tracing.WithAttributes())); err != nil {
		logger.Fatal(ctx, err, "failed to initialize mysql tracing plugin")
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal(ctx, err, "failed to access mysql connection pool")
	}

	sqlDB.SetMaxOpenConns(config.MySQLConfig.MaxOpenConnections)
	sqlDB.SetMaxIdleConns(config.MySQLConfig.MaxIdleConnections)
	sqlDB.SetConnMaxLifetime(config.MySQLConfig.ConnMaxLifetime)

	if err = sqlDB.Ping(); err != nil {
		logger.Fatal(ctx, err, "mysql connection test failed")
	}

	if !config.MySQLConfig.AutoMigrate {
		return db
	}

	migrateDriver, err := migratemysql.WithInstance(sqlDB, &migratemysql.Config{})
	if err != nil {
		logger.Fatal(ctx, err, "failed to initialize mysql migration driver")
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations/mysql", "mysql", migrateDriver)
	if err != nil {
		logger.Fatal(ctx, err, "failed to create migration instance from mysql driver")
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		logger.Fatal(ctx, err, "mysql migration failed")
	}

	return db
}

func (c mysqlConnection) Shutdown(ctx context.Context) error {
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

func (c mysqlConnection) Ping(ctx context.Context) error {
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
