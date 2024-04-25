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
	"regexp"
	"strconv"
	"strings"
)

const (
	audioNameOption = "name"
	audioIdOption   = "id"

	nextIdsButtonID     = "nextIds"
	previousIdsButtonID = "previousIds"
)

const ListCommandName = "list"

type pageContentDirection int8

const (
	left  pageContentDirection = -1
	none  pageContentDirection = 0
	right pageContentDirection = 1
)

type buttonPosition uint8

const (
	previous buttonPosition = 0
	next     buttonPosition = 1
)

type listContentResponse struct {
	content []api.AudioFileInfo
	pattern voice.AudioSearch
}

func (l listContentResponse) Response() *discordgo.InteractionResponseData {
	if len(l.content) == 0 {
		return SimpleMessageResponse{Msg: "Nie znalazłem żadnego takiego pliku"}.Response()
	}

	var msg strings.Builder

	for _, info := range l.content {
		msg.WriteString(fmt.Sprintf("- %s\n", info.String()))
	}

	const customIDFormat = "%s_%d_%s"
	return &discordgo.InteractionResponseData{
		Title: "List of Audios IDs",
		Flags: discordgo.MessageFlagsEphemeral,
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
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Style:    discordgo.SecondaryButton,
						Label:    "Previous",
						CustomID: fmt.Sprintf(customIDFormat, previousIdsButtonID, l.pattern.Type, l.pattern.Value),
					},
					discordgo.Button{
						Style:    discordgo.PrimaryButton,
						Label:    "Next",
						CustomID: fmt.Sprintf(customIDFormat, nextIdsButtonID, l.pattern.Type, l.pattern.Value),
					},
				},
			},
		},
	}
}

type NextListCommandOption struct {
	Cmd *ListCommand
}

func (n NextListCommandOption) MatchCustomID(customID string) bool {
	reg := regexp.MustCompile(fmt.Sprintf("^%s_([0-1])_(a-z0-9)*", nextIdsButtonID))

	return reg.MatchString(customID)
}

func (n NextListCommandOption) Execute(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) (DiscordMessageReceiver, error) {
	select {
	case <-ctx.Done():
		return nil, context.Canceled
	default:
	}

	customerID, err := getCustomIDFromOptionIntegration(*i, next)
	if err != nil {
		return nil, err
	}

	option, err := getOptionFromCustomID(customerID)
	if err != nil {
		return nil, err
	}

	userID := i.Member.User.ID
	result, err := n.Cmd.audioFileInfo(ctx, userID, option)
	if err != nil {
		return nil, err
	}

	n.Cmd.updatePageCounter(len(result), userID, right)

	return listContentResponse{result, option}, nil
}

type PreviousListCommandOption struct {
	Cmd *ListCommand
}

func (p PreviousListCommandOption) MatchCustomID(customID string) bool {
	reg := regexp.MustCompile(fmt.Sprintf("^%s_([0-1])_(a-z0-9)*", previousIdsButtonID))

	return reg.MatchString(customID)
}

func (p PreviousListCommandOption) Execute(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) (DiscordMessageReceiver, error) {
	select {
	case <-ctx.Done():
		return nil, context.Canceled
	default:
	}

	customerID, err := getCustomIDFromOptionIntegration(*i, previous)
	if err != nil {
		return nil, err
	}

	option, err := getOptionFromCustomID(customerID)
	if err != nil {
		return nil, err
	}

	userID := i.Member.User.ID
	result, err := p.Cmd.audioFileInfo(ctx, userID, option)
	if err != nil {
		return nil, err
	}

	p.Cmd.updatePageCounter(len(result), userID, left)

	return listContentResponse{result, option}, nil
}

type ListCommand struct {
	// TODO Clean up pageCounter after a few minutes
	pageCounter map[string]int
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
				MaxLength:   80,
				Type:        discordgo.ApplicationCommandOptionString,
			},
			{
				Name:        audioIdOption,
				Description: "Part of audio ID",
				Required:    false,
				MaxLength:   80,
				Type:        discordgo.ApplicationCommandOptionString,
			},
		},
	}
}

func (l ListCommand) Execute(ctx context.Context, _ *discordgo.Session, i *discordgo.InteractionCreate) (DiscordMessageReceiver, error) {
	userID := i.Member.User.ID
	option := getOption(i.Data.(discordgo.ApplicationCommandInteractionData))
	result, err := l.audioFileInfo(ctx, userID, option)
	if err != nil {
		return nil, err
	}

	l.updatePageCounter(len(result), userID, none)

	return listContentResponse{result, option}, nil
}

func (l *ListCommand) updatePageCounter(resultSize int, userID string, direction pageContentDirection) {
	if _, ok := l.pageCounter[userID]; !ok {
		l.pageCounter[userID] = 0
	}

	if resultSize == 0 || l.pageCounter[userID] < 0 {
		l.pageCounter[userID] = 0
	} else if resultSize != 0 && l.pageCounter[userID] >= 0 {
		l.pageCounter[userID] += int(direction)
	}
}

func (l ListCommand) audioFileInfo(
	ctx context.Context,
	userID string,
	option voice.AudioSearch,
) (result []api.AudioFileInfo, err error) {
	for _, service := range l.services {
		if !service.IsActive() {
			continue
		}

		result, err = service.SearchAudio(ctx, option, uint(l.pageCounter[userID]))
		if err == nil {
			break
		} else {
			slog.With(requestIDKey, ctx.Value(requestIDKey)).WarnContext(ctx, err.Error())
		}
	}

	return result, err
}

func getOption(data discordgo.ApplicationCommandInteractionData) (query voice.AudioSearch) {
	query = voice.AudioSearch{Type: voice.IDType}

	for _, o := range data.Options {
		switch o.Name {
		case audioNameOption:
			query.Type = voice.NameType
		case audioIdOption:
			query.Type = voice.IDType
		default:
			continue
		}

		query.Value = o.Value.(string)

		return
	}

	return
}

func getOptionFromCustomID(customID string) (query voice.AudioSearch, err error) {
	data := strings.Split(customID, "_")[1:]
	query.Value = data[1]

	typeNum, err := strconv.Atoi(data[0])
	if err != nil {
		return
	}

	query.Type = voice.AudioQueryType(typeNum)

	return
}

func getCustomIDFromOptionIntegration(i discordgo.InteractionCreate, buttonPosition buttonPosition) (string, error) {
	actionRow, ok := i.Message.Components[0].(*discordgo.ActionsRow)
	if !ok {
		return "", errors.New("failed cast component to *discordgo.ActionsRow")
	}

	button, ok := actionRow.Components[buttonPosition].(*discordgo.Button)
	if !ok {
		return "", errors.New("failed cast component to *discordgo.Button")
	}

	return button.CustomID, nil
}

func NewListCommand(services ...voice.AudioSearchService) *ListCommand {
	return &ListCommand{
		make(map[string]int),
		services,
	}
}
