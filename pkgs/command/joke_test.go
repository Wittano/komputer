package command

import (
	"context"
	"errors"
	"github.com/wittano/komputer/pkgs/joke"
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
	testServices := []joke.GetService{
		joke.NewJokeDevService(ctx),
		joke.NewHumorAPIService(ctx),
		joke.NewDatabaseJokeService(dumpMongoService{}),
	}

	service, err := selectGetService(ctx, testServices)
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range testServices {
		if v == service {
			return
		}
	}

	t.Fatal("GetService wasn't found")
}

func TestSelectGetServiceButContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	testServices := []joke.GetService{
		joke.NewJokeDevService(ctx),
		joke.NewHumorAPIService(ctx),
		joke.NewDatabaseJokeService(dumpMongoService{}),
	}

	if _, err := selectGetService(ctx, testServices); err == nil {
		t.Fatal("Some services was found, but shouldn't")
	}
}

func TestSelectGetServiceButServicesIsDeactivated(t *testing.T) {
	ctx := context.Background()
	testServices := []joke.GetService{
		joke.NewDatabaseJokeService(dumpMongoService{}),
	}

	if _, err := selectGetService(ctx, testServices); err == nil {
		t.Fatal("Some services was found, but shouldn't")
	}
}
