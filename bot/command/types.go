package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/wittano/komputer/bot/log"
)

const (
	serverGuildKey = "SERVER_GUID"
)

type DiscordSlashCommandHandler interface {
	Command() *discordgo.ApplicationCommand
	DiscordEventHandler
}

type DiscordEventHandler interface {
	Execute(ctx log.Context, s *discordgo.Session, i *discordgo.InteractionCreate) (DiscordMessageReceiver, error)
}

type DiscordOptionMatcher interface {
	Match(customID string) bool
}

type DiscordMessageReceiver interface {
	Response() *discordgo.InteractionResponseData
}
