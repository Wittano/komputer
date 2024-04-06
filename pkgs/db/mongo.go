package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

const mongodbURIKey = "MONGODB_URI"

type MongodbDatabase struct {
	uri string
	ctx context.Context
	db  *mongo.Client
}

func (m *MongodbDatabase) Close() (err error) {
	if m.db == nil {
		return nil
	}

	err = m.db.Disconnect(m.ctx)
	m.db = nil

	return err
}

func (m *MongodbDatabase) Client(ctx context.Context) (*mongo.Client, error) {
	select {
	case <-ctx.Done():
		return nil, context.Canceled
	default:
	}

	if m.uri == "" {
		return nil, errors.New("missing URI to connect database")
	}

	if m.db != nil {
		return m.db, nil
	}

	var err error

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	m.db, err = mongo.Connect(timeoutCtx, options.Client().
		ApplyURI(m.uri).
		SetServerAPIOptions(serverAPI))
	if err != nil {
		m.db = nil
		return nil, err
	}

	return m.db, nil
}

func NewMongodbDatabase(ctx context.Context) (db *MongodbDatabase) {
	db = new(MongodbDatabase)

	if uri, ok := os.LookupEnv(mongodbURIKey); ok {
		db.uri = uri
	}

	db.ctx = ctx

	return
}
