package command

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/wittano/komputer/internal/interaction"
	"os"
	"strings"
)

var WelcomeCommand = DiscordCommand{
	Command: discordgo.ApplicationCommand{
		Name:        "welcome",
		Description: "Welcome command to greetings to you",
		GuildID:     os.Getenv("SERVER_GUID"),
		Type:        discordgo.ChatApplicationCommand,
	},
	Execute: executeWelcomeCommand,
}

func executeWelcomeCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	interaction.CreateDiscordInteractionResponse(ctx, i, s, interaction.CreateDiscordMsg(fmt.Sprintf("Witaj panie %s kapitanie! W czym mogę pomóc?", strings.ToUpper(i.Member.User.Username))))
}
