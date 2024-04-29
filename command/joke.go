package command

import (
	"context"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/wittano/komputer/joke"
	"github.com/wittano/komputer/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
	"math/rand"
	"os"
	"slices"
)

const (
	idOptionKey       = "id"
	categoryOptionKey = "category"
	typeOptionKey     = "type"
	questionOptionKey = "question"
	answerOptionKey   = "answer"
)

const (
	AddJokeCommandName = "add-joke"
	GetJokeCommandName = "joke"
)

const (
	NextJokeButtonName         = "nextJokeButtonId"
	SameJokeCategoryButtonName = "sameJokeButtonId"
	ApologiesButtonName        = "apologiesButtonId"
)

type JokeCommand struct {
	Services []joke.SearchService
}

func (j JokeCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        GetJokeCommandName,
		Description: "Tell me some joke",
		GuildID:     os.Getenv(serverGuildKey),
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			jokeCategoryOption(false),
			jokeTypeOption(false),
			{
				Name:        idOptionKey,
				Description: "Joke ID",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    false,
			},
		},
	}
}

func (j JokeCommand) Execute(ctx context.Context, _ *discordgo.Session, i *discordgo.InteractionCreate) (DiscordMessageReceiver, error) {
	searchQuery := searchParams(ctx, i.Data.(discordgo.ApplicationCommandInteractionData))

	loggerCtx := ctx.(log.Context)

findJoke:
	select {
	case <-ctx.Done():
		return nil, context.Canceled
	default:
	}

	service, err := findService(ctx, j.Services)
	if err != nil {
		return nil, DiscordError{err, "Nie udało mi się, znaleść żadnego żartu"}
	}

	res, err := service.Joke(ctx, searchQuery)
	if err != nil {
		loggerCtx.Logger.Error(err.Error())
		goto findJoke
	}

	return discordJoke{
		username: i.Member.User.Username,
		joke:     res,
	}, nil
}

func findService(ctx context.Context, services []joke.SearchService) (joke.SearchService, error) {
	if len(services) == 1 {
		service := services[0]

		if checker, ok := service.(joke.ActiveChecker); ok && !checker.Active(ctx) {
			return nil, errors.New("all joke services is disabled")
		}

		return service, nil
	}

	i := rand.Int() % len(services)
	service := services[uint8(i)]
	if service == nil {
		services = slices.Delete(services, i, i+1)
		return findService(ctx, services)
	}

	if activeService, ok := service.(joke.ActiveChecker); ok && !activeService.Active(ctx) {
		services = slices.Delete(services, i, i+1)
		return findService(ctx, services)
	}

	return service, nil
}

// Get joke.SearchParams from Discord options
func searchParams(ctx context.Context, data discordgo.ApplicationCommandInteractionData) (query joke.SearchParams) {
	query.Type, query.Category = joke.Single, joke.Any

	for _, o := range data.Options {
		switch o.Name {
		case categoryOptionKey:
			query.Category = joke.Category(o.Value.(string))
		case typeOptionKey:
			query.Type = joke.Type(o.Value.(string))
		case idOptionKey:
			query.ID = o.Value.(primitive.ObjectID)
		default:
			log.Log(ctx, func(log slog.Logger) {
				log.Warn(fmt.Sprintf("Invalid searchOption for %s", o.Name))
			})
		}
	}

	return
}

func jokeCategoryOption(required bool) *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        categoryOptionKey,
		Description: "Joke category",
		Required:    required,
		Choices: []*discordgo.ApplicationCommandOptionChoice{
			{
				Name:  "Programowanie",
				Value: joke.PROGRAMMING,
			},
			{
				Name:  "Różne",
				Value: joke.MISC,
			},
			{
				Name:  "Czarny humor",
				Value: joke.DARK,
			},
			{
				Name:  "YoMamma",
				Value: joke.YOMAMA,
			},
		},
	}
}

func jokeTypeOption(required bool) *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        typeOptionKey,
		Description: "Type of joke",
		Required:    required,
		Choices: []*discordgo.ApplicationCommandOptionChoice{
			{
				Name:  "Single",
				Value: joke.Single,
			},
			{

				Name:  "Two-Part",
				Value: joke.TwoPart,
			},
		},
	}
}

type discordJoke struct {
	username string
	joke     joke.Joke
}

func (j discordJoke) Response() (msg *discordgo.InteractionResponseData) {
	switch j.joke.Type {
	case joke.Single:
		msg = j.singleTypeJoke()
	case joke.TwoPart:
		msg = j.twoPartJoke()
	}

	return
}

func (j discordJoke) singleTypeJoke() (msg *discordgo.InteractionResponseData) {
	embeds := []*discordgo.MessageEmbed{
		{
			Type:        discordgo.EmbedTypeRich,
			Title:       "Joke",
			Description: j.joke.Answer,
			Color:       0x02f5f5,
			Author: &discordgo.MessageEmbedAuthor{
				Name: "komputer",
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "category",
					Value: string(j.joke.Category),
				},
			},
		},
	}

	if j.joke.Category == joke.YOMAMA {
		const muscleManGifURL = "https://media.tenor.com/sgS8GdoZGn8AAAAd/muscle-man-regular-show-muscle-man.gif"

		embeds[0].Image = &discordgo.MessageEmbedImage{
			URL: muscleManGifURL,
		}
	}

	return &discordgo.InteractionResponseData{
		Content:    fmt.Sprintf("BEEP BOOP, Tak jest kapitanie %s!", j.username),
		Components: buttonReactions(),
		Embeds:     embeds,
	}
}

func (j discordJoke) twoPartJoke() *discordgo.InteractionResponseData {
	embeds := []*discordgo.MessageEmbed{
		{
			Type:  discordgo.EmbedTypeRich,
			Title: "Joke",
			Color: 0x02f5f5,
			Author: &discordgo.MessageEmbedAuthor{
				Name: "komputer",
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Question",
					Value:  j.joke.Question,
					Inline: true,
				},
				{
					Name:   "Answer",
					Value:  j.joke.Answer,
					Inline: true,
				},
				{
					Name:  "category",
					Value: string(j.joke.Category),
				},
			},
		},
	}

	if j.joke.Category == joke.YOMAMA {
		embeds[0].Image = &discordgo.MessageEmbedImage{
			URL: "https://media.tenor.com/sgS8GdoZGn8AAAAd/muscle-man-regular-show-muscle-man.gif",
		}
	}

	return &discordgo.InteractionResponseData{
		Content:    fmt.Sprintf("BEEP BOOP, Tak jest Panie kapitanie %s!", j.username),
		Components: buttonReactions(),
		Embeds:     embeds,
	}
}

func buttonReactions() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Style:    discordgo.PrimaryButton,
					Label:    "Przeproś",
					CustomID: ApologiesButtonName,
				},
				discordgo.Button{
					Style:    discordgo.SecondaryButton,
					Label:    "Zabawne powiedz coś podobnego",
					CustomID: SameJokeCategoryButtonName,
				},
				discordgo.Button{
					Style:    discordgo.SecondaryButton,
					Label:    "Zabawne powiedz więcej",
					CustomID: NextJokeButtonName,
				},
			},
		},
	}
}

type AddJokeCommand struct {
	Service joke.AddService
}

func (a AddJokeCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        AddJokeCommandName,
		Description: "Add new joke to server database",
		GuildID:     os.Getenv("SERVER_GUID"),
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			jokeCategoryOption(true),
			jokeTypeOption(true),
			{
				Name:        answerOptionKey,
				Description: "Main part of joke",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
			{
				Name:        questionOptionKey,
				Description: "Question part in two-part joke",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    false,
			},
		},
	}
}

func (a AddJokeCommand) Execute(ctx context.Context, _ *discordgo.Session, i *discordgo.InteractionCreate) (DiscordMessageReceiver, error) {
	newJoke := jokeFromOptions(i.Data.(discordgo.ApplicationCommandInteractionData))

	if newJoke.Answer == "" {
		return nil, DiscordError{Err: errors.New("joke: missing answer"), Msg: "Zrujnowałeś ten żart, Panie Kapitanie"}
	}

	newJoke.ID = primitive.NewObjectID()

	id, err := a.Service.Add(ctx, newJoke)
	if err != nil {
		return nil, err
	}

	return SimpleMessage{Msg: fmt.Sprintf("BEEP BOOP. Dodałem twój żart panie Kapitanie. Jego ID to `%s`", id), Hidden: true}, nil
}

func jokeFromOptions(data discordgo.ApplicationCommandInteractionData) (j joke.Joke) {
	for _, o := range data.Options {
		switch o.Name {
		case categoryOptionKey:
			j.Category = joke.Category(o.Value.(string))
		case typeOptionKey:
			j.Type = joke.Type(o.Value.(string))
		case answerOptionKey:
			j.Answer = o.Value.(string)
		case questionOptionKey:
			j.Question = o.Value.(string)
		}
	}

	return
}

type ApologiesOption struct{}

func (a ApologiesOption) Match(customID string) bool {
	return customID == ApologiesButtonName
}

func (a ApologiesOption) Execute(_ context.Context, _ *discordgo.Session, _ *discordgo.InteractionCreate) (DiscordMessageReceiver, error) {
	return SimpleMessage{Msg: "Przepraszam"}, nil
}

type NextJokeOption struct {
	Services []joke.SearchService
}

func (n NextJokeOption) Match(customID string) bool {
	return customID == NextJokeButtonName
}

func (n NextJokeOption) Execute(ctx context.Context, _ *discordgo.Session, i *discordgo.InteractionCreate) (DiscordMessageReceiver, error) {
	service, err := findService(ctx, n.Services)
	if err != nil {
		return nil, err
	}

	res, err := service.Joke(ctx, joke.SearchParams{Type: randJokeType()})
	if err != nil {
		return nil, err
	}

	return discordJoke{i.Member.User.Username, res}, nil
}

func randJokeType() joke.Type {
	jokeType := joke.Single
	if rand.Int()%2 == 0 {
		jokeType = joke.TwoPart
	}

	return jokeType
}

type SameJokeCategoryOption struct {
	Services []joke.SearchService
}

func (s SameJokeCategoryOption) Match(customID string) bool {
	return customID == SameJokeCategoryButtonName
}

func (s SameJokeCategoryOption) Execute(ctx context.Context, _ *discordgo.Session, i *discordgo.InteractionCreate) (DiscordMessageReceiver, error) {
	fields := i.Message.Embeds[0].Fields
	category := joke.Category(fields[len(fields)-1].Value)

	service, err := findService(ctx, s.Services)
	if err != nil {
		return nil, err
	}

	res, err := service.Joke(ctx, joke.SearchParams{Type: randJokeType(), Category: category})
	if err != nil {
		return nil, err
	}

	return discordJoke{i.Member.User.Username, res}, nil
}
