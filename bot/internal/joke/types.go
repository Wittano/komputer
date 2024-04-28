package joke

import (
	"context"
	"github.com/wittano/komputer/bot/internal"
	"github.com/wittano/komputer/db"
	"net/http"
	"os"
	"time"
)

type AddService interface {
	Add(ctx context.Context, joke Joke) (string, error)
}

type SearchService interface {
	Joke(ctx context.Context, search SearchParams) (Joke, error)
	internal.ActiveChecker
}

func NewJokeDevService(globalCtx context.Context) SearchService {
	client := http.Client{Timeout: time.Second * 1}

	return &DevService{client, true, globalCtx}
}

func NewHumorAPIService(globalCtx context.Context) SearchService {
	client := http.Client{Timeout: time.Second * 1}

	env, ok := os.LookupEnv(humorAPIKey)
	active := ok || env != ""

	return &HumorAPIService{client, active, globalCtx}
}

func NewDatabaseJokeService(database db.MongodbService) DatabaseService {
	return DatabaseService{mongodb: database}
}
