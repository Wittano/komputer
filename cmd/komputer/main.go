package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/wittano/komputer/internal"
	"github.com/wittano/komputer/internal/command"
	"github.com/wittano/komputer/internal/mongo"
	"github.com/wittano/komputer/internal/schedule"
	"github.com/wittano/komputer/internal/voice"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"time"
)

var (
	bot      *discordgo.Session
	commands = map[string]command.DiscordCommand{
		command.WelcomeCommand.String():   command.WelcomeCommand,
		command.JokeCommand.String():      command.JokeCommand,
		command.AddJokeCommand.String():   command.AddJokeCommand,
		command.SpockCommand.String():     command.SpockCommand,
		command.SpockStopCommand.String(): command.SpockStopCommand,
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
		log.Fatalf("failed connect with Discord: %s", err)
	}
}

func init() {
	bot.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		ctx := context.WithValue(context.Background(), "requestID", uuid.New().String())
		deadlineCtx, cancel := context.WithTimeout(ctx, time.Second*3)
		defer cancel()

		if i.Type == discordgo.InteractionMessageComponent {
			if handler, ok := internal.JokeMessageComponentHandler[i.Data.(discordgo.MessageComponentInteractionData).CustomID]; ok {
				slog.InfoContext(deadlineCtx, fmt.Sprintf("User %s execute message component action '%s'", i.Member.User.ID, i.Data.(discordgo.MessageComponentInteractionData).CustomID))
				handler(deadlineCtx, s, i)
				return
			}
		}

		if c, ok := commands[i.ApplicationCommandData().Name]; ok {
			slog.InfoContext(deadlineCtx, fmt.Sprintf("User %s execute slash command '%s'", i.Member.User.ID, i.ApplicationCommandData().Name))
			c.Execute(deadlineCtx, s, i)
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
			log.Fatalf("registration slash command failed: %s", err)
		}
	}
}

func init() {
	bot.AddHandler(voice.HandleVoiceChannelUpdate)
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
	defer schedule.Scheduler.Stop()
	defer mongo.CloseDb()

	schedule.Scheduler.StartAsync()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	slog.Info("Bot is ready!. Press Ctrl+C to exit")
	<-stop
}
