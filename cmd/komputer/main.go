package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	zerolog "github.com/rs/zerolog/log"
	"github.com/wittano/komputer/internal"
	"github.com/wittano/komputer/internal/command"
	"github.com/wittano/komputer/internal/log"
	"github.com/wittano/komputer/internal/mongo"
	"github.com/wittano/komputer/internal/schedule"
	"github.com/wittano/komputer/internal/voice"
	"os"
	"os/signal"
)

var (
	bot      *discordgo.Session
	commands = map[string]command.DiscordCommand{
		command.WelcomeCommand.Command.Name:   command.WelcomeCommand,
		command.JokeCommand.Command.Name:      command.JokeCommand,
		command.AddJokeCommand.Command.Name:   command.AddJokeCommand,
		command.SpockCommand.Command.Name:     command.SpockCommand,
		command.SpockStopCommand.Command.Name: command.SpockStopCommand,
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
		ctx := context.WithValue(context.Background(), "traceID", uuid.New().String())

		if i.Type == discordgo.InteractionMessageComponent {
			if handler, ok := internal.JokeMessageComponentHandler[i.Data.(discordgo.MessageComponentInteractionData).CustomID]; ok {
				log.Info(ctx, fmt.Sprintf("User %s execute message component action '%s'", i.Member.User.ID, i.Data.(discordgo.MessageComponentInteractionData).CustomID))
				handler(ctx, s, i)
				return
			}
		}

		if c, ok := commands[i.ApplicationCommandData().Name]; ok {
			log.Info(ctx, fmt.Sprintf("User %s execute slash command '%s'", i.Member.User.ID, i.ApplicationCommandData().Name))
			c.Execute(ctx, s, i)
		}
	})
}

func init() {
	checkEnvVariables("APPLICATION_ID")

	for _, c := range commands {
		if _, err := bot.ApplicationCommandCreate(
			os.Getenv("APPLICATION_ID"),
			os.Getenv("SERVER_GUID"),
			&c.Command,
		); err != nil {
			log.Error(context.Background(), "Registration slash command failed", err)
		}
	}
}

func init() {
	bot.AddHandler(voice.HandleVoiceChannelUpdate)
}

func checkEnvVariables(vars ...string) {
	for _, v := range vars {
		if _, ok := os.LookupEnv(v); !ok {
			zerolog.Fatal().Msg(fmt.Sprintf("Missing %s varaiable", v))
		}
	}
}

func main() {
	bot.Open()
	defer bot.Close()
	defer schedule.Scheduler.Stop()
	defer mongo.CloseDb()
	defer log.FileLog.Close()

	schedule.Scheduler.StartAsync()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Info(context.Background(), "Bot is ready!. Press Ctrl+C to exit")
	<-stop
}
