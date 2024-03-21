package command

import (
	"context"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/wittano/komputer/pkgs/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
	"math/rand"
	"os"
)

const (
	idOptionKey       = "ID"
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
	Service db.JokeService
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

	res, err := j.Service.Get(ctx, searchQuery)
	if err != nil {
		return nil, err
	}

	return jokeResponse{
		username: i.Member.User.Username,
		joke:     res,
	}, nil
}

func getJokeSearchParameters(ctx context.Context, data discordgo.ApplicationCommandInteractionData) (query db.JokeSearch) {
	query.Type, query.Category = db.Single, db.Any

	for _, o := range data.Options {
		switch o.Name {
		case categoryOptionKey:
			query.Category = db.JokeCategory(o.Value.(string))
		case typeOptionKey:
			query.Type = db.JokeType(o.Value.(string))
		case idOptionKey:
			query.ID = o.Value.(primitive.ObjectID)
		default:
			slog.WarnContext(ctx, fmt.Sprintf("Invalid option for %s", o.Name))
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
				Value: db.PROGRAMMING,
			},
			{
				Name:  "Różne",
				Value: db.MISC,
			},
			{
				Name:  "Czarny humor",
				Value: db.DARK,
			},
			{
				Name:  "YoMamma",
				Value: db.YOMAMA,
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
				Value: db.Single,
			},
			{

				Name:  "Two-Part",
				Value: db.TwoPart,
			},
		},
	}
}

type jokeResponse struct {
	username string
	joke     db.Joke
}

func (j jokeResponse) Response() (msg *discordgo.InteractionResponseData) {
	switch j.joke.Type {
	case db.Single:
		msg = j.createSingleTypeJoke()
	case db.TwoPart:
		msg = j.createTwoPartJoke()
	}

	return
}

func (j jokeResponse) createSingleTypeJoke() (msg *discordgo.InteractionResponseData) {
	embeds := []*discordgo.MessageEmbed{
		{
			Type:        discordgo.EmbedTypeRich,
			Title:       "JokeContainer",
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

	if j.joke.Category == db.YOMAMA {
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
			Title: "JokeContainer",
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

	if j.joke.Category == db.YOMAMA {
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
	Service db.JokeService
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

func (a AddJokeCommand) Execute(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) (DiscordMessageReceiver, error) {
	newJoke := createJokeFromOptions(ctx, i.Data.(discordgo.ApplicationCommandInteractionData))

	if newJoke.Answer == "" {
		return nil, ErrorResponse{err: errors.New("joke: missing answer"), msg: "Kapitanie brakuje żartu"}
	}

	id, err := a.Service.Add(ctx, newJoke)
	if err != nil {
		return nil, err
	}

	return simpleMessageResponse{msg: fmt.Sprintf("BEEP BOOP. Dodałem twój żart panie Kapitanie. Jego ID to `%s`", id.Hex()), hidden: false}, nil
}

func createJokeFromOptions(ctx context.Context, data discordgo.ApplicationCommandInteractionData) (joke db.Joke) {
	for _, o := range data.Options {
		switch o.Name {
		case categoryOptionKey:
			joke.Category = db.JokeCategory(o.Value.(string))
		case typeOptionKey:
			joke.Type = db.JokeType(o.Value.(string))
		case answerOptionKey:
			joke.Answer = o.Value.(string)
		case questionOptionKey:
			joke.Question = o.Value.(string)
		default:
			slog.WarnContext(ctx, "Invalid option for %s", o.Name)
		}
	}

	return
}

type ApologiesOption struct{}

func (a ApologiesOption) Execute(_ context.Context, _ *discordgo.Session, _ *discordgo.InteractionCreate) (DiscordMessageReceiver, error) {
	return simpleMessageResponse{msg: "Przepraszam"}, nil
}

type NextJokeOption struct {
	Service db.JokeService
}

func (n NextJokeOption) Execute(ctx context.Context, _ *discordgo.Session, i *discordgo.InteractionCreate) (DiscordMessageReceiver, error) {
	joke, err := n.Service.Get(ctx, db.JokeSearch{Type: randJokeType()})
	if err != nil {
		return nil, err
	}

	return jokeResponse{i.Member.User.Username, joke}, nil
}

func randJokeType() db.JokeType {
	jokeType := db.Single
	if rand.Int()%2 == 0 {
		jokeType = db.TwoPart
	}

	return jokeType
}

type SameJokeCategoryOption struct {
	Service db.JokeService
}

func (s SameJokeCategoryOption) Execute(ctx context.Context, _ *discordgo.Session, i *discordgo.InteractionCreate) (DiscordMessageReceiver, error) {
	embedFields := i.Message.Embeds[0].Fields
	category := db.JokeCategory(embedFields[len(embedFields)-1].Value)

	joke, err := s.Service.Get(ctx, db.JokeSearch{Type: randJokeType(), Category: category})
	if err != nil {
		return nil, err
	}

	return jokeResponse{i.Member.User.Username, joke}, nil
}
