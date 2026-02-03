package mongodb

import (
	"context"
	"fmt"
	"net/url"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/goodone-dev/go-boilerplate/internal/config"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/logger"
	"github.com/goodone-dev/go-boilerplate/internal/utils/retry"
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
		Host:   fmt.Sprintf("%s:%d", config.Mongo.Host, config.Mongo.Port),
		User:   url.UserPassword(config.Mongo.Username, config.Mongo.Password),
	}

	if len(config.Mongo.MasterHost) > 0 {
		uriMaster = url.URL{
			Scheme: "mongodb",
			Host:   fmt.Sprintf("%s:%d", config.Mongo.MasterHost, config.Mongo.MasterPort),
			User:   url.UserPassword(config.Mongo.MasterUsername, config.Mongo.MasterPassword),
		}
	}

	uriSlave := url.URL{
		Scheme: "mongodb",
		Host:   fmt.Sprintf("%s:%d", config.Mongo.Host, config.Mongo.Port),
		User:   url.UserPassword(config.Mongo.Username, config.Mongo.Password),
	}

	if len(config.Mongo.SlaveHost) > 0 {
		uriSlave = url.URL{
			Scheme: "mongodb",
			Host:   fmt.Sprintf("%s:%d", config.Mongo.SlaveHost, config.Mongo.SlavePort),
			User:   url.UserPassword(config.Mongo.SlaveUsername, config.Mongo.SlavePassword),
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

	conn := &mongoConnection{
		Master: open(ctx, mongoConfig.Master, readpref.Primary()),
		Slave:  open(ctx, mongoConfig.Slave, readpref.Secondary()),
	}

	go conn.Monitor(ctx)

	return conn
}

func open(ctx context.Context, opts *options.ClientOptions, rp *readpref.ReadPref) *mongo.Database {
	// TODO: Enable MongoDB OpenTelemetry monitoring once otelmongo supports mongo-driver v2
	// Currently blocked by: https://github.com/open-telemetry/opentelemetry-go-contrib/issues/
	// The otelmongo package only supports mongo-driver v1.x
	// opts.SetMonitor(otelmongo.NewMonitor())
	opts.SetDirect(true)
	opts.SetRetryWrites(false)
	opts.SetMaxConnIdleTime(time.Duration(config.Mongo.ConnIdleTimeoutMS) * time.Millisecond)
	opts.SetBSONOptions(&options.BSONOptions{
		UseLocalTimeZone: true,
	})
	if config.Mongo.MaxConnPoolSize >= 0 {
		opts.SetMaxPoolSize(uint64(config.Mongo.MaxConnPoolSize))
	}
	if config.Mongo.MinConnPoolSize >= 0 {
		opts.SetMinPoolSize(uint64(config.Mongo.MinConnPoolSize))
	}

	client, err := retry.RetryWithBackoff(ctx, "MongoDB connection", func() (*mongo.Client, error) {
		return mongo.Connect(opts)
	})
	if err != nil {
		logger.Fatal(ctx, err, "❌ MongoDB failed to establish connection after retries").Write()
	}

	_, err = retry.RetryWithBackoff(ctx, "MongoDB connection test", func() (any, error) {
		return nil, client.Ping(ctx, rp)
	})
	if err != nil {
		logger.Fatal(ctx, err, "❌ MongoDB connection test failed").Write()
	}

	mongoDB := client.Database(config.Mongo.Database)
	if !config.Mongo.AutoMigrate {
		return mongoDB
	}

	// TODO: Enable MongoDB migrations once golang-migrate supports mongo-driver v2
	// Currently blocked by: https://github.com/golang-migrate/migrate/pull/1265
	// The migrate mongodb package only supports mongo-driver v1.x
	// Alternative: Consider using github.com/xakep666/mongo-migrate which supports v2
	//
	// Example implementation:
	// migrateDriver, err := migratemongo.WithInstance(client, &migratemongo.Config{
	// 	DatabaseName: config.Mongo.Database,
	// })
	// if err != nil {
	// 	logger.Fatal(ctx, err, "❌ Failed to initialize MongoDB migration driver").Write()
	// }
	//
	// m, err := migrate.NewWithDatabaseInstance("file://migrations/mongodb", "mongodb", migrateDriver)
	// if err != nil {
	// 	logger.Fatal(ctx, err, "❌ Failed to create migration instance from MongoDB driver").Write()
	// }
	//
	// err = m.Up()
	// if err != nil && err != migrate.ErrNoChange {
	// 	logger.Fatal(ctx, err, "❌ MongoDB migration failed").Write()
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

func (c *mongoConnection) Monitor(ctx context.Context) {
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
					logger.Errorf(ctx, err, "🛑 MongoDB connection lost").Write()
					wasLost = true
				}
			} else {
				if wasLost {
					logger.Info(ctx, "✅ MongoDB connection restored").Write()
					wasLost = false
				}
			}
		}
	}
}
