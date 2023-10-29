package command

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/wittano/komputer/internal/interaction"
	"github.com/wittano/komputer/internal/log"
	"github.com/wittano/komputer/internal/mongo"
	"github.com/wittano/komputer/internal/types"
	"os"
)

var (
	JokeCommand = DiscordCommand{
		Command: discordgo.ApplicationCommand{
			Name:        "joke",
			Description: "Tell me some joke",
			GuildID:     os.Getenv("SERVER_GUID"),
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				getJokeCategoryOption(false),
				getJokeTypeOption(false),
			},
		},
		Execute: executeJokeCommand,
	}

	AddJokeCommand = DiscordCommand{
		Command: discordgo.ApplicationCommand{
			Name:        "add-joke",
			Description: "Add new joke to server database",
			GuildID:     os.Getenv("SERVER_GUID"),
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				getJokeCategoryOption(true),
				getJokeTypeOption(true),
				{
					Name:        "content",
					Description: "Main part of joke",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
				{
					Name:        "question",
					Description: "Question part in two-part joke",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    false,
				},
			},
		},
		Execute: executeAddJokeCommand,
	}
)

func executeAddJokeCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	j := mongo.JokeDB{
		GuildID: i.GuildID,
	}

	for _, o := range i.Data.(discordgo.ApplicationCommandInteractionData).Options {
		switch o.Name {
		case "category":
			j.Category = types.JokeCategory(o.Value.(string))
		case "type":
			j.Type = types.JokeType(o.Value.(string))
		case "content":
			j.ContentRes = o.Value.(string)
		case "question":
			j.Question = o.Value.(string)
		default:
			log.Warn(ctx, fmt.Sprintf("Invalid option for %s", o.Name))
		}
	}

	if j.Category == "" {
		interaction.CreateDiscordInteractionResponse(ctx, i, s, interaction.CreateDiscordMsg("BEEP BOOP. Brakuje kategori!"))
		return
	}

	if j.Type == "" {
		interaction.CreateDiscordInteractionResponse(ctx, i, s, interaction.CreateDiscordMsg("BEEP BOOP. Brakuje typu żartu!"))
		return
	}

	if j.ContentRes == "" {
		interaction.CreateDiscordInteractionResponse(ctx, i, s, interaction.CreateDiscordMsg("BEEP BOOP. Gdzie jest żart panie Kapitanie!"))
		return
	}

	id, err := mongo.AddJoke(ctx, j)
	if err != nil {
		log.Error(ctx, "Failed add new joke into database", err)
		interaction.CreateDiscordInteractionResponse(ctx, i, s, interaction.CreateDiscordMsg("BEEP BOOP. Coś poszło nie tak z dodanie twego żartu Kapitanie"))
		return
	}

	msg := &discordgo.InteractionResponseData{
		Flags:   discordgo.MessageFlagsEphemeral,
		Content: fmt.Sprintf("BEEP BOOP. Dodałem twój żart panie Kapitanie. Jego ID to %s", id.Hex()),
	}

	interaction.CreateDiscordInteractionResponse(ctx, i, s, msg)
}

func executeJokeCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	category := types.ANY
	jokeType := types.Single

	for _, o := range i.Data.(discordgo.ApplicationCommandInteractionData).Options {
		switch o.Name {
		case "category":
			category = types.JokeCategory(o.Value.(string))
		case "type":
			jokeType = types.JokeType(o.Value.(string))
		default:
			log.Warn(ctx, fmt.Sprintf("Invalid option for %s", o.Name))
		}
	}

	interaction.SendJoke(ctx, s, i, jokeType, category)
}
