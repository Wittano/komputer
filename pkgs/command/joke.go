package command

import (
	"context"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/wittano/komputer/pkgs/joke"
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
	Services []joke.GetService
}

func (j JokeCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        GetJokeCommandName,
		Description: "Tell me some joke",
		GuildID:     os.Getenv("SERVER_GUID"),
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			getJokeCategoryOption(false),
			getJokeTypeOption(false),
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
	searchQuery := getJokeSearchParameters(ctx, i.Data.(discordgo.ApplicationCommandInteractionData))

findJoke:
	select {
	case <-ctx.Done():
		return nil, context.Canceled
	default:
	}

	service, err := selectGetService(ctx, j.Services)
	if err != nil {
		return nil, ErrorResponse{err, "Nie udało mi się, znaleść żadnego żartu"}
	}

	res, err := service.Get(ctx, searchQuery)
	if err != nil {
		slog.With(requestIDKey, ctx.Value(requestIDKey)).ErrorContext(ctx, err.Error())
		goto findJoke
	}

	return jokeResponse{
		username: i.Member.User.Username,
		joke:     res,
	}, nil
}

func selectGetService(ctx context.Context, getServices []joke.GetService) (joke.GetService, error) {
	if len(getServices) == 1 {
		service := getServices[0]

		if activeService, ok := service.(joke.ActiveService); ok && !activeService.Active(ctx) {
			return nil, errors.New("all joke services is disabled")
		}

		return service, nil
	}

	i := rand.Int() % len(getServices)
	service := getServices[uint8(i)]
	if service == nil {
		getServices = slices.Delete(getServices, i, i+1)
		return selectGetService(ctx, getServices)
	}

	if activeService, ok := service.(joke.ActiveService); ok && !activeService.Active(ctx) {
		getServices = slices.Delete(getServices, i, i+1)
		return selectGetService(ctx, getServices)
	}

	return service, nil
}

func getJokeSearchParameters(ctx context.Context, data discordgo.ApplicationCommandInteractionData) (query joke.SearchParameters) {
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
			slog.With(requestIDKey, ctx.Value(requestIDKey)).WarnContext(ctx, fmt.Sprintf("Invalid option for %s", o.Name))
		}
	}

	return
}

func getJokeCategoryOption(required bool) *discordgo.ApplicationCommandOption {
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

func getJokeTypeOption(required bool) *discordgo.ApplicationCommandOption {
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

type jokeResponse struct {
	username string
	joke     joke.Joke
}

func (j jokeResponse) Response() (msg *discordgo.InteractionResponseData) {
	switch j.joke.Type {
	case joke.Single:
		msg = j.createSingleTypeJoke()
	case joke.TwoPart:
		msg = j.createTwoPartJoke()
	}

	return
}

func (j jokeResponse) createSingleTypeJoke() (msg *discordgo.InteractionResponseData) {
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
		Components: createButtonReactions(),
		Embeds:     embeds,
	}
}

func (j jokeResponse) createTwoPartJoke() *discordgo.InteractionResponseData {
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
		Components: createButtonReactions(),
		Embeds:     embeds,
	}
}

func createButtonReactions() []discordgo.MessageComponent {
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
			getJokeCategoryOption(true),
			getJokeTypeOption(true),
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
	newJoke := createJokeFromOptions(i.Data.(discordgo.ApplicationCommandInteractionData))

	if newJoke.Answer == "" {
		return nil, ErrorResponse{Err: errors.New("joke: missing answer"), Msg: "Zrujnowałeś ten żart, Panie Kapitanie"}
	}

	newJoke.ID = primitive.NewObjectID()

	id, err := a.Service.Add(ctx, newJoke)
	if err != nil {
		return nil, err
	}

	return SimpleMessageResponse{Msg: fmt.Sprintf("BEEP BOOP. Dodałem twój żart panie Kapitanie. Jego ID to `%s`", id), Hidden: true}, nil
}

func createJokeFromOptions(data discordgo.ApplicationCommandInteractionData) (j joke.Joke) {
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

func (a ApologiesOption) Execute(_ context.Context, _ *discordgo.Session, _ *discordgo.InteractionCreate) (DiscordMessageReceiver, error) {
	return SimpleMessageResponse{Msg: "Przepraszam"}, nil
}

type NextJokeOption struct {
	Services []joke.GetService
}

func (n NextJokeOption) Execute(ctx context.Context, _ *discordgo.Session, i *discordgo.InteractionCreate) (DiscordMessageReceiver, error) {
	service, err := selectGetService(ctx, n.Services)
	if err != nil {
		return nil, err
	}

	res, err := service.Get(ctx, joke.SearchParameters{Type: randJokeType()})
	if err != nil {
		return nil, err
	}

	return jokeResponse{i.Member.User.Username, res}, nil
}

func randJokeType() joke.Type {
	jokeType := joke.Single
	if rand.Int()%2 == 0 {
		jokeType = joke.TwoPart
	}

	return jokeType
}

type SameJokeCategoryOption struct {
	Services []joke.GetService
}

func (s SameJokeCategoryOption) Execute(ctx context.Context, _ *discordgo.Session, i *discordgo.InteractionCreate) (DiscordMessageReceiver, error) {
	embedFields := i.Message.Embeds[0].Fields
	category := joke.Category(embedFields[len(embedFields)-1].Value)

	service, err := selectGetService(ctx, s.Services)
	if err != nil {
		return nil, err
	}

	res, err := service.Get(ctx, joke.SearchParameters{Type: randJokeType(), Category: category})
	if err != nil {
		return nil, err
	}

	return jokeResponse{i.Member.User.Username, res}, nil
}
