package audio

import (
	"context"
	"github.com/wittano/komputer/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const audioCollectionName = "audio"

func saveFileDataInDatabase(ctx context.Context, filename string) error {
	client, err := db.Mongodb(ctx).Client(ctx)
	if err != nil {
		return err
	}

	info := db.AudioInfo{
		ID:   primitive.NewObjectID(),
		Path: filename,
	}

	_, err = client.Database(db.DatabaseName).Collection(audioCollectionName).InsertOne(ctx, info)
	return err
}

func GetAudioInfo(ctx context.Context, id string) (result db.AudioInfo, err error) {
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}

	client, err := db.Mongodb(ctx).Client(ctx)
	if err != nil {
		return db.AudioInfo{}, err
	}

	err = client.Database(db.DatabaseName).
		Collection(audioCollectionName).
		FindOne(ctx, bson.D{{"_id", hex}}).
		Decode(&result)

	return
}
