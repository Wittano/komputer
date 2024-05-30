package joke

import (
	"context"
	"errors"
	"fmt"
	komputer "github.com/wittano/komputer/api/proto"
	"github.com/wittano/komputer/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"math/rand"
	"time"
)

const (
	collectionName = "jokes"
	GuildIDKey     = "guildID"
)

type Database struct {
	Mongodb db.MongodbService
}

func (d Database) Active(ctx context.Context) bool {
	const maxTimeoutTime = 500

	ctx, cancel := context.WithTimeout(ctx, maxTimeoutTime*time.Millisecond)
	defer cancel()

	client, err := d.Mongodb.Client(ctx)
	if err != nil {
		return false
	}

	err = client.Ping(ctx, readpref.Nearest(readpref.WithMaxStaleness(maxTimeoutTime)))
	if err != nil {
		return false
	}

	return true
}

func (d Database) Add(ctx context.Context, joke Joke) (string, error) {
	select {
	case <-ctx.Done():
		return "", context.Canceled
	default:
	}

	mongodb, err := d.Mongodb.Client(ctx)
	if err != nil {
		return "", err
	}

	res, err := mongodb.Database(db.DatabaseName).Collection(collectionName).InsertOne(ctx, joke)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (d Database) RandomJoke(ctx context.Context, search SearchParams) (Joke, error) {
	jokes, err := d.Jokes(ctx, search, nil)
	if err != nil {
		return Joke{}, err
	}

	return jokes[rand.Int()%len(jokes)], nil
}

func (d Database) Joke(ctx context.Context, search SearchParams) (Joke, error) {
	jokes, err := d.Jokes(ctx, search, nil)
	if err != nil {
		return Joke{}, err
	}

	return jokes[0], nil
}

func (d Database) Jokes(ctx context.Context, search SearchParams, page *komputer.Pagination) ([]Joke, error) {
	select {
	case <-ctx.Done():
		return nil, context.Canceled
	default:
	}

	if !d.Active(ctx) {
		return nil, errors.New("databases isn't responding")
	}

	mongodb, err := d.Mongodb.Client(ctx)
	if err != nil {
		return nil, err
	}

	if search.Category == "" {
		search.Category = Any
	}

	if search.Type == "" {
		search.Type = Single
	}

	// Create query to database
	const matchQueryKey = "$match"

	var (
		pageSize uint32 = 10
		pageNr   uint32 = 0
	)

	if page != nil {
		if page.Size > 0 {
			pageSize = page.Size
		}
		if page.Page > 0 {
			pageNr = page.Page
		}
	}

	pipeline := mongo.Pipeline{{{
		"$sample", bson.D{{
			"size", pageSize,
		}},
	}},
		{{Key: "$skip", Value: pageNr * pageSize}},
		{{
			matchQueryKey, bson.D{
				{
					"type", search.Type,
				},
				{
					"guild_id", ctx.Value(GuildIDKey),
				},
			},
		}},
	}

	if search.ID != [12]byte{} {
		pipeline = append(pipeline, bson.D{{
			matchQueryKey, bson.D{{
				"_id", search.ID,
			}},
		}})
	}

	if search.Category != "" && search.Category != Any {
		pipeline = append(pipeline, bson.D{{
			matchQueryKey, bson.D{{
				"category", search.Category,
			}},
		}})
	}

	// SearchParams
	res, err := mongodb.Database(db.DatabaseName).Collection(collectionName).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer res.Close(ctx)

	var jokes []Joke
	if err = res.All(ctx, &jokes); err != nil {
		return nil, err
	}

	if len(jokes) == 0 {
		return nil, fmt.Errorf("jokes with category '%s', type '%s' wasn't found", search.Category, search.Type)
	}

	return jokes, nil
}

func (d Database) Remove(ctx context.Context, id string) error {
	select {
	case <-ctx.Done():
		return context.Canceled
	default:
	}

	if !d.Active(ctx) {
		return errors.New("databases isn't responding")
	}

	client, err := d.Mongodb.Client(ctx)
	if err != nil {
		return err
	}

	filter := bson.D{{"$match", bson.D{{
		"_id", id,
	}}}}

	_, err = client.Database(db.DatabaseName).Collection(collectionName).DeleteOne(ctx, filter)
	return err
}
