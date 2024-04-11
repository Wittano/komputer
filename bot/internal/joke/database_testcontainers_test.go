//go:build testcontainers

package joke

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/wittano/komputer/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"testing"
)

type jokeSerchArgs struct {
	name   string
	withID bool
	search SearchParameters
	joke   Joke
}

func createJokeService(ctx context.Context) (container *mongodb.MongoDBContainer, service DatabaseJokeService, err error) {
	container, err = mongodb.RunContainer(ctx, testcontainers.WithEnv(map[string]string{
		"MONGO_INITDB_DATABASE": db.DatabaseName,
	}))
	if err != nil {
		return
	}

	mongodbURI, err := container.ConnectionString(ctx)
	if err != nil {
		return
	}
	os.Setenv(db.MongodbURIKey, mongodbURI)

	service = NewDatabaseJokeService(db.Mongodb(ctx))

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
			search: SearchParameters{},
			joke:   testJoke,
		},
		{
			name:   "find by type",
			search: SearchParameters{Type: TwoPart},
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
			search: SearchParameters{Category: MISC},
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
			search: SearchParameters{}, // I added ID after adding entity to database
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
			search: SearchParameters{Category: MISC},
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
			search: SearchParameters{Type: TwoPart},
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
			search: SearchParameters{Type: TwoPart, Category: YOMAMA},
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
			search: SearchParameters{Type: TwoPart, Category: DARK},
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
