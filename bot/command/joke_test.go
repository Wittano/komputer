package command

import (
	"context"
	"errors"
	"github.com/wittano/komputer/bot/joke"
	"github.com/wittano/komputer/internal"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

type dumpMongoService struct{}

func (d dumpMongoService) Close() error {
	return errors.New("not implemented")
}

func (d dumpMongoService) Client(_ context.Context) (*mongo.Client, error) {
	return nil, errors.New("not implemented")
}

func TestSelectGetService(t *testing.T) {
	ctx := context.Background()
	testServices := []internal.SearchService{
		joke.NewJokeDevService(ctx),
		joke.NewHumorAPIService(ctx),
		joke.NewJokeDatabase(dumpMongoService{}),
	}

	service, err := findService(ctx, testServices)
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range testServices {
		if v == service {
			return
		}
	}

	t.Fatal("SearchService wasn't found")
}

func TestFindJokeService_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	testServices := []internal.SearchService{
		joke.NewJokeDevService(ctx),
		joke.NewHumorAPIService(ctx),
		joke.NewJokeDatabase(dumpMongoService{}),
	}

	if _, err := findService(ctx, testServices); err == nil {
		t.Fatal("Some services was found, but shouldn't")
	}
}

func TestFindJokeService_ServicesIsDeactivated(t *testing.T) {
	ctx := context.Background()
	services := []internal.SearchService{
		joke.NewJokeDatabase(dumpMongoService{}),
	}

	if _, err := findService(ctx, services); err == nil {
		t.Fatal("Some services was found, but shouldn't")
	}
}
