package mongodb

import (
	"context"
	"fmt"
	"net/url"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/goodone-dev/go-boilerplate/internal/config"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/logger"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type mongoConfig struct {
	Master *options.ClientOptions
	Slave  *options.ClientOptions
}

func setConfig() mongoConfig {
	uriMaster := url.URL{
		Scheme: "mongodb",
		Host:   fmt.Sprintf("%s:%d", config.MongoConfig.Host, config.MongoConfig.Port),
		User:   url.UserPassword(config.MongoConfig.Username, config.MongoConfig.Password),
	}

	if len(config.MongoConfig.MasterHost) > 0 {
		uriMaster = url.URL{
			Scheme: "mongodb",
			Host:   fmt.Sprintf("%s:%d", config.MongoConfig.MasterHost, config.MongoConfig.MasterPort),
			User:   url.UserPassword(config.MongoConfig.MasterUsername, config.MongoConfig.MasterPassword),
		}
	}

	uriSlave := url.URL{
		Scheme: "mongodb",
		Host:   fmt.Sprintf("%s:%d", config.MongoConfig.Host, config.MongoConfig.Port),
		User:   url.UserPassword(config.MongoConfig.Username, config.MongoConfig.Password),
	}

	if len(config.MongoConfig.SlaveHost) > 0 {
		uriSlave = url.URL{
			Scheme: "mongodb",
			Host:   fmt.Sprintf("%s:%d", config.MongoConfig.SlaveHost, config.MongoConfig.SlavePort),
			User:   url.UserPassword(config.MongoConfig.SlaveUsername, config.MongoConfig.SlavePassword),
		}
	}

	return mongoConfig{
		Master: options.Client().ApplyURI(uriMaster.String()),
		Slave:  options.Client().ApplyURI(uriSlave.String()),
	}
}

type mongoConnection struct {
	Master *mongo.Database
	Slave  *mongo.Database
}

func Open(ctx context.Context) *mongoConnection {
	mongoConfig := setConfig()

	return &mongoConnection{
		Master: open(ctx, mongoConfig.Master, readpref.Primary()),
		Slave:  open(ctx, mongoConfig.Slave, readpref.Secondary()),
	}
}

func open(ctx context.Context, opts *options.ClientOptions, rp *readpref.ReadPref) *mongo.Database {
	// FIXME: mongodb monitor
	// opts.SetMonitor(otelmongo.NewMonitor())
	opts.SetDirect(true)
	opts.SetRetryWrites(false)
	opts.SetMaxConnIdleTime(time.Duration(config.MongoConfig.ConnIdleTimeoutMS) * time.Millisecond)
	opts.SetBSONOptions(&options.BSONOptions{
		UseLocalTimeZone: true,
	})
	if config.MongoConfig.MaxConnPoolSize >= 0 {
		opts.SetMaxPoolSize(uint64(config.MongoConfig.MaxConnPoolSize))
	}
	if config.MongoConfig.MinConnPoolSize >= 0 {
		opts.SetMinPoolSize(uint64(config.MongoConfig.MinConnPoolSize))
	}

	client, err := mongo.Connect(opts)
	if err != nil {
		logger.Fatal(ctx, err, "❌ Failed to establish MongoDB connection")
	}

	if err := client.Ping(ctx, rp); err != nil {
		logger.Fatal(ctx, err, "❌ MongoDB connection test failed")
	}

	mongoDB := client.Database(config.MongoConfig.Database)
	if !config.MongoConfig.AutoMigrate {
		return mongoDB
	}

	// FIXME: mongodb migration
	// migrateDriver, err := migratemongo.WithInstance(client, &migratemongo.Config{})
	// if err != nil {
	// 	logger.Fatal(ctx, err, "❌️ Could Not Create Migrate Instance For MongoDB")
	// }

	// m, err := migrate.NewWithDatabaseInstance("file://migrations/mongodb", "mongodb", migrateDriver)
	// if err != nil {
	// 	logger.Fatal(ctx, err, "❌️ Could Not Create Migrate Instance For MongoDB")
	// }

	// err = m.Up()
	// if err != nil && err != migrate.ErrNoChange {
	// 	logger.Fatal(ctx, err, "❌️ Could Not Migrate MongoDB")
	// }

	return mongoDB
}

func (c *mongoConnection) Shutdown(ctx context.Context) error {
	if err := close(ctx, c.Master); err != nil {
		return err
	}

	if err := close(ctx, c.Slave); err != nil {
		return err
	}

	return nil
}

func close(ctx context.Context, db *mongo.Database) error {
	return db.Client().Disconnect(ctx)
}

func (c *mongoConnection) Ping(ctx context.Context) error {
	if err := c.Master.Client().Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	if err := c.Slave.Client().Ping(ctx, readpref.Secondary()); err != nil {
		return err
	}

	return nil
}
