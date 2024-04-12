package audio

import (
	"context"
	"errors"
	"github.com/wittano/komputer/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"path/filepath"
	"strings"
)

const audioCollectionName = "audio"

type DatabaseService struct {
	Database db.MongodbService
}

var NotFoundErr = errors.New("audio not found")

func (a DatabaseService) save(ctx context.Context, filename string) (primitive.ObjectID, error) {
	client, err := a.Database.Client(ctx)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	base := filepath.Base(filename)
	id := primitive.NewObjectID()
	realPath := strings.ReplaceAll(filename, base, id.Hex()+".mp3")
	info := db.AudioInfo{
		ID:       id,
		Original: base,
		Path:     realPath,
	}

	result, err := client.Database(db.DatabaseName).Collection(audioCollectionName).InsertOne(ctx, info)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (a DatabaseService) Get(ctx context.Context, id string) (result db.AudioInfo, err error) {
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}

	client, err := a.Database.Client(ctx)
	if err != nil {
		return db.AudioInfo{}, err
	}

	err = client.Database(db.DatabaseName).
		Collection(audioCollectionName).
		FindOne(ctx, bson.D{{"_id", hex}}).
		Decode(&result)

	return
}

func (a DatabaseService) Delete(ctx context.Context, id string) (err error) {
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}

	info, err := a.Get(ctx, id)
	if err != nil {
		return errors.Join(NotFoundErr, err)
	}

	client, err := a.Database.Client(ctx)
	if err != nil {
		return
	}

	_, err = client.Database(db.DatabaseName).
		Collection(audioCollectionName).
		DeleteOne(ctx, bson.D{{"_id", hex}})
	if err != nil {
		return
	}

	if _, err = os.Stat(info.Path); errors.Is(err, os.ErrNotExist) {
		return os.Remove(info.Path)
	}

	return
}
