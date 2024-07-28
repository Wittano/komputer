package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
)

const DatabaseName = "komputer"

type DatabaseGetter interface {
	io.Closer
	Client(ctx context.Context) (*mongo.Client, error)
}
