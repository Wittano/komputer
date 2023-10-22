package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/wittano/komputer/internal/joke"
	"github.com/wittano/komputer/internal/log"
	"github.com/wittano/komputer/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"math/rand"
	"os"
)

type HumorAPIJokeDB struct {
	Category types.JokeCategory
	Joke     joke.HumorAPIRes
}

func AddNewJokesFromHumorAPI(ctx context.Context) {
	categories := []types.JokeCategory{types.ANY, types.PROGRAMMING, types.DARK}

	for {
		c := categories[rand.Int()%len(categories)]

		j, err := joke.GetRandomJokeFromHumorAPI(ctx, c)
		if err != nil && errors.Is(err, joke.HumorAPILimitExceededErr{}) {
			log.Warn(ctx, "Limit of getting jokes from HumorAPI was exceeded")
			return
		} else if err != nil && errors.Is(err, joke.HumorAPIKeyMissingErr{}) {
			log.Warn(ctx, err.Error())
			return
		} else if err != nil {
			log.Error(ctx, "Failed get joke from HumorAPI", err)
			continue
		}

		go func(ctx context.Context, c types.JokeCategory, j joke.HumorAPIRes) {
			hJoke := HumorAPIJokeDB{
				Category: c,
				Joke:     j,
			}

			if err := saveJokeFromRapidAPI(ctx, hJoke); err != nil {
				log.Error(ctx, "Failed save a new joke into database", err)
			}
		}(ctx, c, j)
	}
}

func saveJokeFromRapidAPI(ctx context.Context, jokeAPI HumorAPIJokeDB) error {
	name, ok := os.LookupEnv("MONGODB_DB_NAME")
	if !ok {
		return errors.New("required environment variable 'MONGODB_DB_NAME' is missing")
	}

	filter := bson.D{{
		Key: "externalID", Value: jokeAPI.Joke.ID,
	}}

	res, err := client.Database(name).Collection(jokeCollectionName).Find(ctx, filter)
	if err != nil {
		return err
	}

	if res != nil && res.Current != nil {
		return errors.New(fmt.Sprintf("JokeContainer with externalID %d exists in database", jokeAPI.Joke.ID))
	}

	j := JokeDB{
		Type:       types.Single,
		ContentRes: jokeAPI.Joke.Content(),
		Category:   jokeAPI.Category,
		ExternalID: jokeAPI.Joke.ID,
	}

	_, err = client.Database(name).Collection(jokeCollectionName).InsertOne(ctx, j)
	if err != nil {
		return err
	}

	log.Info(ctx, fmt.Sprintf("JokeContainer with ID %d was saved", jokeAPI.Joke.ID))

	return nil
}
