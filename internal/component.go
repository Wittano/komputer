package internal

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/wittano/komputer/pkg/joke/jokedev"
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
	joke := jokedev.New(ctx, jokedev.JokeType(embedFields[len(embedFields)-1].Value))

	var msg *discordgo.InteractionResponseData

	switch jokeType(len(embedFields)) {
	case singleType:
		msg = CreateJokeMessage(ctx, i.Member.User.Username, joke)
	case twoPartType:
		msg = CreateTwoPartJokeMessage(ctx, i.Member.User.Username, joke)
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: msg,
	})

	if err != nil {
		CreateErrorMsg(ctx, err)
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
		CreateErrorMsg(ctx, err)
	}
}
