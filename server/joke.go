package server

import (
	"context"
	"errors"
	komputer "github.com/wittano/komputer/api/proto"
	"github.com/wittano/komputer/internal"
	"github.com/wittano/komputer/internal/joke"
	"github.com/wittano/komputer/internal/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/emptypb"
)

type jokeServer struct {
	Db mongodb.Service

	komputer.UnimplementedJokeServiceServer
}

type apiJokeParams interface {
	GetId() *komputer.ObjectID
	GetCategory() komputer.Category
	GetType() komputer.Type
}

func (j jokeServer) Find(ctx context.Context, params *komputer.JokeParams) (*komputer.Joke, error) {
	p, err := searchParams(params)
	if err != nil {
		return nil, err
	}

	entity, err := j.Db.Jokes(ctx, p, nil)
	if err != nil {
		return nil, err
	}

	return entity[0].ApiResponse()
}

func (j jokeServer) FindAll(identity *komputer.JokeParamsPagination, server komputer.JokeService_FindAllServer) error {
	p, err := searchParams(identity)
	if err != nil {
		return err
	}

	page := paginationOrDefault(identity.Page)
	jokes, err := j.Db.Jokes(server.Context(), p, page)
	if err != nil {
		return err
	}

	for _, jokeDb := range jokes {
		response, err := jokeDb.ApiResponse()
		if err != nil {
			continue
		}

		err = errors.Join(err, server.Send(response))
	}

	return err
}

func (j jokeServer) Remove(ctx context.Context, id *komputer.ObjectID) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, j.Db.Remove(ctx, id.ObjectId)
}

func (j jokeServer) Add(ctx context.Context, joke *komputer.Joke) (*komputer.JokeID, error) {
	newJoke, err := newJoke(joke)
	if err != nil {
		return nil, err
	}

	id, err := j.Db.Add(ctx, newJoke)
	if err != nil {
		return nil, err
	}

	return &komputer.JokeID{Id: &komputer.JokeID_ObjectId{ObjectId: []byte(id)}}, err
}

func newJoke(j *komputer.Joke) (new joke.DbModel, err error) {
	if j == nil {
		err = errors.New("missing joke data")
		return
	}

	new.Category, err = joke.RawCategory(j.Category)
	new.Type, err = joke.RawType(j.Type)
	new.Answer = j.Answer
	new.GuildID = j.GuildId
	new.ID = primitive.NewObjectID()

	if j.Question != nil {
		new.Question = *j.Question
	}

	return
}

func searchParams(params apiJokeParams) (p internal.SearchParams, err error) {
	if params == nil {
		err = errors.New("missing joke params")
		return
	}

	id := params.GetId()
	if id != nil {
		p.ID, err = primitive.ObjectIDFromHex(id.ObjectId)
	}

	p.Type, err = joke.RawType(params.GetType())
	p.Category, err = joke.RawCategory(params.GetCategory())

	return
}
