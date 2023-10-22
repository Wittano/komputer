package joke

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wittano/komputer/internal/log"
	"github.com/wittano/komputer/internal/types"
	"io"
	"net/http"
	"net/url"
	"os"
)

type HumorAPIRes struct {
	JokeRes string `json:"joke"`
	ID      int64  `json:"id"`
}

func (h HumorAPIRes) Content() string {
	return h.JokeRes
}

type HumorAPILimitExceededErr struct {
}

type HumorAPIKeyMissingErr struct {
}

func (h HumorAPIKeyMissingErr) Error() string {
	return "RAPID_API_KEY is missing"
}

func (h HumorAPILimitExceededErr) Error() string {
	return "Humor API limit exceeded"
}

func GetRandomJokeFromHumorAPI(ctx context.Context, category types.JokeCategory) (HumorAPIRes, error) {
	key, ok := os.LookupEnv("RAPID_API_KEY")
	if !ok {
		return HumorAPIRes{}, HumorAPIKeyMissingErr{}
	}

	u, err := url.Parse("https://humor-jokes-and-memes.p.rapidapi.com/jokes/random?exclude-tags=nsfw&include-tags=" + category.ToHumorAPICategory())
	if err != nil {
		return HumorAPIRes{}, err
	}

	req := http.Request{
		Method: http.MethodGet,
		URL:    u,
		Header: map[string][]string{
			"X-RapidAPI-Key":  {key},
			"X-RapidAPI-Host": {"humor-jokes-and-memes.p.rapidapi.com"},
		},
	}

	res, err := client.Do(&req)
	if err != nil {
		return HumorAPIRes{}, err
	}

	defer res.Body.Close()
	if res.StatusCode == 429 {
		return HumorAPIRes{}, HumorAPILimitExceededErr{}
	} else if res.StatusCode != 200 {
		msg, err := io.ReadAll(res.Body)
		if err != nil {
			log.Error(ctx, "Failed read response body from HumorAPI", err)
			msg = []byte{}
		}

		return HumorAPIRes{}, errors.New(fmt.Sprintf("Failed to get joke from HumorAPI. Status %d, Msg: %s", res.StatusCode, msg))
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return HumorAPIRes{}, err
	}

	var joke HumorAPIRes

	err = json.Unmarshal(resBody, &joke)
	if err != nil {
		return HumorAPIRes{}, err
	}

	return joke, nil
}
