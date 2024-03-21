//go:build testcontainers

package db

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"os"
	"testing"
)

type jokeSerchArgs struct {
	name   string
	withID bool
	search JokeSearch
	joke   Joke
}

func createJokeService(ctx context.Context) (db *mongodb.MongoDBContainer, service JokeService, err error) {
	db, err = mongodb.RunContainer(ctx, testcontainers.WithEnv(map[string]string{
		"MONGO_INITDB_DATABASE": databaseName,
	}))
	if err != nil {
		return
	}

	mongodbURI, err := db.ConnectionString(ctx)
	if err != nil {
		return
	}
	os.Setenv(mongodbURIKey, mongodbURI)

	service = NewJokeService(NewMongodbDatabase(ctx))

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

	if id == [12]byte{} {
		t.Fatal("New joke wasn't saved")
	}
}

func TestJokeService_Get(t *testing.T) {
	const testGuildID = "1"

	args := []jokeSerchArgs{
		{
			name:   "find random joke",
			search: JokeSearch{},
			joke:   testJoke,
		},
		{
			name:   "find by type",
			search: JokeSearch{Type: TwoPart},
			joke: Joke{
				Type:     TwoPart,
				Question: "testTwoPartQuestion",
				Answer:   "testTwoPartAnswer",
				Category: DARK,
				GuildID:  testGuildID,
			},
		},
		{
			name:   "find by category",
			search: JokeSearch{Category: MISC},
			joke: Joke{
				Type:     Single,
				Question: "testQuestion",
				Answer:   "testAnswer",
				Category: MISC,
				GuildID:  testGuildID,
			},
		},
		{
			name:   "find by id",
			withID: true,
			search: JokeSearch{}, // I added ID after adding entity to database
			joke: Joke{
				Type:     Single,
				Question: "testQuestion",
				Answer:   "testAnswer",
				Category: PROGRAMMING,
				GuildID:  testGuildID,
			},
		},
		{
			name:   "find by id and category",
			withID: true,
			search: JokeSearch{Category: MISC},
			joke: Joke{
				Type:     Single,
				Question: "testQuestion",
				Answer:   "testAnswer",
				Category: MISC,
				GuildID:  testGuildID,
			},
		},
		{
			name:   "find by id and type",
			withID: true,
			search: JokeSearch{Type: TwoPart},
			joke: Joke{
				Type:     TwoPart,
				Question: "testTwoPartQuestion",
				Answer:   "testTwoPartAnswer",
				Category: DARK,
				GuildID:  testGuildID,
			},
		},
		{
			name:   "find by type and category",
			search: JokeSearch{Type: TwoPart, Category: YOMAMA},
			joke: Joke{
				Type:     TwoPart,
				Question: "testTwoPartQuestion",
				Answer:   "testTwoPartAnswer",
				Category: YOMAMA,
				GuildID:  testGuildID,
			},
		},
		{
			name:   "find by id, category and type",
			withID: true,
			search: JokeSearch{Type: TwoPart, Category: DARK},
			joke: Joke{
				Type:     TwoPart,
				Question: "testTwoPartQuestion",
				Answer:   "testTwoPartAnswer",
				Category: DARK,
				GuildID:  testGuildID,
			},
		},
	}

	ctx := context.WithValue(context.Background(), GuildIDKey, testGuildID)
	db, service, err := createJokeService(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer terminateDatabase(ctx, t, db)

	// Added jokes to database
	for _, data := range args {
		joke, search, withID := data.joke, data.search, data.withID
		id, err := service.Add(ctx, joke)
		if err != nil {
			t.Fatal(err)
		}
		joke.id = id
		if withID {
			search.ID = id
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
