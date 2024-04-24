package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
)

// TODO Added option to custom NameType of database
const DatabaseName = "komputer"

// TODO Added custom title for audio
type AudioInfo struct {
	ID       primitive.ObjectID `bson:"_id"`
	Path     string             `bson:"path"`
	Original string             `bson:"original_name"`
}

type MongodbService interface {
	io.Closer
	Client(ctx context.Context) (*mongo.Client, error)
}
