package db

import (
	"context"
	"github.com/wittano/komputer/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"strings"
)

// TODO Added option to custom NameType of database
const DatabaseName = "komputer"

// TODO Added custom title for audio
type AudioInfo struct {
	ID       primitive.ObjectID `bson:"_id"`
	Path     string             `bson:"path"`
	Original string             `bson:"original_name"`
}

func (a AudioInfo) ApiAudioFileInfo() api.AudioFileInfo {
	return api.AudioFileInfo{
		ID:       a.ID.Hex(),
		Filename: strings.TrimSuffix(a.Original, ".mp3"),
	}
}

type MongodbService interface {
	io.Closer
	Client(ctx context.Context) (*mongo.Client, error)
}
