package jokedev

import (
	"encoding/json"
	"fmt"
	"github.com/wittano/komputer/internal"
	"io"
	"log"
	"math/rand"
)

type JokeApiDev struct {
	category JokeType
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

type (
	JokeStructureType string
	JokeType          string
)

const (
	Single  JokeStructureType = "single"
	TwoPart JokeStructureType = "twopart"
)

const (
	PROGRAMMING JokeType = "Programming"
	MISC        JokeType = "Misc"
	DARK        JokeType = "Dark"
	PUN         JokeType = "Pun"
	SPOOKY      JokeType = "Spooky"
	CHRISTMAS   JokeType = "Christmas"
	ANY         JokeType = "Any"
)

func New(category JokeType) JokeApiDev {
	return JokeApiDev{
		category: category,
	}
}

func (j JokeApiDev) Category() string {
	return string(j.category)
}

func (j JokeApiDev) Content() string {
	specialJoke := yoMommaJoke()

	var joke jokeApiSingleResponse

	if specialJoke != "" {
		joke = jokeApiSingleResponse{
			Joke:     specialJoke,
			Category: "YoMomma",
		}
	}

	joke, err := j.getSingleJoke()
	if err != nil {
		log.Print(err)
		return ""
	}

	j.category = JokeType(joke.Category)

	return joke.Joke
}

func (j JokeApiDev) ContentTwoPart() (string, string) {
	joke, err := j.getTwoPartJoke()
	if err != nil {
		log.Print(err)
		return "", ""
	}

	j.category = JokeType(joke.Category)

	return joke.Setup, joke.Delivery
}

func (j JokeApiDev) getSingleJoke() (joke jokeApiSingleResponse, err error) {
	res, err := internal.Client.Get(fmt.Sprintf("https://v2.jokeapi.dev/joke/%s?type=%s", j.category, Single))
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

func (j JokeApiDev) getTwoPartJoke() (jokeApiTwoPartResponse, error) {
	res, err := internal.Client.Get(fmt.Sprintf("https://v2.jokeapi.dev/joke/%s?type=%s", j.category, TwoPart))
	if err != nil {
		return jokeApiTwoPartResponse{}, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return jokeApiTwoPartResponse{}, err
	}

	var jokeRes jokeApiTwoPartResponse
	err = json.Unmarshal(resBody, &jokeRes)

	return jokeRes, err
}

func yoMommaJoke() string {
	if rand.Int()%50 == 0 {
		return "Yo Momma so fat32!"
	}

	return ""
}
