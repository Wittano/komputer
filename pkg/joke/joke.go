package joke

import "github.com/wittano/komputer/internal"

type JokeMapper interface {
	ToJoke() (internal.Joke, error)
}
