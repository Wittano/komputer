package bot

import (
	"context"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/wittano/komputer/bot/internal/command"
	"github.com/wittano/komputer/bot/internal/config"
	"github.com/wittano/komputer/bot/internal/db"
	"github.com/wittano/komputer/bot/internal/joke"
	"github.com/wittano/komputer/bot/internal/voice"
	"log/slog"
	"time"
)

const (
	jokeDevServiceID  = 0
	humorAPIServiceID = 1
	databaseServiceID = 2
)

const requestIDKey = "requestID"

type slashCommandHandler struct {
	ctx      context.Context
	commands map[string]command.DiscordSlashCommandHandler
	options  map[string]command.DiscordEventHandler
}

func (sc slashCommandHandler) handleSlashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	const cmdTimeout = time.Second * 2

	deadlineCtx, cancel := context.WithTimeout(sc.ctx, cmdTimeout)
	requestIDCtx := context.WithValue(deadlineCtx, requestIDKey, uuid.New().String())
	ctx := context.WithValue(requestIDCtx, joke.GuildIDKey, i.GuildID)
	defer cancel()

	logger := slog.With(requestIDKey, ctx.Value(requestIDKey))

	userID := i.Member.User.ID
	// Handle options assigned to slash commands
	if i.Type == discordgo.InteractionMessageComponent {
		if option, ok := sc.options[i.Data.(discordgo.MessageComponentInteractionData).CustomID]; ok {
			logger.InfoContext(ctx, fmt.Sprintf("user '%s' select '%s' option", userID, i.Data.(discordgo.MessageComponentInteractionData).CustomID))

			handleEventResponse(ctx, s, i, option)

			return
		}
	}

	// Handle slash commands
	cmdName := i.ApplicationCommandData().Name
	if cmd, ok := sc.commands[cmdName]; ok {
		logger.InfoContext(ctx, fmt.Sprintf("user '%s' execute slash command '%s'", userID, cmdName))

		handleEventResponse(ctx, s, i, cmd)
	} else {
		msg := command.SimpleMessageResponse{Msg: "Kapitanie co chcesz zrobić?", Hidden: true}

		logger.WarnContext(ctx, "someone try execute unknown command")
		command.CreateDiscordInteractionResponse(ctx, i, s, msg)
	}
}

func handleEventResponse(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, event command.DiscordEventHandler) {
	msg, err := event.Execute(ctx, s, i)

	if errors.Is(err, command.ErrorResponse{}) {
		errors.As(err, &msg)
	} else if err != nil {
		slog.With(requestIDKey, ctx.Value(requestIDKey)).ErrorContext(ctx, err.Error())

		msg = command.ErrorResponse{
			Err: err,
		}
	}

	command.CreateDiscordInteractionResponse(ctx, i, s, msg)
}

type DiscordBot struct {
	ctx           context.Context
	bot           *discordgo.Session
	mongodb       db.MongodbService
	spockVoiceChs voice.SpockVoiceChannels
}

func (d *DiscordBot) Start() (err error) {
	if err = d.bot.Open(); err != nil {
		return
	}

	slog.Info("Bot is ready!. Press Ctrl+C to exit")

	return
}

func (d *DiscordBot) Close() (err error) {
	err = d.spockVoiceChs.Close()
	err = d.mongodb.Close()
	err = d.bot.Close()

	return
}

func createJokeGetServices(globalCtx context.Context, database *db.MongodbDatabase) []joke.GetService {
	return []joke.GetService{
		jokeDevServiceID:  joke.NewJokeDevService(globalCtx),
		humorAPIServiceID: joke.NewHumorAPIService(globalCtx),
		databaseServiceID: joke.NewDatabaseJokeService(database),
	}
}

func createCommands(globalCtx context.Context, services []joke.GetService, spockVoiceChns map[string]chan struct{}, guildVoiceChats map[string]voice.ChatInfo) map[string]command.DiscordSlashCommandHandler {
	welcomeCmd := command.WelcomeCommand{}
	addJokeCmd := command.AddJokeCommand{Service: services[databaseServiceID].(joke.DatabaseJokeService)}
	jokeCmd := command.JokeCommand{Services: services}
	spockCmd := command.SpockCommand{GlobalCtx: globalCtx, SpockMusicStopChs: spockVoiceChns, GuildVoiceChats: guildVoiceChats}
	stopSpockCmd := command.SpockStopCommand{SpockMusicStopChs: spockVoiceChns}

	return map[string]command.DiscordSlashCommandHandler{
		command.WelcomeCommandName:   welcomeCmd,
		command.AddJokeCommandName:   addJokeCmd,
		command.GetJokeCommandName:   jokeCmd,
		command.SpockCommandName:     spockCmd,
		command.SpockStopCommandName: stopSpockCmd,
	}
}

func createOptions(services []joke.GetService) map[string]command.DiscordEventHandler {
	apologiesOption := command.ApologiesOption{}
	nextJokeOption := command.NextJokeOption{Services: services}
	sameJokeCategoryOption := command.SameJokeCategoryOption{Services: services}

	return map[string]command.DiscordEventHandler{
		command.ApologiesButtonName:        apologiesOption,
		command.NextJokeButtonName:         nextJokeOption,
		command.SameJokeCategoryButtonName: sameJokeCategoryOption,
	}
}

func NewDiscordBot(ctx context.Context) (*DiscordBot, error) {
	prop, err := config.LoadBotVariables()
	if err != nil {
		return nil, err
	}

	bot, err := discordgo.New(fmt.Sprintf("Bot %s", prop.Token))
	if err != nil {
		return nil, fmt.Errorf("failed connect with Discord: %s", err)
	}

	// Update list of current user on voice channels
	spockVoiceChns := make(voice.SpockVoiceChannels)
	guildVoiceChats := make(map[string]voice.ChatInfo)
	vcHander := voice.ChatHandler{Ctx: ctx, SpockVoiceChns: spockVoiceChns, GuildVoiceChats: guildVoiceChats}

	bot.AddHandler(vcHander.HandleVoiceChannelUpdate)

	// Register slash commands
	database := db.NewMongodbDatabase(ctx)
	getServices := createJokeGetServices(ctx, database)
	commands := createCommands(ctx, getServices, spockVoiceChns, guildVoiceChats)
	for _, c := range commands {
		discordCmd := c.Command()
		if _, err := bot.ApplicationCommandCreate(
			prop.AppID,
			prop.ServerGUID, // If empty, command registers globally
			discordCmd,
		); err != nil {
			return nil, fmt.Errorf("registration '%s' slash command failed: %s", discordCmd.Name, err)
		}
	}

	// General handler for slash commands
	handler := slashCommandHandler{ctx, commands, createOptions(getServices)}

	bot.AddHandler(handler.handleSlashCommand)

	return &DiscordBot{
		ctx:           ctx,
		bot:           bot,
		mongodb:       database,
		spockVoiceChs: spockVoiceChns,
	}, nil
}
