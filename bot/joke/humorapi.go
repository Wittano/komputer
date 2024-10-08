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
	"os"
	"sync"
	"time"
)

const (
	humorApiURL             = "https://humor-jokes-and-memes.p.rapidapi.com/jokes/random?exclude-tags=nsfw&include-tags="
	humorAPIKey             = "RAPID_API_KEY"
	xAPIQuotaLeftHeaderName = "X-API-Quota-Left"
)

type humorAPIResponse struct {
	Content string `json:"joke"`
	ID      int64  `json:"id"`
}

var HumorAPILimitExceededErr = errors.New("humorAPI: daily limit of jokes was exceeded")

type HumorAPIService struct {
	client    http.Client
	active    bool
	m         sync.Mutex
	globalCtx context.Context
}

func (h *HumorAPIService) Active(ctx context.Context) (active bool) {
	select {
	case <-ctx.Done():
		active = false
	default:
		active = h.active
	}

	return
}

func (h *HumorAPIService) RandomJoke(ctx log.Context, search internal.SearchParams) (joke.DbModel, error) {
	select {
	case <-ctx.Done():
		return joke.DbModel{}, context.Canceled
	default:
	}

	if !h.Active(ctx) {
		return joke.DbModel{}, HumorAPILimitExceededErr
	}

	apiKey, ok := os.LookupEnv(humorAPIKey)
	if !ok {
		return joke.DbModel{}, errors.New("humorAPI: missing " + humorAPIKey)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, humorApiURL+humorAPICategory(search.Category), nil)
	if err != nil {
		return joke.DbModel{}, err
	}

	req.Header["X-RapidAPI-Key"] = []string{apiKey}
	req.Header["X-RapidAPI-Host"] = []string{"humor-jokes-and-memes.p.rapidapi.com"}

	res, err := h.client.Do(req)
	if err != nil {
		return joke.DbModel{}, err
	}
	defer res.Body.Close()

	limitExceeded := len(res.Header[xAPIQuotaLeftHeaderName]) > 0 && res.Header[xAPIQuotaLeftHeaderName][0] == "0"

	if res.StatusCode == http.StatusTooManyRequests || res.StatusCode == http.StatusPaymentRequired || limitExceeded {
		h.m.Lock()
		h.active = false
		resetTime := time.Now().Add(time.Hour * 24)

		go unlockService(h.globalCtx, &h.m, &h.active, resetTime)
		h.m.Unlock()

		return joke.DbModel{}, HumorAPILimitExceededErr
	} else if res.StatusCode != http.StatusOK {
		msg, err := io.ReadAll(res.Body)
		if err != nil {
			ctx.Logger.ErrorContext(ctx, "humorAPI: failed read response body", "error", err)
			msg = []byte{}
		}

		return joke.DbModel{}, fmt.Errorf("humorAPI: failed to get joke. status '%d', msg: '%s'", res.StatusCode, msg)
	}

	select {
	case <-ctx.Done():
		return joke.DbModel{}, context.Canceled
	default:
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return joke.DbModel{}, err
	}

	var humorRes humorAPIResponse
	if err = json.Unmarshal(resBody, &humorRes); err != nil {
		return joke.DbModel{}, err
	}

	return joke.DbModel{
		Category: search.Category,
		Type:     joke.Single,
		Answer:   humorRes.Content,
	}, nil
}

func humorAPICategory(category joke.Category) (res string) {
	switch category {
	case joke.PROGRAMMING:
		res = "nerdy"
	case joke.DARK:
		res = "dark"
	case joke.YOMAMA:
		res = "yo_mama"
	case joke.Any:
	default:
		res = "one_liner"
	}

	return
}
