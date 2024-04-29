package config

import (
	"fmt"
	"os"
)

type BotProperties struct {
	Token      string
	AppID      string
	ServerGUID string
}

func NewBotProperties() (prop BotProperties, err error) {
	var ok bool

	prop.Token, ok = os.LookupEnv("DISCORD_BOT_TOKEN")
	if !ok || prop.Token == "" {
		err = fmt.Errorf("missing DISCORD_BOT_TOKEN variable")
		return
	}

	prop.AppID, ok = os.LookupEnv("APPLICATION_ID")
	if !ok || prop.AppID == "" {
		err = fmt.Errorf("missing DISCORD_BOT_TOKEN variable")
		return
	}

	prop.ServerGUID = os.Getenv("SERVER_GUID")

	return
}
