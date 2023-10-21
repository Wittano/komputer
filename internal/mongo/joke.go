package mongo

import (
	"context"
	"errors"
	"github.com/wittano/komputer/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"math/rand"
	"os"
)

type JokeDB struct {
	id         primitive.ObjectID `bson:"_id"`
	Question   string             `bson:"question"`
	ContentRes string             `bson:"content"`
	Type       types.JokeType     `bson:"type"`
	Category   types.JokeCategory `bson:"category"`
	ExternalID int64              `bson:"externalID"`
}

func (j JokeDB) Content() string {
	return j.ContentRes
}

func (j JokeDB) ContentTwoPart() (string, string) {
	return j.Question, j.ContentRes
}

func (j JokeDB) toJokeSearch() JokeSearch {
	return JokeSearch{
		Type:     j.Type,
		Category: j.Category,
	}
}

type JokeSearch struct {
	Type     types.JokeType
	Category types.JokeCategory
}

func AddJoke(ctx context.Context, joke JokeDB) (id primitive.ObjectID, err error) {
	name, ok := os.LookupEnv("MONGODB_DB_NAME")
	if !ok {
		return id, errors.New("required environment variable 'MONGODB_DB_NAME' is missing")
	}

	res, err := client.Database(name).Collection(jokeCollectionName).InsertOne(ctx, joke)

	id = res.InsertedID.(primitive.ObjectID)

	return
}

func GetSingleTypeJoke(ctx context.Context, category types.JokeCategory) (types.Joke, error) {
	return findRandomJoke(ctx, JokeSearch{Type: types.Single, Category: category})
}

func GetTwoPartsTypeJoke(ctx context.Context, category types.JokeCategory) (types.JokeTwoParts, error) {
	return findRandomJoke(ctx, JokeSearch{Type: types.TwoPart, Category: category})
}

func findRandomJoke(ctx context.Context, j JokeSearch) (JokeDB, error) {
	name, ok := os.LookupEnv("MONGODB_DB_NAME")
	if !ok {
		return JokeDB{}, errors.New("required environment variable 'MONGODB_DB_NAME' is missing")
	}

	var category = j.Category
	if category == "" {
		category = types.ANY
	}

	var jokeType = j.Type
	if jokeType == "" {
		jokeType = types.Single
	}

	pipelines := mongo.Pipeline{{{
		"$sample", bson.D{{
			"size", 10,
		}},
	}},
		{{
			"$match", bson.D{{
				"type", jokeType,
			}},
		}},
	}

	if category != "" && category != types.ANY {
		pipelines = append(pipelines, bson.D{{
			"$match", bson.D{{
				"category", category,
			}},
		}})
	}

	c, err := client.Database(name).Collection(jokeCollectionName).Aggregate(ctx, pipelines)
	if err != nil {
		return JokeDB{}, err
	}

	var res []JokeDB
	if err = c.All(ctx, &res); err != nil {
		return JokeDB{}, err
	}

	if len(res) == 0 {
		return JokeDB{}, types.JokeNotFoundErr{Category: category, JokeType: jokeType}
	}

	return res[rand.Int()%len(res)], nil
}
