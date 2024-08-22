package internal

import (
	"context"
	"github.com/wittano/komputer/bot/log"
	"github.com/wittano/komputer/internal/joke"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddService interface {
	Add(ctx context.Context, joke joke.DbModel) (string, error)
}

type SearchService interface {
	// RandomJoke Joke Try to find Joke from Db database. If SearchParams is empty, then function will find 1 random joke
	RandomJoke(ctx log.Context, search SearchParams) (joke.DbModel, error)
	ActiveChecker
}

type ActiveChecker interface {
	Active(ctx context.Context) bool
}

type SearchParams struct {
	Type     joke.Type
	Category joke.Category
	ID       primitive.ObjectID
}
