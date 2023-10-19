package command

import "github.com/bwmarrin/discordgo"

func getJokeCategoryOption(required bool) *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "category",
		Description: "Joke category",
		Required:    required,
		Choices: []*discordgo.ApplicationCommandOptionChoice{
			{
				Name:  "Programowanie",
				Value: "Programming",
			},
			{
				Name:  "Różne",
				Value: "Misc",
			},
			{
				Name:  "Czarny humor",
				Value: "Dark",
			},
			{
				Name:  "Straszne",
				Value: "Spooky",
			},
		},
	}
}

func getJokeTypeOption(required bool) *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "type",
		Description: "Type of joke",
		Required:    required,
		Choices: []*discordgo.ApplicationCommandOptionChoice{
			{
				Name:  "Single",
				Value: "single",
			},
			{

				Name:  "Two-Part",
				Value: "twopart",
			},
		},
	}
}
