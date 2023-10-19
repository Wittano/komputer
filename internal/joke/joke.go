package joke

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
