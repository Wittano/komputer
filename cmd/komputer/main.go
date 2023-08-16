package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/wittano/komputer/internal"
	"github.com/wittano/komputer/pkg/command"
	"log"
	"os"
	"os/signal"
)

var (
	bot      *discordgo.Session
	commands = map[string]command.DiscordCommand{
		command.WelcomeCommand.Command.Name: command.WelcomeCommand,
		command.JokeCommand.Command.Name:    command.JokeCommand,
	}
)

func init() {
	godotenv.Load()
}

func init() {
	checkEnvVariables("DISCORD_BOT_TOKEN")

	var err error

	bot, err = discordgo.New(fmt.Sprintf("Bot %s", os.Getenv("DISCORD_BOT_TOKEN")))
	if err != nil {
		panic(err)
	}
}

func init() {
	bot.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionMessageComponent {
			if handler, ok := internal.JokeMessageComponentHandler[i.Data.(discordgo.MessageComponentInteractionData).CustomID]; ok {
				handler(s, i)
				return
			}
		}

		if c, ok := commands[i.ApplicationCommandData().Name]; ok {
			c.Execute(s, i)
		}
	})
}

func init() {
	checkEnvVariables("APPLICATION_ID", "SERVER_GUID")

	for _, c := range commands {
		if _, err := bot.ApplicationCommandCreate(
			os.Getenv("APPLICATION_ID"),
			os.Getenv("SERVER_GUID"),
			&c.Command,
		); err != nil {
			log.Print(err)
		}
	}
}

func checkEnvVariables(vars ...string) {
	for _, v := range vars {
		if _, ok := os.LookupEnv(v); !ok {
			log.Fatalf("Missing %s varaiable", v)
		}
	}
}

func main() {
	bot.Open()
	defer bot.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Print("Bot is ready!. Press Ctrl+C to exit")
	<-stop
}
