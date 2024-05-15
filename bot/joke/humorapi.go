package joke

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wittano/komputer/bot/log"
	"io"
	"log/slog"
	"net/http"
	"os"
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
	globalCtx context.Context
}

func (h HumorAPIService) Active(ctx context.Context) (active bool) {
	select {
	case <-ctx.Done():
		active = false
	default:
		active = h.active
	}

	return
}

func (h *HumorAPIService) Joke(ctx context.Context, search SearchParams) (Joke, error) {
	select {
	case <-ctx.Done():
		return Joke{}, context.Canceled
	default:
	}

	if !h.Active(ctx) {
		return Joke{}, HumorAPILimitExceededErr
	}

	apiKey, ok := os.LookupEnv(humorAPIKey)
	if !ok {
		return Joke{}, errors.New("humorAPI: missing " + humorAPIKey)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, humorApiURL+humorAPICategory(search.Category), nil)
	if err != nil {
		return Joke{}, err
	}

	req.Header["X-RapidAPI-Key"] = []string{apiKey}
	req.Header["X-RapidAPI-Host"] = []string{"humor-jokes-and-memes.p.rapidapi.com"}

	res, err := h.client.Do(req)
	if err != nil {
		return Joke{}, err
	}
	defer res.Body.Close()

	limitExceeded := len(res.Header[xAPIQuotaLeftHeaderName]) > 0 && res.Header[xAPIQuotaLeftHeaderName][0] == "0"

	if res.StatusCode == http.StatusTooManyRequests || res.StatusCode == http.StatusPaymentRequired || limitExceeded {
		h.active = false
		resetTime := time.Now().Add(time.Hour * 24)

		go unlockService(h.globalCtx, &h.active, resetTime)

		return Joke{}, HumorAPILimitExceededErr
	} else if res.StatusCode != http.StatusOK {
		msg, err := io.ReadAll(res.Body)
		if err != nil {
			log.Log(ctx, func(l slog.Logger) {
				l.ErrorContext(ctx, "humorAPI: failed read response body", "error", err)
			})
			msg = []byte{}
		}

		return Joke{}, fmt.Errorf("humorAPI: failed to get joke. status '%d', msg: '%s'", res.StatusCode, msg)
	}

	select {
	case <-ctx.Done():
		return Joke{}, context.Canceled
	default:
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return Joke{}, err
	}

	var joke humorAPIResponse
	if err = json.Unmarshal(resBody, &joke); err != nil {
		return Joke{}, err
	}

	return Joke{
		Category: search.Category,
		Type:     Single,
		Answer:   joke.Content,
	}, nil
}

func humorAPICategory(category Category) (res string) {
	switch category {
	case PROGRAMMING:
		res = "nerdy"
	case DARK:
		res = "dark"
	case YOMAMA:
		res = "yo_mama"
	case Any:
	default:
		res = "one_liner"
	}

	return
}
