package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/wittano/komputer/internal"
	"github.com/wittano/komputer/internal/command"
	"github.com/wittano/komputer/internal/voice"
	"github.com/wittano/komputer/pkgs/config"
	"github.com/wittano/komputer/pkgs/db"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"time"
)

var (
	commands = map[string]command.DiscordCommand{
		command.WelcomeCommand.String():   command.WelcomeCommand,
		command.JokeCommand.String():      command.JokeCommand,
		command.AddJokeCommand.String():   command.AddJokeCommand,
		command.SpockCommand.String():     command.SpockCommand,
		command.SpockStopCommand.String(): command.SpockStopCommand,
	}
)

type slashCommandHandler struct {
	ctx context.Context
}

type DiscordBot struct {
	ctx     context.Context
	bot     *discordgo.Session
	mongodb db.MongodbService
}

func (d *DiscordBot) Start() (err error) {
	if err = d.bot.Open(); err != nil {
		return
	}

	slog.Info("Bot is ready!. Press Ctrl+C to exit")

	return
}

func (d *DiscordBot) Close() (err error) {
	err = d.mongodb.Close()
	err = d.bot.Close()

	return
}

func newDiscordBot(ctx context.Context) (*DiscordBot, error) {
	prop, err := config.NewBotProperties()
	if err != nil {
		return nil, err
	}

	bot, err := discordgo.New(fmt.Sprintf("Bot %s", prop.Token))
	if err != nil {
		return nil, fmt.Errorf("failed connect with Discord: %s", err)
	}

	// Update list of current user on voice channels
	bot.AddHandler(voice.HandleVoiceChannelUpdate)

	// Register slash commands
	for _, c := range commands {
		if _, err := bot.ApplicationCommandCreate(
			prop.AppID,
			prop.ServerGUID, // If empty, command registers globally
			&c.Command,
		); err != nil {
			return nil, fmt.Errorf("registration slash command failed: %s", err)
		}
	}

	// General handler for slash commands
	handler := slashCommandHandler{ctx}

	bot.AddHandler(handler.handleSlashCommand)

	return &DiscordBot{
		ctx: ctx,
		bot: bot,
		// Create connection with database (optional)
		// Databases doesn't require to running bot
		mongodb: db.NewMongodbDatabase(ctx),
	}, nil
}

func (sc slashCommandHandler) handleSlashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.WithValue(sc.ctx, "requestID", uuid.New().String())
	deadlineCtx, cancel := context.WithTimeout(ctx, time.Second*2)
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
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bot, err := newDiscordBot(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer bot.Close()
	if err := bot.Start(); err != nil {
		log.Fatal(err)
		return
	}

	stop := make(chan os.Signal)
	defer close(stop)

	signal.Notify(stop, os.Interrupt)

	for {
		select {
		case <-ctx.Done():
		case _ = <-stop:
			return
		}
	}
}
