package command

import (
	"context"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/wittano/komputer/audio"
	"os"
	"regexp"
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
	content []string
	pattern string
}

func (l listContent) Response() *discordgo.InteractionResponseData {
	if len(l.content) == 0 {
		return SimpleMessage{Msg: "Nie znalazłem żadnego takiego pliku"}.Response()
	}

	var msg strings.Builder

	for _, info := range l.content {
		msg.WriteString(fmt.Sprintf("- %s\n", info))
	}

	const customIDFormat = "%s_%s"
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
						CustomID: fmt.Sprintf(customIDFormat, previousIdsButtonID, l.pattern),
					},
					discordgo.Button{
						Style:    discordgo.PrimaryButton,
						Label:    "Next",
						CustomID: fmt.Sprintf(customIDFormat, nextIdsButtonID, l.pattern),
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
	return regexp.MustCompile(fmt.Sprintf("^%s_(a-z0-9)*", nextIdsButtonID)).MatchString(customID)
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
	return regexp.MustCompile(fmt.Sprintf("^%s_(a-z0-9)*", previousIdsButtonID)).MatchString(customID)
}

func (p PreviousListCommandOption) Execute(ctx context.Context, _ *discordgo.Session, i *discordgo.InteractionCreate) (DiscordMessageReceiver, error) {
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
	// map of pages number, which user should see
	pageCounter map[string]int
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
	name string,
) (result []string, err error) {
	dirs, err := os.ReadDir(audio.AssertDir())
	if err != nil {
		return nil, err
	}

	page, ok := l.pageCounter[userID]
	if !ok {
		page = 0
	}

	const maxPageSize = 10
	skip := page * maxPageSize

	if len(dirs) <= skip {
		return nil, nil
	}

	for _, s := range dirs[skip:] {
		select {
		case <-ctx.Done():
			return nil, context.Canceled
		default:
			if name == "" && strings.HasPrefix(s.Name(), name) {
				result = append(result, s.Name())
			}
		}
	}

	return result, err
}

func paramsFromDiscord(data discordgo.ApplicationCommandInteractionData) (name string) {
	for _, o := range data.Options {
		switch o.Name {
		case audioNameOption:
			name = o.Value.(string)
		default:
			continue
		}

		return
	}

	return
}

func paramsFromCustomID(customID string) (string, error) {
	data := strings.Split(customID, "_")
	if len(data) != 2 {
		return "", errors.New("invalid customID")
	}

	return data[1], nil
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

func NewListCommand() *ListCommand {
	return &ListCommand{make(map[string]int)}
}
