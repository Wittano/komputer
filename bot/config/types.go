package config

import (
	"errors"
	"fmt"
	"os"
)

type BotProperties struct {
	Token      string
	AppID      string
	ServerGUID string
}

func NewBotProperties() (prop BotProperties, err error) {
	prop.Token, err = loadEnv("DISCORD_BOT_TOKEN")
	if err != nil {
		return
	}

	prop.AppID, err = loadEnv("APPLICATION_ID")
	if err != nil {
		return
	}

	prop.ServerGUID, _ = loadEnv("SERVER_GUID")

	return
}

func loadEnv(name string) (env string, err error) {
	if value, ok := os.LookupEnv(name + "_PATH"); ok {
		return loadFromFile(value)
	} else {
		return loadFromEnvVar(name)
	}
}

func loadFromEnvVar(name string) (env string, err error) {
	env, ok := os.LookupEnv(name)
	if !ok || env == "" {
		return "", fmt.Errorf("missing %s variable", name)
	}

	return
}

func loadFromFile(path string) (env string, err error) {
	if _, err = os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return
	}

	b, err := os.ReadFile(path)
	return string(b), err
}
