package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/wittano/komputer/internal/voice"
	"github.com/wittano/komputer/pkgs/command"
	"github.com/wittano/komputer/pkgs/config"
	"github.com/wittano/komputer/pkgs/db"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"time"
)

type slashCommandHandler struct {
	ctx      context.Context
	commands map[string]command.DiscordSlashCommandHandler
	options  map[string]command.DiscordEventHandler
}

func (sc slashCommandHandler) handleSlashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	deadlineCtx, cancel := context.WithTimeout(sc.ctx, time.Second*2)
	requestIDCtx := context.WithValue(deadlineCtx, "requestID", uuid.New().String())
	ctx := context.WithValue(requestIDCtx, db.GuildIDKey, i.GuildID)
	defer cancel()

	userID := i.Member.User.ID
	// Handle options assigned to slash commands
	if i.Type == discordgo.InteractionMessageComponent {
		if option, ok := sc.options[i.Data.(discordgo.MessageComponentInteractionData).CustomID]; ok {
			slog.InfoContext(ctx, fmt.Sprintf("user '%s' select '%s' option", userID, i.Data.(discordgo.MessageComponentInteractionData).CustomID))

			handleEventResponse(ctx, s, i, option)

			return
		}
	}

	// Handle slash commands
	cmdName := i.ApplicationCommandData().Name
	if cmd, ok := sc.commands[cmdName]; ok {
		slog.InfoContext(ctx, fmt.Sprintf("user '%s' execute slash command '%s'", userID, cmdName))

		handleEventResponse(ctx, s, i, cmd)
	}
}

func handleEventResponse(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, event command.DiscordEventHandler) {
	msg, err := event.Execute(ctx, s, i)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
	}

	if err != nil && errors.Is(err, command.ErrorResponse{}) {
		errors.As(err, &msg)
	}

	command.CreateDiscordInteractionResponse(ctx, i, s, msg)
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
	commands := createCommands(ctx)
	for _, c := range commands {
		if _, err := bot.ApplicationCommandCreate(
			prop.AppID,
			prop.ServerGUID, // If empty, command registers globally
			c.Command(),
		); err != nil {
			return nil, fmt.Errorf("registration slash command failed: %s", err)
		}
	}

	// General handler for slash commands
	handler := slashCommandHandler{ctx, commands, createOptions(ctx)}

	bot.AddHandler(handler.handleSlashCommand)

	return &DiscordBot{
		ctx: ctx,
		bot: bot,
	}, nil
}

func createCommands(ctx context.Context) map[string]command.DiscordSlashCommandHandler {
	mongodb := db.NewMongodbDatabase(ctx)
	jokeService := db.JokeService{Mongodb: mongodb}

	welcomeCmd := command.WelcomeCommand{}
	addJokeCmd := command.AddJokeCommand{Service: jokeService}
	jokeCmd := command.JokeCommand{Service: jokeService}

	return map[string]command.DiscordSlashCommandHandler{
		command.WelcomeCommandName: welcomeCmd,
		command.AddJokeCommandName: addJokeCmd,
		command.GetJokeCommandName: jokeCmd,
	}
}

func createOptions(ctx context.Context) map[string]command.DiscordEventHandler {
	mongodb := db.NewMongodbDatabase(ctx)
	jokeService := db.JokeService{Mongodb: mongodb}

	welcomeCmd := command.ApologiesOption{}
	addJokeCmd := command.NextJokeOption{Service: jokeService}
	jokeCmd := command.SameJokeCategoryOption{Service: jokeService}

	return map[string]command.DiscordEventHandler{
		command.ApologiesButtonName: welcomeCmd,
		command.AddJokeCommandName:  addJokeCmd,
		command.GetJokeCommandName:  jokeCmd,
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
