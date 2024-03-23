package voice

import (
	"context"
	"github.com/bwmarrin/discordgo"
)

type ChatInfo struct {
	UserCount uint
	ChannelID string
}

type ChatHandler struct {
	Ctx             context.Context
	SpockVoiceChns  map[string]chan struct{}
	GuildVoiceChats map[string]ChatInfo
}

func (v *ChatHandler) HandleVoiceChannelUpdate(_ *discordgo.Session, vc *discordgo.VoiceStateUpdate) {
	select {
	case <-v.Ctx.Done():
		return
	default:
	}

	userCount := uint(1)
	info, ok := v.GuildVoiceChats[vc.GuildID]

	if vc.UserID != "" && vc.ChannelID != "" { // User joined to voice channel
		if !ok {
			userCount = info.UserCount + 1
		}

		v.GuildVoiceChats[vc.GuildID] = ChatInfo{
			userCount,
			vc.ChannelID,
		}

		if _, ok = v.SpockVoiceChns[vc.GuildID]; !ok {
			v.SpockVoiceChns[vc.GuildID] = make(chan struct{})
		}
	} else if vc.UserID != "" && vc.ChannelID == "" { // User left the voice channel
		if !ok {
			userCount = info.UserCount - 1
		}

		if userCount == 0 {
			close(v.SpockVoiceChns[vc.GuildID])
			delete(v.GuildVoiceChats, vc.GuildID)
		}
	}
}
