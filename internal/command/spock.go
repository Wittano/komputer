package command

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/wittano/komputer/internal/assets"
	"github.com/wittano/komputer/internal/interaction"
	"github.com/wittano/komputer/internal/voice"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// TODO Replace chan bool by chan struct{} (0-bytes channel)
var SpockMusicStopCh = map[string]chan bool{}

var SpockCommand = DiscordCommand{
	Command: discordgo.ApplicationCommand{
		Name:        "spock",
		Description: "Say funny world",
		GuildID:     os.Getenv("SERVER_GUID"),
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Required:    false,
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "id",
				Choices:     audioIdOptions(),
				Description: "Id of audio asset",
			},
		},
	},
	Execute: execSpookSpeak,
}

// TODO Add context
func audioIdOptions() []*discordgo.ApplicationCommandOptionChoice {
	list, err := assets.Audios()
	if err != nil {
		slog.Error("failed to get audios form assets folder", err)
	}

	result := make([]*discordgo.ApplicationCommandOptionChoice, len(list))
	for i, v := range list {
		result[i] = &discordgo.ApplicationCommandOptionChoice{
			Name:  strings.TrimSuffix(filepath.Base(v), filepath.Ext(v)),
			Value: v,
		}
	}

	return result
}

func execSpookSpeak(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if _, ok := s.VoiceConnections[i.GuildID]; ok {
		interaction.CreateDiscordInteractionResponse(ctx, i, s, interaction.CreateDiscordMsg("Kurwa Spock"))
		return
	}

	userVoiceChat, ok := voice.UserVoiceChatMap[i.Member.User.ID]
	if !ok || userVoiceChat.GuildID != i.GuildID {
		slog.ErrorContext(ctx, fmt.Sprintf("failed find user on voice channels on %s server", i.Member.GuildID), fmt.Errorf("user with ID %s wasn't found on any voice chat on %s server", i.Member.User.ID, i.GuildID))
		interaction.CreateDiscordInteractionResponse(ctx, i, s, interaction.CreateDiscordMsg("BEEP BOOP. Kapitanie gdzie jesteś? Wejdź na kanał głosowy a ja dołącze"))
		return
	}

	go func() {
		s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildVoiceStates)

		voiceJoin, err := s.ChannelVoiceJoin(i.GuildID, userVoiceChat.ChannelID, false, true)
		if err != nil {
			slog.ErrorContext(ctx, "failed join to voice channel", err)
			return
		}
		defer voiceJoin.Close()

		var songPath string

		for _, o := range i.Data.(discordgo.ApplicationCommandInteractionData).Options {
			switch o.Name {
			case "id":
				songPath = o.Value.(string)
			default:
				songPath, err = assets.RandomAudio()
				if err != nil {
					slog.ErrorContext(ctx, "failed find any audio in assets directory", err)
					return
				}
			}
		}

		if songPath == "" {
			songPath, err = assets.RandomAudio()
			if err != nil {
				slog.ErrorContext(ctx, "failed find any audio in assets directory", err)
				return
			}
		}

		stop := make(chan bool)
		SpockMusicStopCh[i.GuildID] = stop

		if err = voice.PlayAudio(voiceJoin, songPath, stop); err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("failed play '%s' songPath", songPath), err)
		}

		voiceJoin.Disconnect()
	}()

	interaction.CreateDiscordInteractionResponse(ctx, i, s, interaction.CreateDiscordMsg("Kurwa Spock"))
}
