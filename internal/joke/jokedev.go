package joke

import (
	"encoding/json"
	"fmt"
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

const (
	Single  JokeType = "single"
	TwoPart JokeType = "twopart"
)

const (
	PROGRAMMING JokeCategory = "Programming"
	MISC        JokeCategory = "Misc"
	DARK        JokeCategory = "Dark"
	PUN         JokeCategory = "Pun"
	SPOOKY      JokeCategory = "Spooky"
	CHRISTMAS   JokeCategory = "Christmas"
	ANY         JokeCategory = "Any"
)

func (j jokeApiSingleResponse) Content() string {
	return j.Joke
}

func (j jokeApiTwoPartResponse) ContentTwoPart() (string, string) {
	return j.Setup, j.Delivery
}

func GetSingleJokeFromJokeDev(category JokeCategory) (joke Joke, err error) {
	res, err := Client.Get(fmt.Sprintf("https://v2.jokeapi.dev/joke/%s?type=%s", category, Single))
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

func GetTwoPartJokeFromJokeDev(category JokeCategory) (joke JokeTwoParts, err error) {
	res, err := Client.Get(fmt.Sprintf("https://v2.jokeapi.dev/joke/%s?type=%s", category, TwoPart))
	if err != nil {
		return
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(resBody, &joke)

	return joke, err
}
