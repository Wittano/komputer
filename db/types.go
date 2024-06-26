package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
)

const DatabaseName = "komputer"

type MongodbService interface {
	io.Closer
	Client(ctx context.Context) (*mongo.Client, error)
}
