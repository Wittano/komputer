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

type mapper interface {
	Joke() Joke
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

func (j singleResponse) Joke() Joke {
	return Joke{
		Answer:   j.Content,
		Type:     Single,
		Category: Category(j.Category),
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

func (j twoPartResponse) Joke() Joke {
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

func (d *DevService) Joke(ctx context.Context, params SearchParams) (Joke, error) {
	select {
	case <-ctx.Done():
		return Joke{}, context.Canceled
	default:
	}

	if !d.Active(ctx) {
		return Joke{}, DevServiceLimitExceededErr
	}

	if params.Category == YOMAMA || params.Category == "" {
		params.Category = Any
	}

	url := fmt.Sprintf(jokeDevAPIUrlTemplate, params.Category, params.Type)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Joke{}, err
	}

	res, err := d.client.Do(req)
	if err != nil {
		return Joke{}, err
	}
	defer res.Body.Close()

	// Check if daily limit exceeded
	isLimitExceeded := len(res.Header[rateLimitRemainingHeaderName]) > 0 &&
		res.Header[rateLimitRemainingHeaderName][0] == "0"

	if res.StatusCode == http.StatusTooManyRequests || isLimitExceeded {
		const rateLimitReset = "RateLimit-Reset"
		resetTime := resetTime(res.Header[rateLimitReset])
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

	var mapper mapper
	switch params.Type {
	case Single:
		singleRes := singleResponse{}

		err = json.Unmarshal(resBody, &singleRes)

		mapper = singleRes
	case TwoPart:
		twoPartRes := &twoPartResponse{}

		err = json.Unmarshal(resBody, &twoPartRes)

		mapper = twoPartRes
	}

	if err != nil {
		return Joke{}, err
	}

	return mapper.Joke(), nil
}

// Prepare reset time after that service should be activated
func resetTime(rateLimitReset []string) (t time.Time) {
	if len(rateLimitReset) > 0 {
		var err error
		if t, err = time.Parse("Sun, 06 Nov 1994 08:49:37 GMT", rateLimitReset[0]); err != nil {
			t = time.Now().Add(24 * time.Hour)
		}
	} else {
		t = time.Now().Add(24 * time.Hour)
	}

	return t
}
