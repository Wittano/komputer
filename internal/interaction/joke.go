package interaction

import (
	"context"
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/wittano/komputer/internal/joke"
	"github.com/wittano/komputer/internal/log"
	"github.com/wittano/komputer/internal/mongo"
	"github.com/wittano/komputer/internal/types"
	"math/rand"
)

type jokeSingleTypeGeneratorFunc func(ctx context.Context, category types.JokeCategory) (types.JokeContainer, error)
type jokeTwoPartGeneratorFunc func(ctx context.Context, category types.JokeCategory) (types.JokeTwoPartsContainer, error)

var (
	jokeSingleTypeGenerator = []jokeSingleTypeGeneratorFunc{
		joke.GetSingleJokeFromJokeDev,
		mongo.GetSingleTypeJoke,
	}

	jokeTwoPartsTypeGenerator = []jokeTwoPartGeneratorFunc{
		joke.GetTwoPartJokeFromJokeDev,
		mongo.GetTwoPartsTypeJoke,
	}
)

func SendJoke(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, t types.JokeType, c types.JokeCategory) {
	var msg *discordgo.InteractionResponseData

	switch t {
	case types.Single:
		j, err := getSingleTypeJokeGenerator()(ctx, c)
		if err != nil {
			log.Error(ctx, "Failed during getting single joke from JokeDev", err)

			if errors.Is(err, types.ErrJokeNotFound{Category: c, JokeType: t}) || errors.Is(err, joke.ErrJokeCategoryNotSupported{}) {
				CreateDiscordInteractionResponse(ctx, i, s, CreateJokeNotFoundMsg(t, c))
			} else {
				CreateDiscordInteractionResponse(ctx, i, s, CreateErrorMsg())
			}

			return
		}

		msg = CreateJokeMessage(i.Member.User.Username, c, j)
	case types.TwoPart:
		j, err := getTwoPartsTypeJokeGenerator()(ctx, c)
		if err != nil {
			log.Error(ctx, "Failed during getting two-part joke from JokeDev", err)

			if errors.Is(err, types.ErrJokeNotFound{Category: c, JokeType: t}) || errors.Is(err, joke.ErrJokeCategoryNotSupported{}) {
				CreateDiscordInteractionResponse(ctx, i, s, CreateJokeNotFoundMsg(t, c))
			} else {
				CreateDiscordInteractionResponse(ctx, i, s, CreateErrorMsg())
			}

			return
		}

		msg = CreateTwoPartJokeMessage(i.Member.User.Username, c, j)
	}

	CreateDiscordInteractionResponse(ctx, i, s, msg)
}

func getSingleTypeJokeGenerator() jokeSingleTypeGeneratorFunc {
	return jokeSingleTypeGenerator[rand.Int()%len(jokeSingleTypeGenerator)]
}

func getTwoPartsTypeJokeGenerator() jokeTwoPartGeneratorFunc {
	return jokeTwoPartsTypeGenerator[rand.Int()%len(jokeTwoPartsTypeGenerator)]
}
