package mysql

import (
	"context"
	"fmt"
	"log"

	"github.com/BagusAK95/go-boilerplate/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

func Open() mysqlConnection {
	mysqlConfig := setConfig()

	return mysqlConnection{
		Master: open(mysqlConfig.Master),
		Slave:  open(mysqlConfig.Slave),
	}
}

func open(mysqlConfig mysql.Config) *gorm.DB {
	db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("❌ Could not to open MySQL connection: %v", err)
	}

	if err := db.Use(tracing.NewPlugin(tracing.WithAttributes())); err != nil {
		log.Fatalf("❌ Could not to use tracing plugin for MySQL: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ Could not to get MySQL connection: %v", err)
	}

	sqlDB.SetMaxOpenConns(config.MySQLConfig.MaxOpenConnections)
	sqlDB.SetMaxIdleConns(config.MySQLConfig.MaxIdleConnections)
	sqlDB.SetConnMaxLifetime(config.MySQLConfig.ConnMaxLifetime)

	if err = sqlDB.Ping(); err != nil {
		log.Fatalf("❌ Could not to ping MySQL database: %v", err)
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
