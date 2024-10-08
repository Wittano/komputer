package joke

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wittano/komputer/bot/log"
	"github.com/wittano/komputer/internal"
	"github.com/wittano/komputer/internal/joke"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	rateLimitRemainingHeaderName = "RateLimit-Remaining"
	jokeDevAPIUrlTemplate        = "https://v2.jokeapi.dev/joke/%s?type=%s"
)

type mapper interface {
	Joke() joke.DbModel
}

type flags struct {
	Nsfw      bool `json:"nsfw"`
	Religious bool `json:"religious"`
	Political bool `json:"political"`
	Racist    bool `json:"racist"`
	Sexist    bool `json:"sexist"`
	Explicit  bool `json:"explicit"`
}

type singleResponse struct {
	Error    bool   `json:"error"`
	Category string `json:"category"`
	Type     string `json:"type"`
	Flags    flags  `json:"flags"`
	Id       int    `json:"id"`
	Safe     bool   `json:"safe"`
	Lang     string `json:"lang"`
	Content  string `json:"joke"`
}

func (j singleResponse) Joke() joke.DbModel {
	return joke.DbModel{
		Answer:   j.Content,
		Type:     joke.Single,
		Category: joke.Category(j.Category),
	}
}

type twoPartResponse struct {
	Error    bool   `json:"error"`
	Category string `json:"category"`
	Type     string `json:"type"`
	Flags    flags  `json:"flags"`
	Id       int    `json:"id"`
	Safe     bool   `json:"safe"`
	Lang     string `json:"lang"`
	Setup    string `json:"setup"`
	Delivery string `json:"delivery"`
}

func (j twoPartResponse) Joke() joke.DbModel {
	return joke.DbModel{
		Question: j.Setup,
		Answer:   j.Delivery,
		Type:     joke.Single,
		Category: joke.Category(j.Category),
	}
}

var DevServiceLimitExceededErr = errors.New("jokedev: current limit exceeded")

type DevService struct {
	client    http.Client
	active    bool
	m         sync.Mutex
	globalCtx context.Context
}

func (d *DevService) Active(ctx context.Context) (active bool) {
	select {
	case <-ctx.Done():
		active = false
	default:
		active = d.active
	}

	return
}

func (d *DevService) RandomJoke(ctx log.Context, params internal.SearchParams) (joke.DbModel, error) {
	select {
	case <-ctx.Done():
		return joke.DbModel{}, context.Canceled
	default:
	}

	if !d.Active(ctx) {
		return joke.DbModel{}, DevServiceLimitExceededErr
	}

	if params.Category == joke.YOMAMA || params.Category == "" {
		params.Category = joke.Any
	}

	url := fmt.Sprintf(jokeDevAPIUrlTemplate, params.Category, params.Type)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return joke.DbModel{}, err
	}

	res, err := d.client.Do(req)
	if err != nil {
		return joke.DbModel{}, err
	}
	defer res.Body.Close()

	// Check if daily limit exceeded
	isLimitExceeded := len(res.Header[rateLimitRemainingHeaderName]) > 0 &&
		res.Header[rateLimitRemainingHeaderName][0] == "0"

	if res.StatusCode == http.StatusTooManyRequests || isLimitExceeded {
		const rateLimitReset = "RateLimit-Reset"
		d.m.Lock()
		resetTime := resetTime(res.Header[rateLimitReset])
		d.active = false

		go unlockService(d.globalCtx, &d.m, &d.active, resetTime)
		d.m.Unlock()

		return joke.DbModel{}, DevServiceLimitExceededErr
	}

	if res.StatusCode >= 400 {
		return joke.DbModel{}, errors.New("jokedev: client or server side error")
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return joke.DbModel{}, err
	}

	var mapper mapper
	switch params.Type {
	case joke.Single:
		singleRes := singleResponse{}

		err = json.Unmarshal(resBody, &singleRes)

		mapper = singleRes
	case joke.TwoPart:
		twoPartRes := &twoPartResponse{}

		err = json.Unmarshal(resBody, &twoPartRes)

		mapper = twoPartRes
	}

	if err != nil {
		return joke.DbModel{}, err
	}

	return mapper.Joke(), nil
}

// Prepare reset time after that service should be activated
func resetTime(rateLimitReset []string) (t time.Time) {
	if len(rateLimitReset) > 0 {
		var err error
		if t, err = time.Parse("Sun, 06 Nov 1994 08:49:37 GMT", rateLimitReset[0]); err == nil {
			return
		}
	}

	return time.Now().Add(24 * time.Hour)
}
