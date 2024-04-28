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

	nextIdsButtonID     = "nextIDs"
	previousIdsButtonID = "previousIDs"
)

const ListCommandName = "list"

type pageDirection int8

const (
	left  pageDirection = -1
	none  pageDirection = 0
	right pageDirection = 1
)

type buttonPosition uint8

const (
	previous buttonPosition = 0
	next     buttonPosition = 1
)

type listContent struct {
	content []api.AudioFileInfo
	pattern voice.SearchParams
}

func (l listContent) Response() *discordgo.InteractionResponseData {
	if len(l.content) == 0 {
		return SimpleMessage{Msg: "Nie znalazłem żadnego takiego pliku"}.Response()
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
	CMD *ListCommand
}

func (n NextListCommandOption) Match(customID string) bool {
	return regexp.MustCompile(fmt.Sprintf("^%s_([0-1])_(a-z0-9)*", nextIdsButtonID)).MatchString(customID)
}

func (n NextListCommandOption) Execute(ctx context.Context, _ *discordgo.Session, i *discordgo.InteractionCreate) (DiscordMessageReceiver, error) {
	select {
	case <-ctx.Done():
		return nil, context.Canceled
	default:
	}

	customID, err := CustomID(*i, next)
	if err != nil {
		return nil, err
	}

	option, err := paramsFromCustomID(customID)
	if err != nil {
		return nil, err
	}

	userID := i.Member.User.ID
	result, err := n.CMD.audioFileInfo(ctx, userID, option)
	if err != nil {
		return nil, err
	}

	n.CMD.update(len(result), userID, right)

	return listContent{result, option}, nil
}

type PreviousListCommandOption struct {
	Cmd *ListCommand
}

func (p PreviousListCommandOption) Match(customID string) bool {
	return regexp.MustCompile(fmt.Sprintf("^%s_([0-1])_(a-z0-9)*", previousIdsButtonID)).MatchString(customID)
}

func (p PreviousListCommandOption) Execute(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) (DiscordMessageReceiver, error) {
	select {
	case <-ctx.Done():
		return nil, context.Canceled
	default:
	}

	customerID, err := CustomID(*i, previous)
	if err != nil {
		return nil, err
	}

	option, err := paramsFromCustomID(customerID)
	if err != nil {
		return nil, err
	}

	userID := i.Member.User.ID
	result, err := p.Cmd.audioFileInfo(ctx, userID, option)
	if err != nil {
		return nil, err
	}

	p.Cmd.update(len(result), userID, left)

	return listContent{result, option}, nil
}

type ListCommand struct {
	// TODO Clean up pageCounter after a few minutes

	// map of pages number, which user should see
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

func (l ListCommand) Execute(
	ctx context.Context,
	_ *discordgo.Session,
	i *discordgo.InteractionCreate,
) (DiscordMessageReceiver, error) {
	userID := i.Member.User.ID
	option := paramsFromDiscord(i.Data.(discordgo.ApplicationCommandInteractionData))
	result, err := l.audioFileInfo(ctx, userID, option)
	if err != nil {
		return nil, err
	}

	l.update(len(result), userID, none)

	return listContent{result, option}, nil
}

// Update pageCounter
func (l *ListCommand) update(size int, userID string, direction pageDirection) {
	if _, ok := l.pageCounter[userID]; !ok {
		l.pageCounter[userID] = 0
	}

	if size == 0 || l.pageCounter[userID] < 0 {
		l.pageCounter[userID] = 0
	} else if size != 0 && l.pageCounter[userID] >= 0 {
		l.pageCounter[userID] += int(direction)
	}
}

func (l ListCommand) audioFileInfo(
	ctx context.Context,
	userID string,
	option voice.SearchParams,
) (result []api.AudioFileInfo, err error) {
	for _, s := range l.services {
		if !s.Active(ctx) {
			continue
		}

		result, err = s.AudioFileInfo(ctx, option, uint(l.pageCounter[userID]))
		if err == nil {
			break
		} else {
			slog.With(requestIDKey, ctx.Value(requestIDKey)).WarnContext(ctx, err.Error())
		}
	}

	return result, err
}

func paramsFromDiscord(data discordgo.ApplicationCommandInteractionData) (s voice.SearchParams) {
	s = voice.SearchParams{Type: voice.IDType}

	for _, o := range data.Options {
		switch o.Name {
		case audioNameOption:
			s.Type = voice.NameType
		case audioIdOption:
			s.Type = voice.IDType
		default:
			continue
		}

		s.Value = o.Value.(string)

		return
	}

	return
}

func paramsFromCustomID(customID string) (s voice.SearchParams, err error) {
	data := strings.Split(customID, "_")[1:]
	s.Value = data[1]

	typeNum, err := strconv.Atoi(data[0])
	if err != nil {
		return
	}

	s.Type = voice.AudioQueryType(typeNum)

	return
}

func CustomID(i discordgo.InteractionCreate, pos buttonPosition) (string, error) {
	actionRow, ok := i.Message.Components[0].(*discordgo.ActionsRow)
	if !ok {
		return "", errors.New("failed cast component to *discordgo.ActionsRow")
	}

	button, ok := actionRow.Components[pos].(*discordgo.Button)
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
