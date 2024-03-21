package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
	"time"
)

var (
	testJokeSearch = JokeSearch{
		Type:     Single,
		Category: ANY,
	}
	testJoke = Joke{
		Question: "testQuestion",
		Answer:   "testAnswer",
		Type:     Single,
		Category: ANY,
		GuildID:  "",
	}
)

func createMTest(t *testing.T) *mtest.T {
	return mtest.New(t, mtest.NewOptions().
		ClientType(mtest.Mock).
		CollectionName(collectionName).
		DatabaseName(databaseName))
}

func TestJokeService_Add(t *testing.T) {
	mt := createMTest(t)

	mt.Run("add new joke", func(t *mtest.T) {
		t.AddMockResponses(mtest.CreateSuccessResponse(bson.E{"ok", "1"},
			bson.E{"_id", primitive.NewObjectID()}))

		ctx := context.Background()

		mongodbService := MongodbDatabase{
			uri,
			ctx,
			t.Client,
		}

		service := NewJokeService(&mongodbService)

		if _, err := service.Add(ctx, testJoke); err != nil {
			mt.Fatal(err)
		}
	})
}

func TestJokeService_AddButContextCancelled(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now())
	cancel()

	mongodbService := MongodbDatabase{
		uri,
		ctx,
		nil,
	}

	if _, err := NewJokeService(&mongodbService).Add(ctx, testJoke); err == nil {
		t.Fatal("Context wasn't cancelled")
	}
}

func TestJokeService_SearchButContextWasCancelled(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now())
	cancel()

	mongodbService := MongodbDatabase{
		uri,
		ctx,
		nil,
	}

	if _, err := NewJokeService(&mongodbService).Get(ctx, testJokeSearch); err == nil {
		t.Fatal("Context wasn't cancelled")
	}
}

func TestJokeService_SearchButNotingFound(t *testing.T) {
	mt := createMTest(t)

	mt.Run("get new joke, but nothing was found", func(t *mtest.T) {
		t.AddMockResponses(mtest.CreateCursorResponse(1, databaseName+"."+collectionName, mtest.FirstBatch, bson.D{}))

		ctx := context.Background()

		mongodbService := MongodbDatabase{
			uri,
			ctx,
			t.Client,
		}

		service := NewJokeService(&mongodbService)

		if _, err := service.Get(ctx, testJokeSearch); err == nil {
			mt.Fatal("Something was found in database, but it shouldn't")
		}
	})
}

func TestJokeService_SearchButFindRandomJoke(t *testing.T) {
	mt := createMTest(t)

	mt.Run("get new joke, but nothing was found", func(t *mtest.T) {
		t.AddMockResponses(mtest.CreateCursorResponse(0, collectionName+".0", mtest.FirstBatch, bson.D{
			{"_id", primitive.NewObjectID()},
			{"question", testJoke.Question},
			{"answer", testJoke.Answer},
			{"type", testJoke.Type},
			{"category", testJoke.Category},
			{"guild_id", testJoke.GuildID},
		}))

		ctx := context.Background()

		mongodbService := MongodbDatabase{
			uri,
			ctx,
			t.Client,
		}

		service := NewJokeService(&mongodbService)

		joke, err := service.Get(ctx, JokeSearch{})
		if err != nil {
			t.Fatal(err)
		}

		if joke.Category != testJoke.Category {
			t.Fatalf("Invalid Category. Expected '%s', Result: '%s'", testJoke.Category, joke.Category)
		}
		if joke.Type != testJoke.Type {
			t.Fatalf("Invalid Type. Expected '%s', Result: '%s'", testJoke.Type, joke.Type)
		}
		if joke.GuildID != testJoke.GuildID {
			t.Fatalf("Invalid GuildID. Expected '%s', Result: '%s'", testJoke.GuildID, joke.GuildID)
		}
		if joke.Question != testJoke.Question {
			t.Fatalf("Invalid Question. Expected '%s', Result: '%s'", testJoke.Question, joke.Question)
		}
		if joke.Answer != testJoke.Answer {
			t.Fatalf("Invalid Answer. Expected '%s', Result: '%s'", testJoke.Answer, joke.Answer)
		}
	})
}
