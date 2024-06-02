package joke

import (
	"context"
	"github.com/wittano/komputer/db"
	"github.com/wittano/komputer/db/joke"
	"net/http"
	"os"
	"sync"
	"time"
)

func NewJokeDevService(globalCtx context.Context) joke.SearchService {
	return &DevService{
		client:    http.Client{Timeout: time.Second},
		active:    true,
		globalCtx: globalCtx,
	}
}

func NewHumorAPIService(globalCtx context.Context) joke.SearchService {
	env, ok := os.LookupEnv(humorAPIKey)

	return &HumorAPIService{
		client:    http.Client{Timeout: time.Second},
		active:    ok || env != "",
		globalCtx: globalCtx,
	}
}

func NewJokeDatabase(database db.MongodbService) joke.Database {
	return joke.Database{Mongodb: database}
}

func unlockService(ctx context.Context, m *sync.Mutex, activeFlag *bool, resetTime time.Time) {
	deadlineCtx, cancel := context.WithDeadline(ctx, resetTime)
	defer cancel()

	for {
		if *activeFlag {
			return
		}

		select {
		case <-deadlineCtx.Done():
			m.Lock()
			*activeFlag = true
			m.Unlock()
			return
		default:
		}
	}
}
