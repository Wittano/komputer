package test

import (
	"context"
	"github.com/wittano/komputer/db"
	"go.mongodb.org/mongo-driver/mongo"
)

// TestMongodbService This service is mock. It shouldn't use in production code
type MongodbService struct {
	client *mongo.Client
	ctx    context.Context
}

func (t MongodbService) Close() error {
	return t.client.Disconnect(t.ctx)
}

func (t MongodbService) Client(_ context.Context) (*mongo.Client, error) {
	return t.client, nil
}

func NewMockedMongodbService(ctx context.Context, client *mongo.Client) db.MongodbService {
	return &MongodbService{client, ctx}
}
