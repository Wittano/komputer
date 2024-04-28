package audio

import (
	"context"
	"github.com/wittano/komputer/db"
	"github.com/wittano/komputer/test"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"mime/multipart"
	"testing"
	"time"
)

func TestUploadRequestedFile(t *testing.T) {
	if err := test.LoadDefaultConfig(t); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	filePath, err := test.CreateTempAudioFiles(t)
	if err != nil {
		t.Fatal(err)
	}

	header, err := test.CreateMultipartFileHeader(filePath)
	if err != nil {
		t.Fatal(err)
	}

	mt := mtest.New(t, mtest.NewOptions().
		ClientType(mtest.Mock).
		CollectionName("audio").
		DatabaseName(db.DatabaseName))
	mt.Run("upload requested file", func(t *mtest.T) {
		t.AddMockResponses(mtest.CreateSuccessResponse(bson.E{Key: "ok", Value: "1"},
			bson.E{Key: "_id", Value: primitive.NewObjectID()}))

		service := UploadService{Db: test.NewMockedMongodbService(ctx, t.Client)}

		if err := service.Upload(ctx, []*multipart.FileHeader{header}); err != nil {
			t.Fatal(err)
		}
	})
}
