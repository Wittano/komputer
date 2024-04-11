//go:build !dev

package config

import (
	"errors"
	"os"
)

func LoadBotVariables() (BotProperties, error) {
	if _, ok := os.LookupEnv("SERVER_GUID"); ok {
		return BotProperties{}, errors.New("variable 'SERVER_GUID' allows only in development mode")
	}

	return NewBotProperties()
}
