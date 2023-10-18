package internal

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/wittano/komputer/internal/joke"
)

func CreateJokeMessage(username string, category joke.JokeType, joke joke.Joke) *discordgo.InteractionResponseData {
	content := joke.Content()

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
						Value: string(category),
					},
				},
			},
		},
	}
}

func CreateTwoPartJokeMessage(username string, category joke.JokeType, joke joke.JokeTwoParts) *discordgo.InteractionResponseData {
	question, answer := joke.ContentTwoPart()

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
						Value: string(category),
					},
				},
			},
		},
	}
}

func CreateErrorMsg() *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{Content: fmt.Sprintf("BEEP BOOP. Coś poszło nie tak :(")}
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
