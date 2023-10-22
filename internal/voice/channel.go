package voice

import (
	"github.com/bwmarrin/discordgo"
)

type UserVoiceChat struct {
	ChannelID string
	GuildID   string
}

var UserVoiceChatMap = map[string]UserVoiceChat{}

func HandleVoiceChannelUpdate(_ *discordgo.Session, vc *discordgo.VoiceStateUpdate) {
	uvc := UserVoiceChat{
		ChannelID: vc.ChannelID,
		GuildID:   vc.GuildID,
	}

	if vc.UserID != "" && vc.ChannelID != "" {
		UserVoiceChatMap[vc.UserID] = uvc
	} else if vc.UserID != "" && vc.ChannelID == "" {
		delete(UserVoiceChatMap, vc.UserID)
	}
}
