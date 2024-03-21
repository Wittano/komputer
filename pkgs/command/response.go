package command

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

const komputerMsgPrefix = "BEEP BOOP. "

type ErrorResponse struct {
	err error
	msg string
}

func (e ErrorResponse) Error() string {
	return e.err.Error()
}

func (e ErrorResponse) Response() *discordgo.InteractionResponseData {
	if e.msg == "" {
		e.msg = komputerMsgPrefix + "Coś poszło nie tak :("
	} else {
		e.msg = komputerMsgPrefix + e.msg
	}

	return &discordgo.InteractionResponseData{Content: e.msg}
}

type simpleMessageResponse struct {
	msg    string
	hidden bool
}

func (s simpleMessageResponse) Response() (msg *discordgo.InteractionResponseData) {
	msg = &discordgo.InteractionResponseData{
		Content: s.msg,
	}

	if !s.hidden {
		msg.Flags = discordgo.MessageFlagsEphemeral
	}

	return
}

func CreateDiscordInteractionResponse(ctx context.Context, i *discordgo.InteractionCreate, s *discordgo.Session, msg DiscordMessageReceiver) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: msg.Response(),
	}); err != nil {
		slog.ErrorContext(ctx, "failed send response to discord user", err)
	}
}
