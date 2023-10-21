package internal

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/wittano/komputer/internal/log"
	"github.com/wittano/komputer/internal/types"
)

func CreateJokeMessage(username string, category types.JokeCategory, joke types.Joke) *discordgo.InteractionResponseData {
	content := joke.Content()

	return &discordgo.InteractionResponseData{
		Content:    fmt.Sprintf("BEEP BOOP, Tak jest kapitanie %s!", username),
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
						Value: string(category),
					},
				},
			},
		},
	}
}

func CreateTwoPartJokeMessage(username string, category types.JokeCategory, joke types.JokeTwoParts) *discordgo.InteractionResponseData {
	question, answer := joke.ContentTwoPart()

	return &discordgo.InteractionResponseData{
		Content:    fmt.Sprintf("BEEP BOOP, Tak jest Panie kapitanie %s!", username),
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
						Value: string(category),
					},
				},
			},
		},
	}
}

func CreateErrorMsg() *discordgo.InteractionResponseData {
	return CreateDiscordMsg("BEEP BOOP. Coś poszło nie tak :(")
}

func CreateDiscordMsg(msg string) *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{Content: msg}
}

func createButtonReactions() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Style:    discordgo.PrimaryButton,
					Label:    "Przeproś",
					CustomID: pleaseButtonId,
				},
				discordgo.Button{
					Style:    discordgo.SecondaryButton,
					Label:    "Zabawne powiedz więcej",
					CustomID: jokeButtonId,
				},
			},
		},
	}
}

func CreateDiscordInteractionResponse(ctx context.Context, i *discordgo.InteractionCreate, s *discordgo.Session, msg *discordgo.InteractionResponseData) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: msg,
	}); err != nil {
		log.Error(ctx, "Failed send response to discord user", err)
	}
}
