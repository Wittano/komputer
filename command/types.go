package command

import (
	"context"
	"github.com/bwmarrin/discordgo"
)

const (
	serverGuildKey = "SERVER_GUID"
	requestIDKey   = "requestID"
)

type DiscordSlashCommandHandler interface {
	Command() *discordgo.ApplicationCommand
	DiscordEventHandler
}

type DiscordEventHandler interface {
	Execute(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) (DiscordMessageReceiver, error)
}

type DiscordOptionMatcher interface {
	Match(customID string) bool
}

type DiscordMessageReceiver interface {
	Response() *discordgo.InteractionResponseData
}
