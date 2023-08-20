package internal

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"github.com/wittano/komputer/pkg/joke"
)

func CreateJokeMessage(ctx context.Context, username string, joke joke.Joke) *discordgo.InteractionResponseData {
	content, err := joke.Content()
	if err != nil {
		return createErrorMsg(ctx, err)
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

func CreateTwoPartJokeMessage(ctx context.Context, username string, joke joke.JokeTwoParts) *discordgo.InteractionResponseData {
	question, answer, err := joke.ContentTwoPart()
	if err != nil {
		return createErrorMsg(ctx, err)
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

func createErrorMsg(ctx context.Context, err error) *discordgo.InteractionResponseData {
	log.Err(err).Str("traceID", ctx.Value("traceID").(string)).Msg("Failed to send message!")

	return &discordgo.InteractionResponseData{Content: fmt.Sprintf("BEEP BOOM. Something went wrong :(")}
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
