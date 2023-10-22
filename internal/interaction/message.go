package interaction

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/wittano/komputer/internal/log"
	"github.com/wittano/komputer/internal/types"
	"strings"
)

func CreateJokeMessage(username string, category types.JokeCategory, joke types.JokeContainer) *discordgo.InteractionResponseData {
	var c types.JokeCategory
	if jc, ok := joke.(types.JokeCategoryContainer); ok {
		c = jc.Category()
	} else {
		c = category
	}

	embeds := []*discordgo.MessageEmbed{
		{
			Type:        discordgo.EmbedTypeRich,
			Title:       "JokeContainer",
			Description: joke.Content(),
			Color:       0x02f5f5,
			Author: &discordgo.MessageEmbedAuthor{
				Name: "komputer",
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "category",
					Value: string(c),
				},
			},
		},
	}

	if c == types.YOMAMA {
		embeds[0].Image = &discordgo.MessageEmbedImage{
			URL: "https://media.tenor.com/sgS8GdoZGn8AAAAd/muscle-man-regular-show-muscle-man.gif",
		}
	}

	return &discordgo.InteractionResponseData{
		Content:    fmt.Sprintf("BEEP BOOP, Tak jest kapitanie %s!", username),
		Components: createButtonReactions(),
		Embeds:     embeds,
	}
}

func CreateTwoPartJokeMessage(username string, category types.JokeCategory, joke types.JokeTwoPartsContainer) *discordgo.InteractionResponseData {
	question, answer := joke.ContentTwoPart()

	var c types.JokeCategory
	if jc, ok := joke.(types.JokeCategoryContainer); ok {
		c = jc.Category()
	} else {
		c = category
	}

	embeds := []*discordgo.MessageEmbed{
		{
			Type:  discordgo.EmbedTypeRich,
			Title: "JokeContainer",
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
					Name:  "category",
					Value: string(c),
				},
			},
		},
	}

	if c == types.YOMAMA {
		embeds[0].Image = &discordgo.MessageEmbedImage{
			URL: "https://media.tenor.com/sgS8GdoZGn8AAAAd/muscle-man-regular-show-muscle-man.gif",
		}
	}

	return &discordgo.InteractionResponseData{
		Content:    fmt.Sprintf("BEEP BOOP, Tak jest Panie kapitanie %s!", username),
		Components: createButtonReactions(),
		Embeds:     embeds,
	}
}

func CreateErrorMsg() *discordgo.InteractionResponseData {
	return CreateDiscordMsg("BEEP BOOP. Coś poszło nie tak :(")
}

func CreateJokeNotFoundMsg(t types.JokeType, c types.JokeCategory) *discordgo.InteractionResponseData {
	return CreateDiscordMsg(fmt.Sprintf("BEEP BOOP. Nie udało mi się znaleść, żartu w kategori %s o typie %s", strings.ToUpper(string(t)), strings.ToUpper(string(c))))
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
					CustomID: "apologiesButtonId",
				},
				discordgo.Button{
					Style:    discordgo.SecondaryButton,
					Label:    "Zabawne powiedz coś podobnego",
					CustomID: "sameJokeButtonId",
				},
				discordgo.Button{
					Style:    discordgo.SecondaryButton,
					Label:    "Zabawne powiedz więcej",
					CustomID: "nextJokeButtonId",
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
