//go:build testcontainers

package db

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/wittano/komputer/pkgs/joke"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"testing"
)

var testJoke = joke.Joke{
	Question: "testQuestion",
	Answer:   "testAnswer",
	Type:     joke.Single,
	Category: joke.Any,
	GuildID:  "",
}

type jokeSerchArgs struct {
	name   string
	withID bool
	search joke.SearchParameters
	joke   joke.Joke
}

func createJokeService(ctx context.Context) (db *mongodb.MongoDBContainer, service joke.DatabaseJokeService, err error) {
	db, err = mongodb.RunContainer(ctx, testcontainers.WithEnv(map[string]string{
		"MONGO_INITDB_DATABASE": joke.DatabaseName,
	}))
	if err != nil {
		return
	}

	mongodbURI, err := db.ConnectionString(ctx)
	if err != nil {
		return
	}
	os.Setenv(mongodbURIKey, mongodbURI)

	service = joke.NewDatabaseJokeService(NewMongodbDatabase(ctx))

	return
}

func terminateDatabase(ctx context.Context, t *testing.T, db *mongodb.MongoDBContainer) {
	if err := db.Terminate(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestJokeService_AddWithRealDatabase(t *testing.T) {
	ctx := context.Background()

	db, service, err := createJokeService(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer terminateDatabase(ctx, t, db)

	id, err := service.Add(ctx, testJoke)
	if err != nil {
		t.Fatal(err)
	}

	if id == "" {
		t.Fatal("New joke wasn't saved")
	}
}

func TestJokeService_Get(t *testing.T) {
	const testGuildID = "1"

	args := []jokeSerchArgs{
		{
			name:   "find random joke",
			search: joke.SearchParameters{},
			joke:   testJoke,
		},
		{
			name:   "find by type",
			search: joke.SearchParameters{Type: joke.TwoPart},
			joke: joke.Joke{
				Type:     joke.TwoPart,
				Question: "testTwoPartQuestion",
				Answer:   "testTwoPartAnswer",
				Category: joke.DARK,
				GuildID:  testGuildID,
			},
		},
		{
			name:   "find by category",
			search: joke.SearchParameters{Category: joke.MISC},
			joke: joke.Joke{
				Type:     joke.Single,
				Question: "testQuestion",
				Answer:   "testAnswer",
				Category: joke.MISC,
				GuildID:  testGuildID,
			},
		},
		{
			name:   "find by id",
			withID: true,
			search: joke.SearchParameters{}, // I added ID after adding entity to database
			joke: joke.Joke{
				Type:     joke.Single,
				Question: "testQuestion",
				Answer:   "testAnswer",
				Category: joke.PROGRAMMING,
				GuildID:  testGuildID,
			},
		},
		{
			name:   "find by id and category",
			withID: true,
			search: joke.SearchParameters{Category: joke.MISC},
			joke: joke.Joke{
				Type:     joke.Single,
				Question: "testQuestion",
				Answer:   "testAnswer",
				Category: joke.MISC,
				GuildID:  testGuildID,
			},
		},
		{
			name:   "find by id and type",
			withID: true,
			search: joke.SearchParameters{Type: joke.TwoPart},
			joke: joke.Joke{
				Type:     joke.TwoPart,
				Question: "testTwoPartQuestion",
				Answer:   "testTwoPartAnswer",
				Category: joke.DARK,
				GuildID:  testGuildID,
			},
		},
		{
			name:   "find by type and category",
			search: joke.SearchParameters{Type: joke.TwoPart, Category: joke.YOMAMA},
			joke: joke.Joke{
				Type:     joke.TwoPart,
				Question: "testTwoPartQuestion",
				Answer:   "testTwoPartAnswer",
				Category: joke.YOMAMA,
				GuildID:  testGuildID,
			},
		},
		{
			name:   "find by id, category and type",
			withID: true,
			search: joke.SearchParameters{Type: joke.TwoPart, Category: joke.DARK},
			joke: joke.Joke{
				Type:     joke.TwoPart,
				Question: "testTwoPartQuestion",
				Answer:   "testTwoPartAnswer",
				Category: joke.DARK,
				GuildID:  testGuildID,
			},
		},
	}

	ctx := context.WithValue(context.Background(), joke.GuildIDKey, testGuildID)
	db, service, err := createJokeService(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer terminateDatabase(ctx, t, db)

	// Added jokes to database
	for _, data := range args {
		newJoke, search, withID := data.joke, data.search, data.withID
		id, err := service.Add(ctx, newJoke)
		if err != nil {
			t.Fatal(err)
		}
		newJoke.ID, err = primitive.ObjectIDFromHex(id)
		if err != nil {
			t.Fatal(err)
		}

		if withID {
			search.ID = newJoke.ID
		}
	}

	for _, data := range args {
		search := data.search
		t.Run(data.name, func(t *testing.T) {
			if _, err := service.Get(ctx, search); err != nil {
				t.Fatal(err)
			}
		})
	}
}
