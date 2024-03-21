package db

import (
	"context"
	"os"
	"testing"
)

const uri = "mongodb+srv://username:password@server/database"

func TestNewMongodbDatabase(t *testing.T) {
	os.Setenv(mongodbURIKey, uri)

	if db := NewMongodbDatabase(context.Background()); db.uri != uri {
		t.Fatalf("Invalid URI. Expected: '%s', Result: '%s'", uri, db.uri)
	}
}

func TestNewMongodbDatabaseButMongodbURIMissing(t *testing.T) {
	if db := NewMongodbDatabase(context.Background()); db.uri != uri {
		t.Fatalf("Invalid URI. Expected: '%s', Result: '%s'", uri, db.uri)
	}
}

func TestMongodbDatabase_ClientButURIMissing(t *testing.T) {
	ctx := context.Background()

	db := NewMongodbDatabase(ctx)

	if _, err := db.Client(ctx); err == nil {
		t.Fatal("mongodb connection was established!")
	}
}

func TestMongodbDatabase_ClientButConnectionFailed(t *testing.T) {
	ctx := context.Background()
	os.Setenv(mongodbURIKey, uri)

	db := NewMongodbDatabase(ctx)

	if _, err := db.Client(ctx); err == nil {
		t.Fatal("mongodb connection was established!")
	}
}
