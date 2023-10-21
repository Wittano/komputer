package types

import (
	"fmt"
	"math/rand"
)

type (
	JokeType     string
	JokeCategory string
)

type Joke interface {
	Content() string
}

type JokeTwoParts interface {
	ContentTwoPart() (string, string)
}

type JokeNotFoundErr struct {
	Category JokeCategory
	JokeType JokeType
}

func (j JokeNotFoundErr) Error() string {
	return fmt.Sprintf("Joke \"%s\" from category \"%s\" wasn't found", j.JokeType, j.Category)
}

const (
	Single  JokeType = "single"
	TwoPart JokeType = "twopart"
)

// TODO Add yoMama category
const (
	PROGRAMMING JokeCategory = "Programming"
	MISC        JokeCategory = "Misc"
	DARK        JokeCategory = "Dark"
	ANY         JokeCategory = "Any"
)

func (j JokeCategory) ToHumorAPICategory() string {
	var c string

	switch j {
	case PROGRAMMING:
		c = "nerdy"
	case ANY:
		c = "one_liner"
	case DARK:
		c = "dark"
	default:
		c = "one_liner"
	}

	return c
}

func GetRandomCategory() JokeCategory {
	c := []JokeCategory{PROGRAMMING, MISC, DARK, ANY}

	return c[rand.Int()%len(c)]
}
