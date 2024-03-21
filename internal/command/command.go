package command

import (
	"context"
	"github.com/bwmarrin/discordgo"
)

type discordHandler func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate)

type DiscordCommand struct {
	Command discordgo.ApplicationCommand
	Execute discordHandler
}

func (d DiscordCommand) String() string {
	return d.Command.Name
}
