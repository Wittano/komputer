package joke

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	rateLimitRemainingHeaderName = "RateLimit-Remaining"
	jokeDevAPIUrlTemplate        = "https://v2.jokeapi.dev/joke/%s?type=%s"
)

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
	Error    bool         `json:"error"`
	Category string       `json:"category"`
	Type     string       `json:"type"`
	Flags    jokeApiFlags `json:"flags"`
	Id       int          `json:"id"`
	Safe     bool         `json:"safe"`
	Lang     string       `json:"lang"`
	Content  string       `json:"joke"`
}

func (j jokeApiSingleResponse) Joke() Joke {
	return Joke{
		Answer:   j.Content,
		Type:     Single,
		Category: Category(j.Category),
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
	active    bool
	globalCtx context.Context
}

func (d DevService) Active(ctx context.Context) (active bool) {
	select {
	case <-ctx.Done():
		active = false
	default:
		active = d.active
	}

	return
}

func (d *DevService) Get(ctx context.Context, search SearchParameters) (Joke, error) {
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

	jokeDevApiURL := fmt.Sprintf(jokeDevAPIUrlTemplate, search.Category, search.Type)
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
	isLimitExceeded := len(res.Header[rateLimitRemainingHeaderName]) > 0 && res.Header[rateLimitRemainingHeaderName][0] == "0"

	if res.StatusCode == http.StatusTooManyRequests || isLimitExceeded {
		const rateLimitReset = "RateLimit-Reset"
		resetTime := prepareResetTime(res.Header[rateLimitReset])
		d.active = false

		go unlockService(d.globalCtx, &d.active, resetTime)

		return Joke{}, DevServiceLimitExceededErr
	}

	if res.StatusCode >= 400 {
		return Joke{}, errors.New("jokedev: client or server side error")
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return Joke{}, err
	}

	var jokeMapper jokeMapper
	switch search.Type {
	case Single:
		singleRes := jokeApiSingleResponse{}

		err = json.Unmarshal(resBody, &singleRes)

		jokeMapper = singleRes
	case TwoPart:
		twoPartRes := &jokeApiTwoPartResponse{}

		err = json.Unmarshal(resBody, &twoPartRes)

		jokeMapper = twoPartRes
	}

	if err != nil {
		return Joke{}, err
	}

	return jokeMapper.Joke(), nil
}

func prepareResetTime(rateLimitReset []string) (resetTime time.Time) {
	var err error
	if len(rateLimitReset) > 0 {
		resetTime, err = time.Parse("Sun, 06 Nov 1994 08:49:37 GMT", rateLimitReset[0])
		if err != nil {
			resetTime = time.Now().Add(24 * time.Hour)
		}
	} else {
		resetTime = time.Now().Add(24 * time.Hour)
	}

	return resetTime
}
