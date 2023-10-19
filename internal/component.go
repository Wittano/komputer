package internal

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/wittano/komputer/internal/joke"
	"github.com/wittano/komputer/internal/log"
)

type messageComponentHandler func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate)

type jokeType uint8

const (
	singleType  jokeType = 1
	twoPartType jokeType = 3
)

const (
	pleaseButtonId = "ApologiesButtonId"
	jokeButtonId   = "jokeButtonId"
)

var (
	JokeMessageComponentHandler = map[string]messageComponentHandler{
		pleaseButtonId: apologiseMe,
		jokeButtonId:   nextJoke,
	}
)

func nextJoke(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	embedFields := i.Message.Embeds[0].Fields
	category := joke.JokeCategory(embedFields[len(embedFields)-1].Value)

	var msg *discordgo.InteractionResponseData

	switch jokeType(len(embedFields)) {
	case singleType:
		j, err := joke.GetSingleJokeFromJokeDev(category)
		if err != nil {
			log.Error(ctx, "Failed during getting single joke from JokeDev", err)

			CreateErrorMsg()

			return
		}

		msg = CreateJokeMessage(i.Member.User.Username, category, j)
	case twoPartType:
		j, err := joke.GetTwoPartJokeFromJokeDev(category)
		if err != nil {
			log.Error(ctx, "Failed during getting two-part joke from JokeDev", err)

			CreateErrorMsg()

			return
		}

		msg = CreateTwoPartJokeMessage(i.Member.User.Username, category, j)
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: msg,
	})

	if err != nil {
		log.Error(ctx, "Failed create Discord interaction response", err)

		CreateErrorMsg()
	}
}

func apologiseMe(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	sendMessage(ctx, "Przepraszam", s, i)
}

func sendMessage(ctx context.Context, msg string, s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})

	if err != nil {
		log.Error(ctx, "Failed create Discord interaction response", err)

		CreateErrorMsg()
	}
}
