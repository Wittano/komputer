package command

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/wittano/komputer/internal/assets"
	"github.com/wittano/komputer/pkgs/voice"
	"log/slog"
	"os"
)

const (
	SpockCommandName = "spock"
	idOptionName     = "id"
)

type SpockCommand struct {
	GlobalCtx         context.Context
	SpockMusicStopChs map[string]chan struct{}
	GuildVoiceChats   map[string]voice.ChatInfo
}

func (sc SpockCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        SpockCommandName,
		Description: "Say funny world",
		GuildID:     os.Getenv("SERVER_GUID"),
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

func (sc SpockCommand) Execute(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) (DiscordMessageReceiver, error) {
	const spockQuote = "Kurwa Spock"
	defaultMsg := simpleMessageResponse{spockQuote, false}

	if _, ok := s.VoiceConnections[i.GuildID]; ok {
		return defaultMsg, nil
	}

	info, ok := sc.GuildVoiceChats[i.GuildID]
	if !ok || info.UserCount == 0 {
		slog.ErrorContext(ctx, fmt.Sprintf("failed find user on voice channels on %s server", i.Member.GuildID), fmt.Errorf("user with ID %s wasn't found on any voice chat on %s server", i.Member.User.ID, i.GuildID))

		return simpleMessageResponse{msg: "Kapitanie gdzie jesteś? Wejdź na kanał głosowy a ja dołącze"}, nil
	}

	go sc.playAudio(ctx, info.ChannelID, s, i)

	return defaultMsg, nil
}

func (sc SpockCommand) playAudio(ctx context.Context, channelID string, s *discordgo.Session, i *discordgo.InteractionCreate) {
	select {
	case <-ctx.Done():
		return
	default:
	}

	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildVoiceStates)

	voiceChat, err := s.ChannelVoiceJoin(i.GuildID, channelID, false, true)
	if err != nil {
		slog.ErrorContext(ctx, "failed join to voice channel", err)

		return
	}
	defer func(voiceChat *discordgo.VoiceConnection) {
		if err := voiceChat.Disconnect(); err != nil {
			slog.ErrorContext(ctx, "failed disconnect bot from voice channel", err)
		}
	}(voiceChat)
	defer voiceChat.Close()

	songPath, err := songPath(i.Data.(discordgo.ApplicationCommandInteractionData))
	if err != nil {
		slog.ErrorContext(ctx, "failed find song path", err)

		return
	}

	stopCh, ok := sc.SpockMusicStopChs[i.GuildID]
	if !ok {
		slog.ErrorContext(ctx, fmt.Sprintf("failed find user on voice channels on %s server", i.Member.GuildID), fmt.Errorf("user with ID %s wasn't found on any voice chat on %s server", i.Member.User.ID, i.GuildID))

		return
	}

	var (
		playingCtx context.Context
		cancel     context.CancelFunc
	)

	if audioDuration, err := voice.DuractionAudio(songPath); err != nil {
		slog.WarnContext(ctx, "failed calculated audio duration", err)

		playingCtx, cancel = context.WithCancel(context.Background())
	} else {
		playingCtx, cancel = context.WithTimeout(context.Background(), audioDuration)
	}

	defer cancel()

	if err = voice.PlayAudio(playingCtx, voiceChat, songPath, stopCh); err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("failed play '%s' songPath", songPath), err)
	}
}

func songPath(data discordgo.ApplicationCommandInteractionData) (path string, err error) {
	for _, o := range data.Options {
		switch o.Name {
		case idOptionName:
			path = o.Value.(string)
		default:
			path, err = assets.RandomAudio()
		}
	}

	if path == "" && err == nil {
		path, err = assets.RandomAudio()
	}

	return
}
