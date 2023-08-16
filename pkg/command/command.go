package command

import "github.com/bwmarrin/discordgo"

type discordHandler func(s *discordgo.Session, i *discordgo.InteractionCreate)

type DiscordCommand struct {
	Command discordgo.ApplicationCommand
	Execute discordHandler
}
