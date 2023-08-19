package command

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	zerolog "github.com/rs/zerolog/log"
	"os"
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
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Welcome %s! How can I help you?", i.Member.User.Username),
		},
	})

	if err != nil {
		zerolog.Err(err).Ctx(ctx).Str("traceID", ctx.Value("traceID").(string)).Msg("Failed send response message")
	}
}
