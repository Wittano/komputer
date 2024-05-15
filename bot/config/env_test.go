package config

import (
	"os"
	"testing"
)

const (
	discordToken = "123"
	appID        = "456"
)

func TestLoadBotVariables(t *testing.T) {
	os.Setenv("DISCORD_BOT_TOKEN", discordToken)
	os.Setenv("APPLICATION_ID", appID)

	prop, err := NewBotProperties()
	if err != nil {
		t.Fatal(err)
	}

	if prop.Token != discordToken {
		t.Fatalf("Invalid token! Expected: '%s', Result: '%s'", discordToken, prop.Token)
	}

	if prop.AppID != appID {
		t.Fatalf("Invalid token! Expected: '%s', Result: '%s'", appID, prop.AppID)
	}

	if prop.ServerGUID != "" {
		t.Fatalf("Invalid ServerGUID! Expected: '%s', Result: '%s'", "", prop.ServerGUID)
	}
}

func TestLoadBotVariables_DiscordTokenMissing(t *testing.T) {
	os.Setenv("DISCORD_BOT_TOKEN", "")
	os.Setenv("APPLICATION_ID", appID)

	if _, err := LoadBotVariables(); err == nil {
		t.Fatal("Properties was loaded!")
	}
}

func TestLoadBotVariables_AppIDMissing(t *testing.T) {
	os.Setenv("DISCORD_BOT_TOKEN", discordToken)
	os.Setenv("APPLICATION_ID", "")

	if _, err := LoadBotVariables(); err == nil {
		t.Fatal("Properties was loaded!")
	}
}
