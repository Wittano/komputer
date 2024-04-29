//go:build dev

package config

import "github.com/joho/godotenv"

func LoadBotVariables() (BotProperties, error) {
	godotenv.Load()

	return NewBotProperties()
}
