package command

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

const komputerMsgPrefix = "BEEP BOOP. "

type ErrorResponse struct {
	Err error
	Msg string
}

func (e ErrorResponse) Error() string {
	return e.Err.Error()
}

func (e ErrorResponse) Response() *discordgo.InteractionResponseData {
	if e.Msg == "" {
		e.Msg = komputerMsgPrefix + "Coś poszło nie tak :("
	} else {
		e.Msg = komputerMsgPrefix + e.Msg
	}

	return &discordgo.InteractionResponseData{Content: e.Msg}
}

type SimpleMessageResponse struct {
	Msg    string
	Hidden bool
}

func (s SimpleMessageResponse) Response() (msg *discordgo.InteractionResponseData) {
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
		slog.With(requestIDKey, ctx.Value(requestIDKey)).ErrorContext(ctx, "failed send response to discord user", "error", err)
	}
}
