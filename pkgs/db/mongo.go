package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

type MongoDBDatabase struct {
	uri string
	ctx context.Context
	db  *mongo.Client
}

func (m *MongoDBDatabase) Close() (err error) {
	err = m.db.Disconnect(m.ctx)
	m.db = nil

	return err
}

func (m *MongoDBDatabase) Client(ctx context.Context) (*mongo.Client, error) {
	if m.uri == "" {
		return nil, errors.New("missing URI to connect database")
	}

	if m.db != nil {
		return m.db, nil
	}

	var err error

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	m.db, err = mongo.Connect(ctx, options.Client().
		ApplyURI(m.uri).
		SetServerAPIOptions(serverAPI))
	if err != nil {
		m.db = nil
		return nil, err
	}

	return m.db, nil
}

func NewMongoDatabase(ctx context.Context) (db *MongoDBDatabase) {
	db = new(MongoDBDatabase)

	if uri, ok := os.LookupEnv("MONGODB_URI"); ok {
		db.uri = uri
	}

	db.ctx = ctx

	return
}
