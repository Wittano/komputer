package command

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/wittano/komputer/internal"
	"github.com/wittano/komputer/internal/joke"
	"github.com/wittano/komputer/internal/log"
	"github.com/wittano/komputer/internal/mongo"
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
	var j mongo.JokeDb

	for _, o := range i.Data.(discordgo.ApplicationCommandInteractionData).Options {
		switch o.Name {
		case "category":
			j.Category = joke.JokeCategory(o.Value.(string))
		case "type":
			j.Type = joke.JokeType(o.Value.(string))
		case "content":
			j.ContentRes = o.Value.(string)
		case "question":
			j.Question = o.Value.(string)
		default:
			log.Warn(ctx, fmt.Sprintf("Invalid option for %s", o.Name))
		}
	}

	if j.Category == "" {
		internal.CreateDisacordInteractionResponse(ctx, i, s, internal.CreateDiscordMsg("BEEP BOOP. Brakuje kategori!"))
		return
	}

	if j.Type == "" {
		internal.CreateDisacordInteractionResponse(ctx, i, s, internal.CreateDiscordMsg("BEEP BOOP. Brakuje typu żartu!"))
		return
	}

	if j.ContentRes == "" {
		internal.CreateDisacordInteractionResponse(ctx, i, s, internal.CreateDiscordMsg("BEEP BOOP. Gdzie jest żart panie Kapitanie!"))
		return
	}

	id, err := mongo.AddJoke(ctx, j)
	if err != nil {
		log.Error(ctx, "Failed add new joke into database", err)
		internal.CreateDisacordInteractionResponse(ctx, i, s, internal.CreateDiscordMsg("BEEP BOOP. Coś poszło nie tak z dodanie twego żartu Kapitanie"))
		return
	}

	// TODO Add only-user show this message
	internal.CreateDisacordInteractionResponse(ctx, i, s, internal.CreateDiscordMsg(fmt.Sprintf("BEEP BOOP. Dodałem twój żart panie Kapitanie. Jego ID to %s", id.Hex())))
}

func executeJokeCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	category := joke.ANY
	jokeType := joke.Single

	for _, o := range i.Data.(discordgo.ApplicationCommandInteractionData).Options {
		switch o.Name {
		case "category":
			category = joke.JokeCategory(o.Value.(string))
		case "type":
			jokeType = joke.JokeType(o.Value.(string))
		default:
			log.Warn(ctx, fmt.Sprintf("Invalid option for %s", o.Name))
		}
	}

	var msg *discordgo.InteractionResponseData

	switch jokeType {
	case joke.Single:
		j, err := joke.GetSingleJokeFromJokeDev(category)
		if err != nil {
			log.Error(ctx, "Failed during single j from JokeDev", err)

			internal.CreateDisacordInteractionResponse(ctx, i, s, internal.CreateErrorMsg())
			return
		}

		msg = internal.CreateJokeMessage(i.Member.User.Username, category, j)
	case joke.TwoPart:
		j, err := joke.GetTwoPartJokeFromJokeDev(category)
		if err != nil {
			log.Error(ctx, "Failed during two-part j from JokeDev", err)

			internal.CreateDisacordInteractionResponse(ctx, i, s, internal.CreateErrorMsg())
			return
		}

		msg = internal.CreateTwoPartJokeMessage(i.Member.User.Username, category, j)
	}

	internal.CreateDisacordInteractionResponse(ctx, i, s, msg)
}
