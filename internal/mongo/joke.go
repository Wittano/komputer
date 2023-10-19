package mongo

import (
	"context"
	"errors"
	"github.com/wittano/komputer/internal/joke"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
)

const collectionName = "joke"

type JokeDb struct {
	id       primitive.ObjectID `bson:"_id"`
	Question string             `bson:"question"`
	Content  string             `bson:"content"`
	Type     joke.JokeType      `bson:"type"`
	Category joke.JokeCategory  `bson:"category"`
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
