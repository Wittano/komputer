package command

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"os"
)

const SpockStopCommandName = "stop"

type SpockStopCommand struct {
	SpockMusicStopChs map[string]chan struct{}
}

func (ssc SpockStopCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        SpockStopCommandName,
		Description: "Stop playing song by discord",
		GuildID:     os.Getenv(serverGuildKey),
		Type:        discordgo.ChatApplicationCommand,
	}
}

func (ssc SpockStopCommand) Execute(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) (res DiscordMessageReceiver, _ error) {
	res = SimpleMessageResponse{Msg: "Przepraszam"}

	select {
	case <-ctx.Done():
		return
	default:
		if _, ok := s.VoiceConnections[i.GuildID]; ok {
			ssc.SpockMusicStopChs[i.GuildID] <- struct{}{}
		}
	}

	return
}
