package voice

import "github.com/bwmarrin/discordgo"

var UsersOnVoiceChat = map[string]string{}

func HandleVoiceChannelUpdate(_ *discordgo.Session, vc *discordgo.VoiceStateUpdate) {
	if vc.UserID != "" && vc.ChannelID != "" {
		UsersOnVoiceChat[vc.UserID] = vc.ChannelID
	}
}
