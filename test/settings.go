package test

import (
	"github.com/wittano/komputer/web/settings"
	"path/filepath"
	"testing"
)

func LoadDefaultConfig(t *testing.T) error {
	const defaultConfigFileName = "config.yml"
	path := filepath.Join(t.TempDir(), defaultConfigFileName)

	if err := settings.Load(path); err != nil {
		return err
	}

	return settings.Config.Update(settings.Settings{AssetDir: filepath.Join(t.TempDir(), "assets")})
}
