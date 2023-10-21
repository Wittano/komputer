package joke

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wittano/komputer/internal/types"
	"io"
)

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
	Joke     string       `json:"joke"`
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

func (j jokeApiSingleResponse) Content() string {
	return j.Joke
}

func (j jokeApiTwoPartResponse) ContentTwoPart() (string, string) {
	return j.Setup, j.Delivery
}

func GetSingleJokeFromJokeDev(_ context.Context, category types.JokeCategory) (joke types.Joke, err error) {
	res, err := client.Get(fmt.Sprintf("https://v2.jokeapi.dev/joke/%s?type=%s", category, types.Single))
	if err != nil {
		return
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	var jokeRes jokeApiSingleResponse
	err = json.Unmarshal(resBody, &jokeRes)

	return jokeRes, err
}

func GetTwoPartJokeFromJokeDev(_ context.Context, category types.JokeCategory) (types.JokeTwoParts, error) {
	res, err := client.Get(fmt.Sprintf("https://v2.jokeapi.dev/joke/%s?type=%s", category, types.TwoPart))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var joke jokeApiTwoPartResponse
	err = json.Unmarshal(resBody, &joke)

	return joke, err
}
