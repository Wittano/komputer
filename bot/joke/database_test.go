package joke

import (
	"context"
	"github.com/wittano/komputer/db"
	"github.com/wittano/komputer/test"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
	"time"
)

var (
	testParams = SearchParams{
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

func createMTest(t *testing.T) *mtest.T {
	return mtest.New(t, mtest.NewOptions().
		ClientType(mtest.Mock).
		CollectionName(collectionName).
		DatabaseName(db.DatabaseName))
}

func TestDatabaseService_Add(t *testing.T) {
	mt := createMTest(t)

	mt.Run("add new joke", func(t *mtest.T) {
		t.AddMockResponses(mtest.CreateSuccessResponse(bson.E{Key: "ok", Value: "1"},
			bson.E{Key: "_id", Value: primitive.NewObjectID()}))

		ctx := context.Background()

		service := DatabaseService{test.NewMockedMongodbService(ctx, t.Client)}

		if _, err := service.Add(ctx, testJoke); err != nil {
			mt.Fatal(err)
		}
	})
}

func TestDatabaseService_AddWithContextCancelled(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now())
	cancel()

	service := DatabaseService{test.NewMockedMongodbService(ctx, nil)}

	if _, err := service.Add(ctx, testJoke); err == nil {
		t.Fatal("Context wasn't cancelled")
	}
}

func TestDatabaseService_JokeWithContextCancelled(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now())
	cancel()

	service := DatabaseService{test.NewMockedMongodbService(ctx, nil)}

	if _, err := service.Joke(ctx, testParams); err == nil {
		t.Fatal("Context wasn't cancelled")
	}
}

func TestDatabaseService_JokeReturnEmptyWithoutError(t *testing.T) {
	mt := createMTest(t)

	mt.Run("get new joke, but nothing was found", func(t *mtest.T) {
		t.AddMockResponses(mtest.CreateCursorResponse(1, db.DatabaseName+"."+collectionName, mtest.FirstBatch, bson.D{}))

		ctx := context.Background()

		service := DatabaseService{test.NewMockedMongodbService(ctx, t.Client)}

		if _, err := service.Joke(ctx, testParams); err == nil {
			mt.Fatal("Something was found in database, but it shouldn't")
		}
	})
}

func TestDatabaseService_FindRandomJoke(t *testing.T) {
	mt := createMTest(t)

	mt.Run("get new joke, but nothing was found", func(t *mtest.T) {
		t.AddMockResponses(mtest.CreateSuccessResponse(bson.E{Key: "ok", Value: "1"},
			bson.E{Key: "_id", Value: primitive.NewObjectID()}))
		t.AddMockResponses(mtest.CreateCursorResponse(0, collectionName+".0", mtest.FirstBatch, bson.D{
			{"_id", primitive.NewObjectID()},
			{"question", testJoke.Question},
			{"answer", testJoke.Answer},
			{"type", testJoke.Type},
			{"category", testJoke.Category},
			{"guild_id", testJoke.GuildID},
		}))

		ctx := context.Background()

		service := DatabaseService{test.NewMockedMongodbService(ctx, t.Client)}

		joke, err := service.Joke(ctx, SearchParams{})
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

func TestDatabaseService_ActiveWithContextCancelled(t *testing.T) {
	mt := createMTest(t)
	mt.Run("testing active database", func(t *mtest.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		service := DatabaseService{test.NewMockedMongodbService(ctx, t.Client)}

		if service.Active(ctx) {
			t.Fatal("service can still running and handle new requests")
		}
	})
}

func TestDatabaseService_Active(t *testing.T) {
	mt := createMTest(t)
	mt.Run("testing active database", func(t *mtest.T) {
		t.AddMockResponses(mtest.CreateSuccessResponse(bson.E{Key: "ok", Value: "1"},
			bson.E{Key: "_id", Value: primitive.NewObjectID()}))

		ctx := context.Background()

		service := DatabaseService{test.NewMockedMongodbService(ctx, t.Client)}

		if !service.Active(ctx) {
			t.Fatal("service isn't responding")
		}
	})
}

func TestUnlockService(t *testing.T) {
	testFlag := false
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	resetTime := time.Now()

	unlockService(ctx, &testFlag, resetTime)

	if testFlag != true {
		t.Fatal("Service doesn't unlock")
	}
}

func TestUnlockService_ParentContextCancelled(t *testing.T) {
	testFlag := false
	ctx, cancel := context.WithCancel(context.Background())
	resetTime := time.Now().Add(1 * time.Hour)

	cancel()
	unlockService(ctx, &testFlag, resetTime)

	if testFlag != true {
		t.Fatal("Service doesn't unlock")
	}
}
