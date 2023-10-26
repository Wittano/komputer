package command

import (
	"context"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/wittano/komputer/internal/assets"
	"github.com/wittano/komputer/internal/interaction"
	"github.com/wittano/komputer/internal/log"
	"github.com/wittano/komputer/internal/voice"
	"math/rand"
	"os"
)

var SpockMusicStopCh = map[string]chan bool{}

var SpockCommand = DiscordCommand{
	Command: discordgo.ApplicationCommand{
		Name:        "spock",
		Description: "Say funny world",
		GuildID:     os.Getenv("SERVER_GUID"),
		Type:        discordgo.ChatApplicationCommand,
	},
	Execute: execSpookSpeak,
}

func execSpookSpeak(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if _, ok := s.VoiceConnections[i.GuildID]; ok {
		interaction.CreateDiscordInteractionResponse(ctx, i, s, interaction.CreateDiscordMsg("Kurwa Spock"))
		return
	}

	userVoiceChat, ok := voice.UserVoiceChatMap[i.Member.User.ID]
	if !ok || userVoiceChat.GuildID != i.GuildID {
		log.Error(ctx, fmt.Sprintf("Failed find user on voice channels on %s server", i.Member.GuildID), errors.New(fmt.Sprintf("user with ID %s wasn't found on any voice chat on %s server", i.Member.User.ID, i.GuildID)))
		interaction.CreateDiscordInteractionResponse(ctx, i, s, interaction.CreateDiscordMsg("BEEP BOOP. Kapitanie gdzie jesteś? Wejdź na kanał głosowy a ja dołącze"))
		return
	}

	go func() {
		s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildVoiceStates)

		voiceJoin, err := s.ChannelVoiceJoin(i.GuildID, userVoiceChat.ChannelID, false, true)
		if err != nil {
			log.Error(ctx, "Failed join to voice channel", err)
			return
		}
		defer voiceJoin.Close()

		path, err := assets.GetAudioPaths("spock")
		if err != nil || len(path) == 0 {
			log.Error(ctx, "Failed find any assert \"spock\" in assets directory", err)
			return
		}

		ch := make(chan bool)
		SpockMusicStopCh[i.GuildID] = ch

		songPath := path[rand.Int()%len(path)]
		if err = voice.PlayAudio(voiceJoin, songPath, ch); err != nil {
			log.Error(ctx, fmt.Sprintf("Failed play '%s' audio", songPath), err)
		}

		voiceJoin.Disconnect()
	}()

	interaction.CreateDiscordInteractionResponse(ctx, i, s, interaction.CreateDiscordMsg("Kurwa Spock"))
}
