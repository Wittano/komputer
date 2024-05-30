package joke

import (
	"context"
	"github.com/wittano/komputer/db"
	"github.com/wittano/komputer/db/joke"
	"net/http"
	"os"
	"time"
)

func NewJokeDevService(globalCtx context.Context) joke.SearchService {
	client := http.Client{Timeout: time.Second * 1}

	return &DevService{client, true, globalCtx}
}

func NewHumorAPIService(globalCtx context.Context) joke.SearchService {
	client := http.Client{Timeout: time.Second * 1}

	env, ok := os.LookupEnv(humorAPIKey)
	active := ok || env != ""

	return &HumorAPIService{client, active, globalCtx}
}

func NewDatabaseJokeService(database db.MongodbService) joke.Database {
	return joke.Database{Mongodb: database}
}

func unlockService(ctx context.Context, activeFlag *bool, resetTime time.Time) {
	deadlineCtx, cancel := context.WithDeadline(ctx, resetTime)
	defer cancel()

	for {
		if *activeFlag {
			return
		}

		select {
		case <-deadlineCtx.Done():
			*activeFlag = true
			return
		default:
		}
	}
}
