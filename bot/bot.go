package bot

import (
	"context"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/wittano/komputer/bot/command"
	"github.com/wittano/komputer/bot/config"
	"github.com/wittano/komputer/bot/joke"
	"github.com/wittano/komputer/bot/log"
	"github.com/wittano/komputer/bot/voice"
	"github.com/wittano/komputer/db"
	dbJoke "github.com/wittano/komputer/db/joke"
	"log/slog"
	"time"
)

const (
	jokeDevServiceID  = 0
	humorAPIServiceID = 1
	databaseServiceID = 2
)

type slashCommandHandler struct {
	ctx      context.Context
	commands map[string]command.DiscordSlashCommandHandler
	options  []command.DiscordEventHandler
}

func (sc slashCommandHandler) handleSlashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	const cmdTimeout = time.Second * 2

	requestID := uuid.New().String()
	loggerCtx := log.NewContext(requestID)
	deadlineCtx, cancel := context.WithTimeout(loggerCtx, cmdTimeout)
	ctx := context.WithValue(deadlineCtx, dbJoke.GuildIDKey, i.GuildID)
	defer cancel()

	userID := i.Member.User.ID
	// Handle options assigned to slash commands
	if i.Type == discordgo.InteractionMessageComponent {
		customID := i.Data.(discordgo.MessageComponentInteractionData).CustomID

		for _, option := range sc.options {
			if matcher, ok := option.(command.DiscordOptionMatcher); ok && matcher.Match(customID) {
				log.Log(ctx, func(log slog.Logger) {
					log.InfoContext(ctx, fmt.Sprintf("user '%s' select '%s' option", userID, customID))
				})

				handleEventResponse(ctx, s, i, option)

				return
			}
		}
	}

	defer func() {
		if r := recover(); r != nil {
			cancel()
			log.Log(ctx, func(log slog.Logger) {
				log.Error("unexpected error during handle command", r)
			})
		}
	}()

	// Handle slash commands
	name := i.ApplicationCommandData().Name
	if cmd, ok := sc.commands[name]; ok {
		log.Log(ctx, func(l slog.Logger) {
			l.InfoContext(ctx, fmt.Sprintf("user '%s' execute slash command '%s'", userID, name))
		})

		handleEventResponse(ctx, s, i, cmd)
	} else {
		msg := command.SimpleMessage{Msg: "Kapitanie co chcesz zrobić?", Hidden: true}

		log.Log(ctx, func(l slog.Logger) {
			l.WarnContext(ctx, "someone try execute unknown command")
		})
		command.CreateDiscordInteractionResponse(ctx, i, s, msg)
	}
}

func handleEventResponse(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, event command.DiscordEventHandler) {
	msg, err := event.Execute(ctx, s, i)

	if errors.Is(err, command.DiscordError{}) {
		errors.As(err, &msg)
	} else if err != nil {
		ctx.(log.Context).Logger.ErrorContext(ctx, err.Error())

		msg = command.DiscordError{
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

func createJokeGetServices(globalCtx context.Context, database *db.MongodbDatabase) []dbJoke.SearchService {
	return []dbJoke.SearchService{
		jokeDevServiceID:  joke.NewJokeDevService(globalCtx),
		humorAPIServiceID: joke.NewHumorAPIService(globalCtx),
		databaseServiceID: joke.NewDatabaseJokeService(database),
	}
}

func createCommands(
	globalCtx context.Context,
	services []dbJoke.SearchService,
	spockVoice map[string]chan struct{},
	guildVoiceChats map[string]voice.ChatInfo,
) map[string]command.DiscordSlashCommandHandler {
	welcome := command.WelcomeCommand{}
	addJoke := command.AddJokeCommand{Service: services[databaseServiceID].(dbJoke.Database)}
	getJoke := command.JokeCommand{Services: services}
	spock := command.SpockCommand{
		GlobalCtx:       globalCtx,
		MusicStopChs:    spockVoice,
		GuildVoiceChats: guildVoiceChats,
	}
	stop := command.StopCommand{MusicStopChs: spockVoice}
	list := command.NewListCommand()

	return map[string]command.DiscordSlashCommandHandler{
		command.WelcomeCommandName: welcome,
		command.AddJokeCommandName: addJoke,
		command.GetJokeCommandName: getJoke,
		command.SpockCommandName:   spock,
		command.StopCommandName:    stop,
		command.ListCommandName:    list,
	}
}

func createOptions(
	services []dbJoke.SearchService,
	commands map[string]command.DiscordSlashCommandHandler,
) []command.DiscordEventHandler {
	apologies := command.ApologiesOption{}
	nextJoke := command.NextJokeOption{Services: services}
	sameJokeCategory := command.SameJokeCategoryOption{Services: services}

	listCommand := commands[command.ListCommandName].(*command.ListCommand)
	nextList := command.NextListCommandOption{CMD: listCommand}
	previousList := command.PreviousListCommandOption{Cmd: listCommand}

	return []command.DiscordEventHandler{
		apologies,
		nextJoke,
		sameJokeCategory,
		nextList,
		previousList,
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
	vcHandler := voice.ChatHandler{Ctx: ctx, SpockVoiceChns: spockVoiceChns, GuildVoiceChats: guildVoiceChats}

	bot.AddHandler(vcHandler.HandleVoiceChannelUpdate)

	// Register slash commands
	database := db.Mongodb(ctx)
	services := createJokeGetServices(ctx, database)
	commands := createCommands(ctx, services, spockVoiceChns, guildVoiceChats)
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
	handler := slashCommandHandler{ctx, commands, createOptions(services, commands)}

	bot.AddHandler(handler.handleSlashCommand)

	return &DiscordBot{
		ctx:           ctx,
		bot:           bot,
		mongodb:       database,
		spockVoiceChs: spockVoiceChns,
	}, nil
}
