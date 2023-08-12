package command

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
)

var WelcomeCommand = DiscordCommand{
	Command: discordgo.ApplicationCommand{
		Name:        "welcome",
		Description: "Welcome command to greetings to you",
		GuildID:     os.Getenv("SERVER_GUID"),
		Type:        discordgo.ChatApplicationCommand,
	},
	Execute: execute,
}

func execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Welcome %s! How can I help you?", i.Member.User.Username),
		},
	})

	if err != nil {
		log.Print(err)
	}
}
