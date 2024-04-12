package test

import (
	"context"
	"github.com/wittano/komputer/db"
	"go.mongodb.org/mongo-driver/mongo"
)

// TestMongodbService This service is mock. It shouldn't use in production code
type TestMongodbService struct {
	client *mongo.Client
	ctx    context.Context
}

func (t TestMongodbService) Close() error {
	return t.client.Disconnect(t.ctx)
}

func (t TestMongodbService) Client(_ context.Context) (*mongo.Client, error) {
	return t.client, nil
}

func NewMockedMognodbService(ctx context.Context, client *mongo.Client) db.MongodbService {
	return &TestMongodbService{client, ctx}
}
