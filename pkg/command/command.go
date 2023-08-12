package command

import "github.com/bwmarrin/discordgo"

type DiscordCommand struct {
	Command discordgo.ApplicationCommand
	Execute func(s *discordgo.Session, i *discordgo.InteractionCreate)
}
