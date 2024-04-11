package command

import (
	"context"
	"github.com/bwmarrin/discordgo"
)

const requestIDKey = "requestID"

type DiscordSlashCommandHandler interface {
	Command() *discordgo.ApplicationCommand
	DiscordEventHandler
}

type DiscordEventHandler interface {
	Execute(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) (DiscordMessageReceiver, error)
}

type DiscordMessageReceiver interface {
	Response() *discordgo.InteractionResponseData
}
