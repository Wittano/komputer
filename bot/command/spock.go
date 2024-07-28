package command

import (
	"context"
	"fmt"
	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"github.com/wittano/komputer/internal/audio"
	"github.com/wittano/komputer/bot/log"
	"github.com/wittano/komputer/bot/voice"
	"log/slog"
	"os"
)

const (
	SpockCommandName = "spock"
	nameOptionName   = "name"
)

type SpockCommand struct {
	GlobalCtx       context.Context
	MusicStopChs    map[string]chan struct{}
	GuildVoiceChats map[string]voice.ChatInfo
}

func (sc SpockCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        SpockCommandName,
		Description: "Say funny world",
		GuildID:     os.Getenv(serverGuildKey),
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Required:    false,
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        nameOptionName,
				Description: "ID of audio asset",
			},
		},
	}
}

func (sc SpockCommand) Execute(
	ctx context.Context,
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) (DiscordMessageReceiver, error) {
	const spockQuote = "Kurwa Spock"
	msg := SimpleMessage{spockQuote, false}

	if _, ok := s.VoiceConnections[i.GuildID]; ok {
		return msg, nil
	}

	info, ok := sc.GuildVoiceChats[i.GuildID]
	if !ok || info.UserCount == 0 {
		log.Log(ctx, func(l slog.Logger) {
			l.Error(fmt.Sprintf(fmt.Sprintf("user with ID '%s' wasn't found on any voice chat on '%s' server", i.Member.User.ID, i.GuildID)))
		})

		return SimpleMessage{Msg: "Kapitanie gdzie jesteś? Wejdź na kanał głosowy a ja dołącze"}, nil
	}

	path, err := audioPath(i.Data.(discordgo.ApplicationCommandInteractionData))
	if err != nil {
		log.Log(ctx, func(l slog.Logger) {
			l.Error("failed find song path", "error", err)
		})

		return SimpleMessage{Msg: "Panie kapitanie, nie znalazłem utworu"}, nil
	}

	go sc.playAudio(log.NewCtxWithRequestID(ctx), s, i, info.ChannelID, path)

	return msg, nil
}

func (sc SpockCommand) playAudio(
	loggerCtx log.Context,
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	channelID string,
	path string,
) {
	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildVoiceStates)

	voiceChat, err := s.ChannelVoiceJoin(i.GuildID, channelID, false, true)
	if err != nil {
		loggerCtx.Logger.Error("failed join to voice channel", "error", err)

		return
	}
	defer func(voiceChat *discordgo.VoiceConnection) {
		if err := voiceChat.Disconnect(); err != nil {
			loggerCtx.Logger.Error("failed disconnect discord from voice channel", "error", err)
		}
	}(voiceChat)
	defer voiceChat.Close()

	stopCh, ok := sc.MusicStopChs[i.GuildID]
	if !ok {
		loggerCtx.Logger.Error(fmt.Sprintf("failed find user on voice channels on '%s' server", i.Member.GuildID), "error", fmt.Errorf("user with ID '%s' wasn't found on any voice chat on '%s' server", i.Member.User.ID, i.GuildID))

		return
	}

	var (
		playCtx context.Context
		cancel  context.CancelFunc
	)

	if audioDuration, err := audio.Duration(path); err != nil {
		loggerCtx.Logger.Warn("failed calculated audio duration", "error", err)

		playCtx, cancel = context.WithCancel(loggerCtx)
	} else {
		playCtx, cancel = context.WithTimeout(loggerCtx, audioDuration)
	}
	defer cancel()

	if err = dgvoice.PlayAudioFile(playCtx, voiceChat, path, stopCh); err != nil {
		loggerCtx.Logger.ErrorContext(playCtx, fmt.Sprintf("failed play '%s' audioPath", path), "error", err)
	}
}

func audioPath(data discordgo.ApplicationCommandInteractionData) (path string, err error) {
	var name string
	for _, o := range data.Options {
		switch o.Name {
		case nameOptionName:
			name = o.Value.(string)
		default:
		}
	}

	if name == "" {
		path, err = audio.RandomAudioName()
	} else {
		path, err = audio.Path(path)
	}

	return
}
