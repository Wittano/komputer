package test

import (
	"github.com/wittano/komputer/web/settings"
	"path/filepath"
	"testing"
)

func LoadDefaultConfig(t *testing.T) error {
	const defaultConfigFileName = "config.yml"
	configFile := filepath.Join(t.TempDir(), defaultConfigFileName)

	if err := settings.Load(configFile); err != nil {
		return err
	}

	return settings.Config.Update(settings.Settings{AssetDir: filepath.Join(t.TempDir(), "assets")})
}
