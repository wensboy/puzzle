package database

import (
	"context"

	"github.com/wendisx/puzzle/pkg/clog"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	_default_mongo_dsn = "mongodb://<username>:<password>@<host:port>/?<options>"
)

type (
	// MongoDB is just a layer of mongo encapsulation.
	MongoDB struct {
		client *mongo.Client
	}
)

// InitMongo return new instance of mongodb.
func InitMongo(dsn string) MongoDB {
	if dsn == "" {
		dsn = _default_mongo_dsn
	}
	client, err := mongo.Connect(options.Client().ApplyURI(dsn))
	if err != nil {
		clog.Panic(err.Error())
	}
	return MongoDB{
		client: client,
	}
}

// In most scenarios, the client needs to be closed in defer,
// which may be related to the location where initialization is called.
func (db *MongoDB) Close() {
	if err := db.client.Disconnect(context.TODO()); err != nil {
		clog.Panic(err.Error())
	}
}
