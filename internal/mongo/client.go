package mongo

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/wittano/komputer/internal/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

var client *mongo.Client

func init() {
	var err error

	if err = godotenv.Load(); err != nil {
		log.Error(context.Background(), "No .env file found", err)
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Error(context.Background(), "You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable", nil)
	}

	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(fmt.Sprintf(uri)))
	if err != nil {
		log.Error(context.Background(), "Failed to connect MongoDB database. Check error message", err)
	}
}

func CloseDb() {
	if err := client.Disconnect(context.Background()); err != nil {
		panic(err)
	}
}
