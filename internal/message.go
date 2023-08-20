package internal

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/wittano/komputer/pkg/joke"
	"log"
)

type messageComponentHandler func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate)

const (
	PleaseButtonId = "ApologiesButtonId"
	FunnyButtonId  = "FunnyButtonId"
)

var (
	JokeMessageComponentHandler = map[string]messageComponentHandler{
		PleaseButtonId: apologiseMe,
		FunnyButtonId:  funnyMe,
	}
)

func funnyMe(_ context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	sendMessage("Funny", s, i)
}

func apologiseMe(_ context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	sendMessage("Przepraszam", s, i)
}

func sendMessage(msg string, s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})

	if err != nil {
		log.Print(err)
	}
}

func CreateJokeMessage(username string, joke joke.Joke) *discordgo.InteractionResponseData {
	content, err := joke.Content()
	if err != nil {
		return createErrorMsg()
	}

	return &discordgo.InteractionResponseData{
		Content:    fmt.Sprintf("BEEP BOOP, Yes my captin %s!", username),
		Components: createButtonReactions(),
		Embeds: []*discordgo.MessageEmbed{
			{
				Type:        discordgo.EmbedTypeRich,
				Title:       "Joke",
				Description: content,
				Color:       0x02f5f5,
				Author: &discordgo.MessageEmbedAuthor{
					Name: "komputer",
				},
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "category",
						Value: joke.Category(),
					},
				},
			},
		},
	}
}

func CreateTwoPartJokeMessage(username string, joke joke.JokeTwoParts) *discordgo.InteractionResponseData {
	question, answer, err := joke.ContentTwoPart()
	if err != nil {
		return createErrorMsg()
	}

	return &discordgo.InteractionResponseData{
		Content:    fmt.Sprintf("BEEP BOOP, Yes my captin %s!", username),
		Components: createButtonReactions(),
		Embeds: []*discordgo.MessageEmbed{
			{
				Type:  discordgo.EmbedTypeRich,
				Title: "Joke",
				Color: 0x02f5f5,
				Author: &discordgo.MessageEmbedAuthor{
					Name: "komputer",
				},
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Question",
						Value:  question,
						Inline: true,
					},
					{
						Name:   "Answer",
						Value:  answer,
						Inline: true,
					},
					{
						Name:  "Category",
						Value: joke.Category(),
					},
				},
			},
		},
	}
}

func createErrorMsg() *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{Content: fmt.Sprintf("BEEP BOOM. Something went wrong :(")}
}

func createButtonReactions() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Style:    discordgo.PrimaryButton,
					Label:    "Przeproś",
					CustomID: PleaseButtonId,
				},
				discordgo.Button{
					Style:    discordgo.SecondaryButton,
					Label:    "Funny",
					CustomID: FunnyButtonId,
				},
			},
		},
	}
}
