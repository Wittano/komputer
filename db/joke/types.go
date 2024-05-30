package joke

import (
	"context"
	"errors"
	komputer "github.com/wittano/komputer/api/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddService interface {
	Add(ctx context.Context, joke Joke) (string, error)
}

type SearchService interface {
	// RandomJoke Joke Try to find Joke from Mongodb database. If SearchParams is empty, then function will find 1 random joke
	RandomJoke(ctx context.Context, search SearchParams) (Joke, error)
	ActiveChecker
}

type ActiveChecker interface {
	Active(ctx context.Context) bool
}

type (
	Type     string
	Category string
)

const (
	Single  Type = "single"
	TwoPart Type = "twopart"
)

const (
	PROGRAMMING Category = "Programming"
	MISC        Category = "Misc"
	DARK        Category = "Dark"
	YOMAMA      Category = "YoMama"
	Any         Category = "Any"
)

func (t Type) ApiType() (ty komputer.Type, err error) {
	switch t {
	case Single:
		ty = komputer.Type_SINGLE
	case TwoPart:
		ty = komputer.Type_TWO_PART
	default:
		err = errors.New("joke: unknown type")
	}

	return
}

func (c Category) ApiCategory() (ca komputer.Category, err error) {
	switch c {
	case Any:
		ca = komputer.Category_Any
	case DARK:
		ca = komputer.Category_DARK
	case PROGRAMMING:
		ca = komputer.Category_PROGRAMMING
	case YOMAMA:
		ca = komputer.Category_YOMAMA
	case MISC:
		ca = komputer.Category_MISC
	default:
		err = errors.New("joke: unknown category")
	}

	return
}

func RawType(api komputer.Type) (t Type, err error) {
	switch api {
	case komputer.Type_SINGLE:
		t = Single
	case komputer.Type_TWO_PART:
		t = TwoPart
	default:
		err = errors.New("joke: unknown type")
	}

	return
}

func RawCategory(api komputer.Category) (c Category, err error) {
	switch api {
	case komputer.Category_Any:
		c = Any
	case komputer.Category_DARK:
		c = DARK
	case komputer.Category_PROGRAMMING:
		c = PROGRAMMING
	case komputer.Category_YOMAMA:
		c = YOMAMA
	case komputer.Category_MISC:
		c = MISC
	default:
		err = errors.New("joke: unknown category")
	}

	return
}

type Joke struct {
	ID       primitive.ObjectID `bson:"_id"`
	Question string             `bson:"question"`
	Answer   string             `bson:"answer"`
	Type     Type               `bson:"type"`
	Category Category           `bson:"category"`
	GuildID  string             `bson:"guild_id"`
}

func (j Joke) ApiResponse() (*komputer.Joke, error) {
	ty, err := j.Type.ApiType()
	if err != nil {
		return nil, err
	}

	ca, err := j.Category.ApiCategory()
	if err != nil {
		return nil, err
	}

	return &komputer.Joke{
		Id:       &komputer.ObjectID{ObjectId: j.ID.Hex()},
		Answer:   j.Answer,
		Question: &j.Question,
		Type:     ty,
		Category: ca,
		GuildId:  j.GuildID,
	}, nil
}

type SearchParams struct {
	Type     Type
	Category Category
	ID       primitive.ObjectID
}
