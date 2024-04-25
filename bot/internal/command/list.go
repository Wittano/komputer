package command

import (
	"context"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/wittano/komputer/api"
	"github.com/wittano/komputer/bot/internal/voice"
	"log/slog"
	"os"
	"strings"
)

const (
	audioNameOption = "name"
	audioIdOption   = "id"
)

const ListCommandName = "list"

type listContentResponse struct {
	content []api.AudioFileInfo
}

func (l listContentResponse) Response() *discordgo.InteractionResponseData {
	if len(l.content) == 0 {
		return SimpleMessageResponse{Msg: "Nie znalazłem żadnego takiego pliku"}.Response()
	}

	var msg strings.Builder

	for _, info := range l.content {
		msg.WriteString(fmt.Sprintf("- %s\n", info.String()))
	}

	return &discordgo.InteractionResponseData{
		Title: "List of Audios IDs",
		Embeds: []*discordgo.MessageEmbed{
			{
				Type:  discordgo.EmbedTypeRich,
				Title: "List of audio names and ids",
				Color: 0x02f5f5,
				Author: &discordgo.MessageEmbedAuthor{
					Name: "komputer",
				},
				Description: msg.String(),
			},
		},
	}
}

type ListCommand struct {
	// TODO Clean up pageCounter after a few minutes
	pageCounter map[string]uint
	services    []voice.AudioSearchService
}

func (l ListCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        ListCommandName,
		Description: "Show list of available audios",
		Type:        discordgo.ChatApplicationCommand,
		GuildID:     os.Getenv(serverGuildKey),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        audioNameOption,
				Description: "Original audio name. Suffix .mp3 isn't necessary",
				Required:    false,
				Type:        discordgo.ApplicationCommandOptionString,
			},
			{
				Name:        audioIdOption,
				Description: "Part of audio ID",
				Required:    false,
				Type:        discordgo.ApplicationCommandOptionString,
			},
		},
	}
}

func (l ListCommand) Execute(ctx context.Context, _ *discordgo.Session, i *discordgo.InteractionCreate) (DiscordMessageReceiver, error) {
	option, err := getOption(i.Data.(discordgo.ApplicationCommandInteractionData))
	if err != nil {
		return nil, err
	}

	userID := i.Member.User.ID
	page, _ := l.pageCounter[userID]
	var result []api.AudioFileInfo

	for _, service := range l.services {
		if !service.IsActive() {
			continue
		}

		result, err = service.SearchAudio(ctx, option, page)
		if err == nil {
			break
		} else {
			slog.With(requestIDKey, ctx.Value(requestIDKey)).WarnContext(ctx, err.Error())
		}
	}

	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		l.pageCounter[userID] = 0
	} else {
		l.pageCounter[userID] += 1
	}

	return listContentResponse{result}, nil
}

func getOption(data discordgo.ApplicationCommandInteractionData) (voice.AudioSearch, error) {
	for _, o := range data.Options {
		var audioType voice.AudioQueryType

		switch o.Name {
		case audioNameOption:
			audioType = voice.NameType
		case audioIdOption:
			audioType = voice.IDType
		default:
			continue
		}

		return voice.AudioSearch{
			Type:  audioType,
			Value: o.Value.(string),
		}, nil
	}

	return voice.AudioSearch{}, errors.New("unknown option")
}

func NewListCommand(services ...voice.AudioSearchService) ListCommand {
	return ListCommand{
		make(map[string]uint),
		services,
	}
}
