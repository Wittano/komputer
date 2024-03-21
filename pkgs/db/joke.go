package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"math/rand"
)

const (
	databaseName   = "komputer"
	collectionName = "jokes"
	GuildIDKey     = "guildID"
)

type (
	JokeType     string
	JokeCategory string
)

func (j JokeCategory) ToHumorAPICategory() string {
	switch j {
	case PROGRAMMING:
		return "nerdy"
	case Any:
		return "one_liner"
	case DARK:
		return "dark"
	case YOMAMA:
		return "yo_mama"
	default:
		return "one_liner"
	}
}

const (
	Single  JokeType = "single"
	TwoPart JokeType = "twopart"
)

const (
	PROGRAMMING JokeCategory = "Programming"
	MISC        JokeCategory = "Misc"
	DARK        JokeCategory = "Dark"
	YOMAMA      JokeCategory = "YoMama"
	Any         JokeCategory = "Any"
)

type Joke struct {
	id       primitive.ObjectID `bson:"_id"`
	Question string             `bson:"question"`
	Answer   string             `bson:"answer"`
	Type     JokeType           `bson:"type"`
	Category JokeCategory       `bson:"category"`
	GuildID  string             `bson:"guild_id"`
}

type JokeSearch struct {
	Type     JokeType
	Category JokeCategory
	ID       primitive.ObjectID
}

// JokeService TODO Export JokeService as interface, after refactor external API
type JokeService struct {
	Mongodb MongodbService
}

func (j JokeService) Add(ctx context.Context, joke Joke) (primitive.ObjectID, error) {
	select {
	case <-ctx.Done():
		return [12]byte{}, context.Canceled
	default:
	}

	db, err := j.Mongodb.Client(ctx)
	if err != nil {
		return [12]byte{}, err
	}

	res, err := db.Database(databaseName).Collection(collectionName).InsertOne(ctx, joke)
	if err != nil {
		return [12]byte{}, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}

func (j JokeService) Get(ctx context.Context, search JokeSearch) (Joke, error) {
	select {
	case <-ctx.Done():
		return Joke{}, context.Canceled
	default:
	}

	db, err := j.Mongodb.Client(ctx)
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

	// Search
	res, err := db.Database(databaseName).Collection(collectionName).Aggregate(ctx, pipeline)
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
