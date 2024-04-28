package voice

import (
	"context"
	"github.com/bwmarrin/discordgo"
)

type ChatInfo struct {
	UserCount uint
	ChannelID string
}

type SpockVoiceChannels map[string]chan struct{}

func (s SpockVoiceChannels) Close() error {
	for _, ch := range s {
		close(ch)
	}

	return nil
}

type ChatHandler struct {
	Ctx             context.Context
	SpockVoiceChns  SpockVoiceChannels
	GuildVoiceChats map[string]ChatInfo
}

func (v *ChatHandler) HandleVoiceChannelUpdate(_ *discordgo.Session, vc *discordgo.VoiceStateUpdate) {
	select {
	case <-v.Ctx.Done():
		return
	default:
	}

	count := uint(1)
	info, ok := v.GuildVoiceChats[vc.GuildID]

	if vc.UserID != "" && vc.ChannelID != "" { // User joined to voice channel
		if !ok {
			count = info.UserCount + 1
		}

		v.GuildVoiceChats[vc.GuildID] = ChatInfo{
			count,
			vc.ChannelID,
		}

		if _, ok = v.SpockVoiceChns[vc.GuildID]; !ok {
			v.SpockVoiceChns[vc.GuildID] = make(chan struct{})
		}
	} else if vc.UserID != "" && vc.ChannelID == "" { // User left the voice channel
		if !ok {
			count = info.UserCount - 1
		}

		if count == 0 {
			close(v.SpockVoiceChns[vc.GuildID])
			delete(v.GuildVoiceChats, vc.GuildID)
		}
	}
}
