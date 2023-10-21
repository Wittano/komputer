package mongo

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/wittano/komputer/internal/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

var client *mongo.Client

const jokeCollectionName = "jokes"

func init() {
	var err error

	if err = godotenv.Load(); err != nil {
		log.Error(context.Background(), "No .env file found", err)
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal(context.Background(), "You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable", nil)
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	client, err = mongo.Connect(context.Background(), options.Client().
		ApplyURI(uri).
		SetServerAPIOptions(serverAPI))
	if err != nil {
		log.Fatal(context.Background(), "Failed to connect MongoDB database. Check error message", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(context.Background(), "Failed connect to MongoDB", err)
	}
}

func CloseDb() {
	if err := client.Disconnect(context.Background()); err != nil {
		panic(err)
	}
}
