package types

type (
	JokeType     string
	JokeCategory string // TODO Add yoMama category
)

const (
	Single  JokeType = "single"
	TwoPart JokeType = "twopart"
)

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

type Joke interface {
	Content() string
}

type JokeTwoParts interface {
	ContentTwoPart() (string, string)
}
