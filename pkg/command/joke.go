package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/wittano/komputer/internal"
	"github.com/wittano/komputer/pkg/joke/jokedev"
	"log"
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
						Name:  "Programming",
						Value: "Programming",
					},
					{
						Name:  "Misc",
						Value: "Misc",
					},
					{
						Name:  "Dark",
						Value: "Dark",
					},
					{
						Name:  "Pun",
						Value: "Pun",
					},
					{
						Name:  "Spooky",
						Value: "Spooky",
					},
					{
						Name:  "Christmas",
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

func executeJokeCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	category := jokedev.ANY
	jokeType := jokedev.Single

	for _, o := range i.Data.(discordgo.ApplicationCommandInteractionData).Options {
		switch o.Name {
		case "category":
			category = jokedev.JokeType(o.Value.(string))
		case "type":
			jokeType = jokedev.JokeStructureType(o.Value.(string))
		}
	}

	joke := jokedev.New(category)

	var msg *discordgo.InteractionResponseData

	switch jokeType {
	case jokedev.Single:
		msg = internal.CreateJokeMessage(i.Member.User.Username, joke)
	case jokedev.TwoPart:
		msg = internal.CreateTwoPartJokeMessage(i.Member.User.Username, joke)
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: msg,
	})

	if err != nil {
		log.Print(err)
	}
}
