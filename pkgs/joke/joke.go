package joke

import (
	"context"
	"fmt"
	"github.com/wittano/komputer/internal/file"
	"github.com/wittano/komputer/pkgs/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"math/rand"
	"time"
)

const (
	DatabaseName   = "komputer"
	collectionName = "jokes"
	GuildIDKey     = "guildID"
)

type (
	Type     string
	Category string
)

const (
	Single  Type = "single"
	TwoPart Type = "twopart"
)

const (
	PROGRAMMING Category = "Programming"
	MISC        Category = "Misc"
	DARK        Category = "Dark"
	YOMAMA      Category = "YoMama"
	Any         Category = "Any"
)

type Joke struct {
	ID       primitive.ObjectID `bson:"_id"`
	Question string             `bson:"question"`
	Answer   string             `bson:"answer"`
	Type     Type               `bson:"type"`
	Category Category           `bson:"category"`
	GuildID  string             `bson:"guild_id"`
}

type SearchParameters struct {
	Type     Type
	Category Category
	ID       primitive.ObjectID
}

type DatabaseJokeService struct {
	mongodb db.MongodbService
}

func (d DatabaseJokeService) Active(ctx context.Context) bool {
	client, err := d.mongodb.Client(ctx)
	if err != nil {
		return false
	}

	err = client.Ping(ctx, readpref.Nearest(readpref.WithMaxStaleness(500)))
	if err != nil {
		return false
	}

	return true
}

func (d DatabaseJokeService) Add(ctx context.Context, joke Joke) (string, error) {
	select {
	case <-ctx.Done():
		return "", context.Canceled
	default:
	}

	mongodb, err := d.mongodb.Client(ctx)
	if err != nil {
		return "", err
	}

	res, err := mongodb.Database(DatabaseName).Collection(collectionName).InsertOne(ctx, joke)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (d DatabaseJokeService) Get(ctx context.Context, search SearchParameters) (Joke, error) {
	select {
	case <-ctx.Done():
		return Joke{}, context.Canceled
	default:
	}

	mongodb, err := d.mongodb.Client(ctx)
	if err != nil {
		return Joke{}, err
	}

	if search.Category == "" {
		search.Category = Any
	}

	if search.Type == "" {
		search.Type = Single
	}

	// Create query to database
	const matchQueryKey = "$match"

	pipeline := mongo.Pipeline{{{
		"$sample", bson.D{{
			"size", 10,
		}},
	}},
		{{
			matchQueryKey, bson.D{
				{
					"type", search.Type,
				},
				{
					"guild_id", ctx.Value(GuildIDKey),
				},
			},
		}},
	}

	if search.ID != [12]byte{} {
		pipeline = append(pipeline, bson.D{{
			matchQueryKey, bson.D{{
				"_id", search.ID,
			}},
		}})
	}

	if search.Category != "" && search.Category != Any {
		pipeline = append(pipeline, bson.D{{
			matchQueryKey, bson.D{{
				"category", search.Category,
			}},
		}})
	}

	// SearchParameters
	res, err := mongodb.Database(DatabaseName).Collection(collectionName).Aggregate(ctx, pipeline)
	if err != nil {
		return Joke{}, err
	}
	defer res.Close(ctx)

	var jokes []Joke
	if err = res.All(ctx, &jokes); err != nil {
		return Joke{}, err
	}

	if len(jokes) == 0 {
		return Joke{}, fmt.Errorf("jokes with category '%s', type '%s' wasn't found", search.Category, search.Type)
	}

	return jokes[rand.Int()%len(jokes)], nil
}

func lockJokeService(ctx context.Context, name string, resetTime time.Time) {
	file.CreateLockForService(ctx, name)

	deadlineCtx, cancel := context.WithDeadline(ctx, resetTime)
	defer cancel()

	for {
		select {
		case <-deadlineCtx.Done():
			file.RemoveLockForService(ctx, name)
		}
	}
}
