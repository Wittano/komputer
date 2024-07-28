package joke

import (
	"context"
	"github.com/wittano/komputer/internal"
	"github.com/wittano/komputer/internal/mongodb"
	"net/http"
	"os"
	"sync"
	"time"
)

func NewJokeDevService(globalCtx context.Context) internal.SearchService {
	return &DevService{
		client:    http.Client{Timeout: time.Second},
		active:    true,
		globalCtx: globalCtx,
	}
}

func NewHumorAPIService(globalCtx context.Context) internal.SearchService {
	env, ok := os.LookupEnv(humorAPIKey)

	return &HumorAPIService{
		client:    http.Client{Timeout: time.Second},
		active:    ok || env != "",
		globalCtx: globalCtx,
	}
}

func NewJokeDatabase(database mongodb.DatabaseGetter) mongodb.Service {
	return mongodb.Service{Db: database}
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
