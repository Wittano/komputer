package command

import (
	"context"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/wittano/komputer/bot/log"
	"github.com/wittano/komputer/internal"
	"github.com/wittano/komputer/internal/joke"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
	"os"
	"slices"
)

const (
	idOptionKey       = "id"
	categoryOptionKey = "category"
	typeOptionKey     = "type"
)

const (
	GetJokeCommandName = "joke"
)

const (
	NextJokeButtonName         = "nextJokeButtonId"
	SameJokeCategoryButtonName = "sameJokeButtonId"
	ApologiesButtonName        = "apologiesButtonId"
)

type JokeCommand struct {
	Services []internal.SearchService
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
				Description: "Jokes ID",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    false,
			},
		},
	}
}

func (j JokeCommand) Execute(ctx log.Context, _ *discordgo.Session, i *discordgo.InteractionCreate) (DiscordMessageReceiver, error) {
	searchQuery := searchParams(ctx, i.Data.(discordgo.ApplicationCommandInteractionData))

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

	res, err := service.RandomJoke(ctx, searchQuery)
	if err != nil {
		ctx.Logger.Error(err.Error())
		goto findJoke
	}

	return discordJoke{
		username: i.Member.User.Username,
		joke:     res,
	}, nil
}

func findService(ctx log.Context, services []internal.SearchService) (internal.SearchService, error) {
	if len(services) == 1 {
		service := services[0]

		if checker, ok := service.(internal.ActiveChecker); ok && !checker.Active(ctx) {
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

	if activeService, ok := service.(internal.ActiveChecker); ok && !activeService.Active(ctx) {
		services = slices.Delete(services, i, i+1)
		return findService(ctx, services)
	}

	return service, nil
}

// Get internal.SearchParams from Discord options
func searchParams(ctx log.Context, data discordgo.ApplicationCommandInteractionData) (query internal.SearchParams) {
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
			ctx.Logger.Warn(fmt.Sprintf("Invalid searchOption for %s", o.Name))
		}
	}

	return
}

func jokeCategoryOption(required bool) *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        categoryOptionKey,
		Description: "Jokes category",
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
		Description: "RawType of joke",
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
	joke     joke.DbModel
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
			Title:       "Jokes",
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
			Title: "Jokes",
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

type ApologiesOption struct{}

func (a ApologiesOption) Match(customID string) bool {
	return customID == ApologiesButtonName
}

func (a ApologiesOption) Execute(_ log.Context, _ *discordgo.Session, _ *discordgo.InteractionCreate) (DiscordMessageReceiver, error) {
	return SimpleMessage{Msg: "Przepraszam"}, nil
}

type NextJokeOption struct {
	Services []internal.SearchService
}

func (n NextJokeOption) Match(customID string) bool {
	return customID == NextJokeButtonName
}

func (n NextJokeOption) Execute(ctx log.Context, _ *discordgo.Session, i *discordgo.InteractionCreate) (DiscordMessageReceiver, error) {
	service, err := findService(ctx, n.Services)
	if err != nil {
		return nil, err
	}

	res, err := service.RandomJoke(ctx, internal.SearchParams{Type: randJokeType()})
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
	Services []internal.SearchService
}

func (s SameJokeCategoryOption) Match(customID string) bool {
	return customID == SameJokeCategoryButtonName
}

func (s SameJokeCategoryOption) Execute(ctx log.Context, _ *discordgo.Session, i *discordgo.InteractionCreate) (DiscordMessageReceiver, error) {
	fields := i.Message.Embeds[0].Fields
	category := joke.Category(fields[len(fields)-1].Value)

	service, err := findService(ctx, s.Services)
	if err != nil {
		return nil, err
	}

	res, err := service.RandomJoke(ctx, internal.SearchParams{Type: randJokeType(), Category: category})
	if err != nil {
		return nil, err
	}

	return discordJoke{i.Member.User.Username, res}, nil
}
