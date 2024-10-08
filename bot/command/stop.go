package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/wittano/komputer/bot/log"
	"os"
)

const StopCommandName = "stop"

type StopCommand struct {
	MusicStopChs map[string]chan struct{}
}

func (sc StopCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        StopCommandName,
		Description: "Stop playing song by discord",
		GuildID:     os.Getenv(serverGuildKey),
		Type:        discordgo.ChatApplicationCommand,
	}
}

func (sc StopCommand) Execute(
	ctx log.Context,
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) (res DiscordMessageReceiver, _ error) {
	res = SimpleMessage{Msg: "Przepraszam"}

	select {
	case <-ctx.Done():
		return
	default:
		if _, ok := s.VoiceConnections[i.GuildID]; ok {
			sc.MusicStopChs[i.GuildID] <- struct{}{}
		}
	}

	return
}
