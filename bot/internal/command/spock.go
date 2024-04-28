package command

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/wittano/komputer/bot/internal/api"
	"github.com/wittano/komputer/bot/internal/voice"
	"log/slog"
	"os"
)

const (
	SpockCommandName = "spock"
	idOptionName     = "id"
)

type SpockCommand struct {
	GlobalCtx       context.Context
	MusicStopChs    map[string]chan struct{}
	GuildVoiceChats map[string]voice.ChatInfo
	WebAPI          api.WebClient
	Storage         voice.BotLocalStorage
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
				Name:        idOptionName,
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

	audioID, err := audioID(i.Data.(discordgo.ApplicationCommandInteractionData))
	if err != nil {
		logger.ErrorContext(ctx, "failed find song path", "error", err)

		return nil, nil
	}

	query := voice.SearchParams{Type: voice.IDType, Value: audioID}
	if _, err = sc.Storage.Get(ctx, query); err != nil && sc.WebAPI.Active(ctx) {
		go func() {
			logger.Info("download audio with id " + audioID)
			_, err := sc.WebAPI.DownloadAudio(audioID)
			if err != nil {
				logger.Error(fmt.Sprintf("failed download audio. %s", err))
				return
			}

			logger.Info("success download audio with id " + audioID)
			sc.playAudio(logger, s, i, info.ChannelID, audioID)
		}()

		return SimpleMessage{Msg: "Panie Kapitanie. Pobieram utwór. Proszę poczekać"}, nil
	} else if err != nil && !sc.WebAPI.Active(ctx) {
		return nil, err
	} else {
		go sc.playAudio(logger, s, i, info.ChannelID, audioID)
	}

	return msg, nil
}

func (sc SpockCommand) playAudio(
	l *slog.Logger,
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	channelID string,
	audioID string,
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

	if audioDuration, err := voice.Duration(audioID); err != nil {
		l.Warn("failed calculated audio duration", "error", err)

		playCtx, cancel = context.WithCancel(context.Background())
	} else {
		playCtx, cancel = context.WithTimeout(context.Background(), audioDuration)
	}
	defer cancel()

	if err = voice.Play(playCtx, voiceChat, voice.Path(audioID), stopCh); err != nil {
		l.ErrorContext(playCtx, fmt.Sprintf("failed play '%s' audioID", audioID), "error", err)
	}
}

func audioID(data discordgo.ApplicationCommandInteractionData) (path string, err error) {
	for _, o := range data.Options {
		switch o.Name {
		case idOptionName:
			path = o.Value.(string)
		default:
			path, err = voice.RandomAudioID()
		}
	}

	if path == "" && err == nil {
		path, err = voice.RandomAudioID()
	}

	return
}
