package db

import (
	"context"
	"github.com/wittano/komputer/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"strings"
)

const DatabaseName = "komputer"

type AudioInfo struct {
	ID       primitive.ObjectID `bson:"_id"`
	Path     string             `bson:"path"`
	Original string             `bson:"original_name"`
}

func (a AudioInfo) AudioFileInfo() api.AudioFileInfo {
	return api.AudioFileInfo{
		ID:       a.ID.Hex(),
		Filename: strings.TrimSuffix(a.Original, ".mp3"),
	}
}

type MongodbService interface {
	io.Closer
	Client(ctx context.Context) (*mongo.Client, error)
}
