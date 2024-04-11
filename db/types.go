package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
)

// TODO Added option to custom name of database
const DatabaseName = "komputer"

type AudioInfo struct {
	ID   primitive.ObjectID `bson:"_id"`
	Path string             `bson:"path"`
}

type MongodbService interface {
	io.Closer
	Client(ctx context.Context) (*mongo.Client, error)
}
