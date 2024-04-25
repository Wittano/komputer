package audio

import (
	"context"
	"errors"
	"fmt"
	"github.com/wittano/komputer/api"
	"github.com/wittano/komputer/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"path/filepath"
	"strings"
)

const audioCollectionName = "audio"

// TODO Change service into singleton
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

func (a DatabaseService) AudioFilesInfo(ctx context.Context, searchType, value string, page int) (fileInfos []api.AudioFileInfo, err error) {
	client, err := a.Database.Client(ctx)
	if err != nil {
		return []api.AudioFileInfo{}, err
	}

	keyName := "_id"
	if searchType == "name" {
		keyName = "original_name"
	}

	const maxPageSize = 10

	filter := bson.D{{keyName, primitive.Regex{Pattern: fmt.Sprintf("^[\\w]*%s[\\w\\.]+", value)}}}
	opts := options.Find().SetLimit(maxPageSize).SetSkip(int64(maxPageSize * page))
	cursor, err := client.Database(db.DatabaseName).Collection(audioCollectionName).Find(ctx, filter, opts)
	if err != nil {
		return []api.AudioFileInfo{}, err
	}
	defer cursor.Close(ctx)

	fileInfos = make([]api.AudioFileInfo, 0, maxPageSize)
	for cursor.TryNext(ctx) {
		var info db.AudioInfo
		if err := bson.Unmarshal(cursor.Current, &info); err != nil {
			return nil, err
		}

		fileInfos = append(fileInfos, info.ApiAudioFileInfo())
		if len(fileInfos) == maxPageSize {
			break
		}
	}

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
