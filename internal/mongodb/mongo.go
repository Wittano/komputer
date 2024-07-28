package mongodb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

const MongodbURIKey = "MONGODB_URI"

type Database struct {
	uri string
	ctx context.Context
	db  *mongo.Client
}

func (m *Database) Close() (err error) {
	if m.db == nil {
		return nil
	}

	err = m.db.Disconnect(m.ctx)
	m.db = nil

	return
}

func (m *Database) Client(ctx context.Context) (*mongo.Client, error) {
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

var db *Database

func NewMongodb(ctx context.Context) *Database {
	if db != nil {
		return db
	}

	db = new(Database)
	db.uri, _ = os.LookupEnv(MongodbURIKey)
	db.ctx = ctx

	return db
}
