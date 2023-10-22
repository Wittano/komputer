package command

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/wittano/komputer/internal/interaction"
	"os"
)

var SpockStopCommand = DiscordCommand{
	Command: discordgo.ApplicationCommand{
		Name:        "stop",
		Description: "Stop playing song by bot",
		GuildID:     os.Getenv("SERVER_GUID"),
		Type:        discordgo.ChatApplicationCommand,
	},
	Execute: execSpookStopSpeak,
}

func execSpookStopSpeak(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if _, ok := s.VoiceConnections[i.GuildID]; ok {
		SpockMusicStopCh[i.GuildID] <- true
	}

	interaction.CreateDiscordInteractionResponse(ctx, i, s, interaction.CreateDiscordMsg("BEEP BOOP. Przepraszam"))
}
