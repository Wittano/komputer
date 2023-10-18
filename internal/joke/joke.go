package joke

type (
	JokeStructureType string
	JokeType          string
)

type Joke interface {
	Content() string
}

type JokeTwoParts interface {
	ContentTwoPart() (string, string)
}
