package mongodb

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/goodone-dev/go-boilerplate/internal/config"
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

func Open(ctx context.Context) mongoConnection {
	mongoConfig := setConfig()

	return mongoConnection{
		Master: open(ctx, mongoConfig.Master),
		Slave:  open(ctx, mongoConfig.Slave),
	}
}

func open(ctx context.Context, opts *options.ClientOptions) *mongo.Database {
	// FIXME: mongodb monitor
	// opts.SetMonitor(otelmongo.NewMonitor())
	opts.SetDirect(true)
	opts.SetRetryWrites(false)
	opts.SetMaxPoolSize(uint64(config.MongoConfig.MaxConnPoolSize))
	opts.SetMinPoolSize(uint64(config.MongoConfig.MinConnPoolSize))
	opts.SetMaxConnIdleTime(time.Duration(config.MongoConfig.ConnIdleTimeoutMS) * time.Millisecond)
	opts.SetBSONOptions(&options.BSONOptions{
		UseLocalTimeZone: true,
	})

	client, err := mongo.Connect(opts)
	if err != nil {
		log.Fatalf("❌ Could not to connect MongoDB connection: %v", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("❌ Could not to ping MongoDB database: %v", err)
	}

	mongoDB := client.Database(config.MongoConfig.Database)
	if !config.MongoConfig.AutoMigrate {
		return mongoDB
	}

	// FIXME: mongodb migration
	// migrateDriver, err := migratemongo.WithInstance(client, &migratemongo.Config{})
	// if err != nil {
	// 	log.Fatalf("❌ Could not to create migrate instance for MongoDB:%v", err)
	// }

	// m, err := migrate.NewWithDatabaseInstance("file://migrations/mongodb", "mongodb", migrateDriver)
	// if err != nil {
	// 	log.Fatalf("❌ Could not to create migrate instance for MongoDB:%v", err)
	// }

	// err = m.Up()
	// if err != nil && err != migrate.ErrNoChange {
	// 	log.Fatalf("❌ Could not to migrate MongoDB:%v", err)
	// }

	return mongoDB
}

func (c mongoConnection) Shutdown(ctx context.Context) error {
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
