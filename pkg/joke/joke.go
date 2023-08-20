package joke

type JokeCategory interface {
	Category() string
}

type Joke interface {
	JokeCategory
	Content() (string, error)
}

type JokeTwoParts interface {
	JokeCategory
	ContentTwoPart() (string, string, error)
}
