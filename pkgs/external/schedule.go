package external

import (
	"context"
	"github.com/wittano/komputer/internal/mongo"
	"time"
)

const dayInterval = time.Hour * 24

func StartDownloadingJokeFromHumorAPI(ctx context.Context) {
	timer := time.NewTimer(dayInterval)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			break
		case <-timer.C:
			mongo.AddNewJokesFromHumorAPI(ctx)
		}
	}
}
