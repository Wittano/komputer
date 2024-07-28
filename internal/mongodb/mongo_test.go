package mongodb

import (
	"context"
	"os"
	"testing"
)

const testUri = "NewMongodb+srv://username:password@server/database"

func TestNewMongodbDatabase(t *testing.T) {
	os.Setenv(MongodbURIKey, testUri)

	if db := NewMongodb(context.Background()); db.uri != testUri {
		t.Fatalf("Invalid URI. Expected: '%s', Result: '%s'", testUri, db.uri)
	}
}

func TestNewMongodbDatabaseButMongodbURIMissing(t *testing.T) {
	if db := NewMongodb(context.Background()); db.uri != testUri {
		t.Fatalf("Invalid URI. Expected: '%s', Result: '%s'", testUri, db.uri)
	}
}

func TestMongodbDatabase_ClientButURIMissing(t *testing.T) {
	ctx := context.Background()

	db := NewMongodb(ctx)

	if _, err := db.Client(ctx); err == nil {
		t.Fatal("NewMongodb connection was established!")
	}
}

func TestMongodbDatabase_ClientButConnectionFailed(t *testing.T) {
	ctx := context.Background()
	os.Setenv(MongodbURIKey, testUri)

	db := NewMongodb(ctx)

	if _, err := db.Client(ctx); err == nil {
		t.Fatal("NewMongodb connection was established!")
	}
}
