package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/wittano/komputer/internal/joke"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
	"os"
)

const collectionName = "joke"

type JokeDb struct {
	id         primitive.ObjectID `bson:"_id"`
	Question   string             `bson:"question"`
	ContentRes string             `bson:"content"`
	Type       joke.JokeType      `bson:"type"`
	Category   joke.JokeCategory  `bson:"category"`
}

func (j JokeDb) Content() string {
	return j.ContentRes
}

func (j JokeDb) ContentTwoPart() (string, string) {
	return j.Question, j.ContentRes
}

func (j JokeDb) toJokeSearch() JokeSearch {
	return JokeSearch{
		Type:     j.Type,
		Category: j.Category,
	}
}

type JokeSearch struct {
	Type     joke.JokeType
	Category joke.JokeCategory
}

type JokeNotFoundErr struct {
	category joke.JokeCategory
	jokeType joke.JokeType
}

func (j JokeNotFoundErr) Error() string {
	return fmt.Sprintf("Joke %s from category %s wasn't found", j.jokeType, j.category)
}

func AddJoke(ctx context.Context, joke JokeDb) (id primitive.ObjectID, err error) {
	name, ok := os.LookupEnv("MONGODB_DB_NAME")
	if !ok {
		return id, errors.New("required environment variable 'MONGODB_DB_NAME' is missing")
	}

	res, err := client.Database(name).Collection(collectionName).InsertOne(ctx, joke)

	id = res.InsertedID.(primitive.ObjectID)

	return
}

func FindRandomJoke(ctx context.Context, j JokeSearch) (JokeDb, error) {
	name, ok := os.LookupEnv("MONGODB_DB_NAME")
	if !ok {
		return JokeDb{}, errors.New("required environment variable 'MONGODB_DB_NAME' is missing")
	}

	var category = j.Category
	if category == "" {
		category = joke.ANY
	}

	var jokeType = j.Type
	if jokeType == "" {
		jokeType = joke.Single
	}

	pipelines := bson.D{
		{
			"type", jokeType,
		},
	}

	if category != "" && category != joke.ANY {
		pipelines = append(pipelines, bson.E{
			Key: "category", Value: category,
		})
	}

	c, err := client.Database(name).Collection(collectionName).Find(ctx, pipelines)
	if err != nil {
		return JokeDb{}, err
	}

	var res []JokeDb
	if err = c.All(ctx, &res); err != nil {
		return JokeDb{}, err
	}

	if len(res) == 0 {
		return JokeDb{}, JokeNotFoundErr{category, jokeType}
	}

	return res[rand.Int()%len(res)], nil
}
