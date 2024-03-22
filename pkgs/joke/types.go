package joke

import (
	"context"
	"github.com/wittano/komputer/pkgs/db"
	"net/http"
	"time"
)

type AddService interface {
	Add(ctx context.Context, joke Joke) (string, error)
}

type GetService interface {
	Get(ctx context.Context, search SearchParameters) (Joke, error)
	ActiveService
}

type ActiveService interface {
	Active(ctx context.Context) bool
}

func NewJokeDevService(globalCtx context.Context) GetService {
	client := http.Client{Timeout: time.Second * 1}

	return DevService{client, globalCtx}
}

func NewHumorAPIService(globalCtx context.Context) GetService {
	client := http.Client{Timeout: time.Second * 1}

	return HumorAPIService{client, globalCtx}
}

func NewDatabaseJokeService(database db.MongodbService) DatabaseJokeService {
	return DatabaseJokeService{mongodb: database}
}
