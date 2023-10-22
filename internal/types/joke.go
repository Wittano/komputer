package types

import (
	"fmt"
	"math/rand"
)

type (
	JokeType     string
	JokeCategory string
)

func (j JokeCategory) ToJokeDevCategory() JokeCategory {
	if j == YOMAMA {
		return ANY
	}

	return j
}

func (j JokeCategory) ToHumorAPICategory() string {
	switch j {
	case PROGRAMMING:
		return "nerdy"
	case ANY:
		return "one_liner"
	case DARK:
		return "dark"
	case YOMAMA:
		return "yo_mama"
	default:
		return "one_liner"
	}
}

type JokeContainer interface {
	Content() string
}

type JokeTwoPartsContainer interface {
	ContentTwoPart() (string, string)
}

type JokeCategoryContainer interface {
	Category() JokeCategory
}

type ErrJokeNotFound struct {
	Category JokeCategory
	JokeType JokeType
}

func (j ErrJokeNotFound) Error() string {
	return fmt.Sprintf("JokeContainer \"%s\" from category \"%s\" wasn't found", j.JokeType, j.Category)
}

const (
	Single  JokeType = "single"
	TwoPart JokeType = "twopart"
)

const (
	PROGRAMMING JokeCategory = "Programming"
	MISC        JokeCategory = "Misc"
	DARK        JokeCategory = "Dark"
	YOMAMA      JokeCategory = "YoMama"
	ANY         JokeCategory = "Any"
)

func GetRandomCategory() JokeCategory {
	c := []JokeCategory{PROGRAMMING, MISC, DARK, ANY}

	return c[rand.Int()%len(c)]
}
