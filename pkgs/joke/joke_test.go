package joke

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
	"time"
)

var (
	testJokeSearch = SearchParameters{
		Type:     Single,
		Category: Any,
	}
	testJoke = Joke{
		Question: "testQuestion",
		Answer:   "testAnswer",
		Type:     Single,
		Category: Any,
		GuildID:  "",
	}
)

type testMongodbService struct {
	client *mongo.Client
	ctx    context.Context
}

func (t testMongodbService) Close() error {
	return t.client.Disconnect(t.ctx)
}

func (t testMongodbService) Client(_ context.Context) (*mongo.Client, error) {
	return t.client, nil
}

func createMTest(t *testing.T) *mtest.T {
	return mtest.New(t, mtest.NewOptions().
		ClientType(mtest.Mock).
		CollectionName(collectionName).
		DatabaseName(DatabaseName))
}

func TestJokeService_Add(t *testing.T) {
	mt := createMTest(t)

	mt.Run("add new joke", func(t *mtest.T) {
		t.AddMockResponses(mtest.CreateSuccessResponse(bson.E{"ok", "1"},
			bson.E{"_id", primitive.NewObjectID()}))

		ctx := context.Background()

		mongodbService := testMongodbService{
			t.Client,
			ctx,
		}

		service := DatabaseJokeService{&mongodbService}

		if _, err := service.Add(ctx, testJoke); err != nil {
			mt.Fatal(err)
		}
	})
}

func TestJokeService_AddButContextCancelled(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now())
	cancel()

	mongodbService := testMongodbService{
		nil,
		ctx,
	}

	service := DatabaseJokeService{&mongodbService}

	if _, err := service.Add(ctx, testJoke); err == nil {
		t.Fatal("Context wasn't cancelled")
	}
}

func TestJokeService_SearchButContextWasCancelled(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now())
	cancel()

	mongodbService := testMongodbService{
		nil,
		ctx,
	}

	service := DatabaseJokeService{&mongodbService}

	if _, err := service.Get(ctx, testJokeSearch); err == nil {
		t.Fatal("Context wasn't cancelled")
	}
}

func TestJokeService_SearchButNotingFound(t *testing.T) {
	mt := createMTest(t)

	mt.Run("get new joke, but nothing was found", func(t *mtest.T) {
		t.AddMockResponses(mtest.CreateCursorResponse(1, DatabaseName+"."+collectionName, mtest.FirstBatch, bson.D{}))

		ctx := context.Background()

		mongodbService := testMongodbService{
			t.Client,
			ctx,
		}

		service := DatabaseJokeService{&mongodbService}

		if _, err := service.Get(ctx, testJokeSearch); err == nil {
			mt.Fatal("Something was found in database, but it shouldn't")
		}
	})
}

func TestJokeService_SearchButFindRandomJoke(t *testing.T) {
	mt := createMTest(t)

	mt.Run("get new joke, but nothing was found", func(t *mtest.T) {
		t.AddMockResponses(mtest.CreateSuccessResponse(bson.E{"ok", "1"},
			bson.E{"_id", primitive.NewObjectID()}))
		t.AddMockResponses(mtest.CreateCursorResponse(0, collectionName+".0", mtest.FirstBatch, bson.D{
			{"_id", primitive.NewObjectID()},
			{"question", testJoke.Question},
			{"answer", testJoke.Answer},
			{"type", testJoke.Type},
			{"category", testJoke.Category},
			{"guild_id", testJoke.GuildID},
		}))

		ctx := context.Background()

		mongodbService := testMongodbService{
			t.Client,
			ctx,
		}

		service := DatabaseJokeService{&mongodbService}

		joke, err := service.Get(ctx, SearchParameters{})
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

func TestDatabaseJokeService_ActiveButContextCancelled(t *testing.T) {
	mt := createMTest(t)
	mt.Run("testing active database", func(t *mtest.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		mongodbService := testMongodbService{
			t.Client,
			ctx,
		}

		service := DatabaseJokeService{&mongodbService}

		if service.Active(ctx) {
			t.Fatal("service can still running and handle new requests")
		}
	})
}

func TestDatabaseJokeService_Active(t *testing.T) {
	mt := createMTest(t)
	mt.Run("testing active database", func(t *mtest.T) {
		t.AddMockResponses(mtest.CreateSuccessResponse(bson.E{"ok", "1"},
			bson.E{"_id", primitive.NewObjectID()}))

		ctx := context.Background()

		mongodbService := testMongodbService{
			t.Client,
			ctx,
		}

		service := DatabaseJokeService{&mongodbService}

		if !service.Active(ctx) {
			t.Fatal("service isn't responding")
		}
	})
}
