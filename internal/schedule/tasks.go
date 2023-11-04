package schedule

import (
	"context"
	"github.com/go-co-op/gocron"
	"github.com/wittano/komputer/internal/log"
	"github.com/wittano/komputer/internal/mongo"
	"os"
	"time"
)

var Scheduler = gocron.NewScheduler(time.UTC)

func init() {
	if _, ok := os.LookupEnv("RAPID_API_KEY"); ok {
		_, err := Scheduler.Every(1).Day().Do(mongo.AddNewJokesFromHumorAPI, context.Background())
		if err != nil {
			log.Error(context.Background(), "Failed to start getting new jokes from HumorAPI", err)
		}
	}
}
