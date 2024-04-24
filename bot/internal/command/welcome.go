package command

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"strings"
)

type WelcomeCommand struct{}

const WelcomeCommandName = "welcome"

func (w WelcomeCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        WelcomeCommandName,
		Description: "Welcome command to greetings to you",
		GuildID:     os.Getenv(serverGuildKey),
		Type:        discordgo.ChatApplicationCommand,
	}
}

func (w WelcomeCommand) Execute(_ context.Context, _ *discordgo.Session, i *discordgo.InteractionCreate) (DiscordMessageReceiver, error) {
	return SimpleMessageResponse{Msg: fmt.Sprintf("Witaj panie %s kapitanie! W czym mogę pomóc?", strings.ToUpper(i.Member.User.Username))}, nil
}
