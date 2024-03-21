package db

import (
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"io"
)

type MongoDBService interface {
	io.Closer
	Client(ctx context.Context) (*mongo.Client, error)
}
