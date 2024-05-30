package db

import (
	"context"
	"os"
	"testing"
)

const testUri = "Mongodb+srv://username:password@server/database"

func TestNewMongodbDatabase(t *testing.T) {
	os.Setenv(MongodbURIKey, testUri)

	if db := Mongodb(context.Background()); db.uri != testUri {
		t.Fatalf("Invalid URI. Expected: '%s', Result: '%s'", testUri, db.uri)
	}
}

func TestNewMongodbDatabaseButMongodbURIMissing(t *testing.T) {
	if db := Mongodb(context.Background()); db.uri != testUri {
		t.Fatalf("Invalid URI. Expected: '%s', Result: '%s'", testUri, db.uri)
	}
}

func TestMongodbDatabase_ClientButURIMissing(t *testing.T) {
	ctx := context.Background()

	db := Mongodb(ctx)

	if _, err := db.Client(ctx); err == nil {
		t.Fatal("Mongodb connection was established!")
	}
}

func TestMongodbDatabase_ClientButConnectionFailed(t *testing.T) {
	ctx := context.Background()
	os.Setenv(MongodbURIKey, testUri)

	db := Mongodb(ctx)

	if _, err := db.Client(ctx); err == nil {
		t.Fatal("Mongodb connection was established!")
	}
}
