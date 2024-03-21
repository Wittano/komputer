package db

import (
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"io"
)

type MongodbService interface {
	io.Closer
	Client(ctx context.Context) (*mongo.Client, error)
}
