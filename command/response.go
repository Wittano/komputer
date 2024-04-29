package command

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/wittano/komputer/log"
	"log/slog"
)

const komputerMsgPrefix = "BEEP BOOP. "

type DiscordError struct {
	Err error
	Msg string
}

func (e DiscordError) Error() string {
	return e.Err.Error()
}

func (e DiscordError) Response() *discordgo.InteractionResponseData {
	if e.Msg == "" {
		e.Msg = komputerMsgPrefix + "Coś poszło nie tak :("
	} else {
		e.Msg = komputerMsgPrefix + e.Msg
	}

	return &discordgo.InteractionResponseData{Content: e.Msg}
}

type SimpleMessage struct {
	Msg    string
	Hidden bool
}

func (s SimpleMessage) Response() (msg *discordgo.InteractionResponseData) {
	msg = &discordgo.InteractionResponseData{
		Content: s.Msg,
	}

	if s.Hidden {
		msg.Flags = discordgo.MessageFlagsEphemeral
	}

	return
}

func CreateDiscordInteractionResponse(ctx context.Context, i *discordgo.InteractionCreate, s *discordgo.Session, msg DiscordMessageReceiver) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: msg.Response(),
	}); err != nil {
		log.Log(ctx, func(l slog.Logger) {
			l.Error("failed send response to discord user", "error", err)
		})
	}
}
