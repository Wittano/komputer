package joke

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wittano/komputer/internal/file"
	"io"
	"net/http"
	"time"
)

const devServiceName = "jokedev"

type jokeMapper interface {
	Joke() Joke
}

type jokeApiFlags struct {
	Nsfw      bool `json:"nsfw"`
	Religious bool `json:"religious"`
	Political bool `json:"political"`
	Racist    bool `json:"racist"`
	Sexist    bool `json:"sexist"`
	Explicit  bool `json:"explicit"`
}

type jokeApiSingleResponse struct {
	Error       bool         `json:"error"`
	CategoryRes string       `json:"category"`
	Type        string       `json:"type"`
	Flags       jokeApiFlags `json:"flags"`
	Id          int          `json:"id"`
	Safe        bool         `json:"safe"`
	Lang        string       `json:"lang"`
	Content     string       `json:"joke"`
}

func (j jokeApiSingleResponse) Joke() Joke {
	return Joke{
		Answer:   j.Content,
		Type:     Single,
		Category: Category(j.CategoryRes),
	}
}

type jokeApiTwoPartResponse struct {
	Error    bool         `json:"error"`
	Category string       `json:"category"`
	Type     string       `json:"type"`
	Flags    jokeApiFlags `json:"flags"`
	Id       int          `json:"id"`
	Safe     bool         `json:"safe"`
	Lang     string       `json:"lang"`
	Setup    string       `json:"setup"`
	Delivery string       `json:"delivery"`
}

func (j jokeApiTwoPartResponse) Joke() Joke {
	return Joke{
		Question: j.Setup,
		Answer:   j.Delivery,
		Type:     Single,
		Category: Category(j.Category),
	}
}

var DevServiceLimitExceededErr = errors.New("jokedev: current limit exceeded")

type DevService struct {
	client    http.Client
	globalCtx context.Context
}

func (d DevService) Active(_ context.Context) bool {
	return !file.IsServiceLocked(devServiceName)
}

func (d DevService) Get(ctx context.Context, search SearchParameters) (Joke, error) {
	select {
	case <-ctx.Done():
		return Joke{}, context.Canceled
	default:
	}

	if !d.Active(ctx) {
		return Joke{}, DevServiceLimitExceededErr
	}

	if search.Category == YOMAMA {
		search.Category = Any
	}

	jokeDevApiURL := fmt.Sprintf("https://v2.jokeapi.dev/joke/%s?type=%s", search.Category, search.Category)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, jokeDevApiURL, nil)
	if err != nil {
		return Joke{}, err
	}

	res, err := d.client.Do(req)
	if err != nil {
		return Joke{}, err
	}
	defer res.Body.Close()

	// Check if daily limit exceeded
	if res.StatusCode == http.StatusTooManyRequests || res.Header["RateLimit-Remaining"][0] == "0" {
		resetTime, err := time.Parse("Sun, 06 Nov 1994 08:49:37 GMT", res.Header["RateLimit-Reset"][0])
		if err != nil {
			return Joke{}, err
		}

		go lockJokeService(d.globalCtx, devServiceName, resetTime)
		return Joke{}, DevServiceLimitExceededErr
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return Joke{}, err
	}

	var jokeMapper jokeMapper
	switch search.Type {
	case Single:
		jokeMapper = jokeApiSingleResponse{}
	case TwoPart:
		jokeMapper = jokeApiTwoPartResponse{}
	}

	err = json.Unmarshal(resBody, &jokeMapper)
	if err != nil {
		return Joke{}, err
	}

	return jokeMapper.Joke(), nil
}
