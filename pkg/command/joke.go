package command

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/wittano/komputer/internal"
	"github.com/wittano/komputer/internal/log"
	"github.com/wittano/komputer/pkg/joke/jokedev"
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
	category := jokedev.ANY
	jokeType := jokedev.Single

	for _, o := range i.Data.(discordgo.ApplicationCommandInteractionData).Options {
		switch o.Name {
		case "category":
			category = jokedev.JokeType(o.Value.(string))
		case "type":
			jokeType = jokedev.JokeStructureType(o.Value.(string))
		default:
			log.Warn(ctx, fmt.Sprintf("Invalid option for %s", o.Name))
		}
	}

	joke := jokedev.New(ctx, category)

	var msg *discordgo.InteractionResponseData

	switch jokeType {
	case jokedev.Single:
		msg = internal.CreateJokeMessage(ctx, i.Member.User.Username, joke)
	case jokedev.TwoPart:
		msg = internal.CreateTwoPartJokeMessage(ctx, i.Member.User.Username, joke)
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: msg,
	})

	if err != nil {
		log.Error(ctx, "Failed to send response", err)
	}
}
