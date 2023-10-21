package internal

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/wittano/komputer/internal/interaction"
	"github.com/wittano/komputer/internal/types"
	"math/rand"
)

type messageComponentHandler func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate)
type jokeType uint8

const singleType jokeType = 1

const (
	PleaseButtonId   = "apologiesButtonId"
	SameJokeButtonId = "sameJokeButtonId"
	NextJokeButtonId = "nextJokeButtonId"
)

var JokeMessageComponentHandler = map[string]messageComponentHandler{
	PleaseButtonId:   apologiseMe,
	SameJokeButtonId: nextJokeWithSameCategory,
	NextJokeButtonId: nextJoke,
}

func nextJokeWithSameCategory(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	embedFields := i.Message.Embeds[0].Fields
	c := types.JokeCategory(embedFields[len(embedFields)-1].Value)

	var t types.JokeType
	if len(embedFields) == int(singleType) {
		t = types.Single
	} else {
		t = types.TwoPart
	}

	interaction.SendJoke(ctx, s, i, t, c)
}

func nextJoke(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	var t types.JokeType
	if rand.Int()%2 == 0 {
		t = types.Single
	} else {
		t = types.TwoPart
	}

	c := types.GetRandomCategory()

	interaction.SendJoke(ctx, s, i, t, c)
}

func apologiseMe(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	interaction.CreateDiscordInteractionResponse(ctx, i, s, interaction.CreateDiscordMsg("Przepraszam"))
}
