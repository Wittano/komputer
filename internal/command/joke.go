package command

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/wittano/komputer/internal"
	"github.com/wittano/komputer/internal/joke"
	"github.com/wittano/komputer/internal/log"
	"os"
)

var JokeCommand = DiscordCommand{
	Command: discordgo.ApplicationCommand{
		Name:        "joke",
		Description: "Tell me some joke",
		GuildID:     os.Getenv("SERVER_GUID"),
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "category",
				Description: "Joke category",
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Programowanie",
						Value: "Programming",
					},
					{
						Name:  "Różne",
						Value: "Misc",
					},
					{
						Name:  "Czarny humor",
						Value: "Dark",
					},
					{
						Name:  "Pun",
						Value: "Pun",
					},
					{
						Name:  "Straszne",
						Value: "Spooky",
					},
					{
						Name:  "Świątecznie",
						Value: "Christmas",
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "type",
				Description: "Type of joke",
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Single",
						Value: "single",
					},
					{

						Name:  "Two-Part",
						Value: "twopart",
					},
				},
			},
		},
	},
	Execute: executeJokeCommand,
}

func executeJokeCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	category := joke.ANY
	jokeType := joke.Single

	for _, o := range i.Data.(discordgo.ApplicationCommandInteractionData).Options {
		switch o.Name {
		case "category":
			category = joke.JokeType(o.Value.(string))
		case "type":
			jokeType = joke.JokeStructureType(o.Value.(string))
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

			internal.CreateErrorMsg()
			return
		}

		msg = internal.CreateJokeMessage(i.Member.User.Username, category, j)
	case joke.TwoPart:
		j, err := joke.GetTwoPartJokeFromJokeDev(category)
		if err != nil {
			log.Error(ctx, "Failed during two-part j from JokeDev", err)

			internal.CreateErrorMsg()
			return
		}

		msg = internal.CreateTwoPartJokeMessage(i.Member.User.Username, category, j)
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: msg,
	})

	if err != nil {
		log.Error(ctx, "Failed to send response", err)
	}
}
