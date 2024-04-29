package command

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/wittano/komputer/audio"
	"github.com/wittano/komputer/voice"
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
		slog.With(requestIDKey, ctx.Value(requestIDKey)).ErrorContext(ctx, fmt.Sprintf(fmt.Sprintf("user with ID '%s' wasn't found on any voice chat on '%s' server", i.Member.User.ID, i.GuildID)))

		return SimpleMessage{Msg: "Kapitanie gdzie jesteś? Wejdź na kanał głosowy a ja dołącze"}, nil
	}

	logger := slog.With(requestIDKey, ctx.Value(requestIDKey))

	path, err := audioPath(i.Data.(discordgo.ApplicationCommandInteractionData))
	if err != nil {
		logger.ErrorContext(ctx, "failed find song path", "error", err)

		return SimpleMessage{Msg: "Panie kapitanie, nie znalazłem utworu"}, nil
	}

	go sc.playAudio(logger, s, i, info.ChannelID, path)

	return msg, nil
}

func (sc SpockCommand) playAudio(
	l *slog.Logger,
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	channelID string,
	path string,
) {
	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildVoiceStates)

	voiceChat, err := s.ChannelVoiceJoin(i.GuildID, channelID, false, true)
	if err != nil {
		l.Error("failed join to voice channel", "error", err)

		return
	}
	defer func(voiceChat *discordgo.VoiceConnection) {
		if err := voiceChat.Disconnect(); err != nil {
			l.Error("failed disconnect discord from voice channel", "error", err)
		}
	}(voiceChat)
	defer voiceChat.Close()

	stopCh, ok := sc.MusicStopChs[i.GuildID]
	if !ok {
		l.Error(fmt.Sprintf("failed find user on voice channels on '%s' server", i.Member.GuildID), "error", fmt.Errorf("user with ID '%s' wasn't found on any voice chat on '%s' server", i.Member.User.ID, i.GuildID))

		return
	}

	var (
		playCtx context.Context
		cancel  context.CancelFunc
	)

	if audioDuration, err := audio.Duration(path); err != nil {
		l.Warn("failed calculated audio duration", "error", err)

		playCtx, cancel = context.WithCancel(context.Background())
	} else {
		playCtx, cancel = context.WithTimeout(context.Background(), audioDuration)
	}
	defer cancel()

	if err = voice.Play(playCtx, voiceChat, path, stopCh); err != nil {
		l.ErrorContext(playCtx, fmt.Sprintf("failed play '%s' audioPath", path), "error", err)
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

	if name == "" && err == nil {
		path, err = audio.RandomPath()
	} else {
		path = audio.Path(path)
	}

	return
}
