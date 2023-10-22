package joke

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wittano/komputer/internal/types"
	"io"
	"regexp"
	"strings"
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
	Error       bool         `json:"error"`
	CategoryRes string       `json:"category"`
	Type        string       `json:"type"`
	Flags       jokeApiFlags `json:"flags"`
	Id          int          `json:"id"`
	Safe        bool         `json:"safe"`
	Lang        string       `json:"lang"`
	Joke        string       `json:"joke"`
}

type jokeApiTwoPartResponse struct {
	Error       bool         `json:"error"`
	CategoryRes string       `json:"category"`
	Type        string       `json:"type"`
	Flags       jokeApiFlags `json:"flags"`
	Id          int          `json:"id"`
	Safe        bool         `json:"safe"`
	Lang        string       `json:"lang"`
	Setup       string       `json:"setup"`
	Delivery    string       `json:"delivery"`
}

type ErrJokeCategoryNotSupported struct{}

func (e ErrJokeCategoryNotSupported) Error() string {
	return "category isn't support by JokeDevAPI"
}

func (j jokeApiSingleResponse) Content() string {
	return j.Joke
}

var YoMamaRegex = regexp.MustCompile("(mama)|(momma)|(mother)")

func (j jokeApiSingleResponse) Category() types.JokeCategory {
	if YoMamaRegex.MatchString(j.Joke) {
		return types.YOMAMA
	}

	return types.JokeCategory(j.CategoryRes)
}

func (j jokeApiTwoPartResponse) ContentTwoPart() (string, string) {
	return j.Setup, j.Delivery
}

func (j jokeApiTwoPartResponse) Category() types.JokeCategory {
	quest := strings.ToLower(j.Delivery)
	content := strings.ToLower(j.Setup)

	if YoMamaRegex.MatchString(quest) || YoMamaRegex.MatchString(content) {
		return types.YOMAMA
	}

	return types.JokeCategory(j.CategoryRes)
}

func GetSingleJokeFromJokeDev(_ context.Context, category types.JokeCategory) (joke types.JokeContainer, err error) {
	if category == types.YOMAMA {
		return nil, ErrJokeCategoryNotSupported{}
	}

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

func GetTwoPartJokeFromJokeDev(_ context.Context, category types.JokeCategory) (types.JokeTwoPartsContainer, error) {
	if category == types.YOMAMA {
		return nil, ErrJokeCategoryNotSupported{}
	}

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
