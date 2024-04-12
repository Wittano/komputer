package settings

import (
	"errors"
	"gopkg.in/yaml.v3"
	"os"
)

const (
	DefaultAssertDir    = "assets"
	DefaultSettingsPath = "settings.yml"
)

const defaultMaxFileSize = 8 * (1 << 20) // 8MB in bytes

type UploadSettings struct {
	MaxFileCount int64 `yaml:"max_file_count" json:"max_file_count"`
	MaxFileSize  int64 `yaml:"max_file_size" json:"max_file_size"`
}

type Settings struct {
	AssetDir string         `yaml:"asset_dir" json:"asset_dir"`
	Upload   UploadSettings `yaml:"upload" json:"upload"`
}

func (s *Settings) Update(new Settings) error {
	if new.AssetDir != "" && s.AssetDir != new.AssetDir {
		if err := os.MkdirAll(new.AssetDir, 0700); err != nil {
			return err
		}

		err := moveToNewLocation(s.AssetDir, new.AssetDir)
		if err != nil {
			return err
		}

		s.AssetDir = new.AssetDir
	}

	if new.Upload.MaxFileCount != 0 && s.Upload.MaxFileCount != new.Upload.MaxFileCount {
		s.Upload.MaxFileCount = new.Upload.MaxFileCount
	}

	if new.Upload.MaxFileSize != 0 && s.Upload.MaxFileSize != new.Upload.MaxFileSize {
		s.Upload.MaxFileSize = new.Upload.MaxFileSize
	}

	return nil
}

func (s Settings) CheckFileCountLimit(count int) bool {
	return count >= 1 && int64(count) <= s.Upload.MaxFileCount
}

var Config *Settings

func Load(path string) error {
	if Config != nil {
		return nil
	}

	settingPath := DefaultSettingsPath
	if path != "" {
		settingPath = path
	}

	if _, err := os.Stat(settingPath); errors.Is(err, os.ErrNotExist) {
		Config, err = defaultSettings(settingPath)
		return err
	}

	f, err := os.Open(settingPath)
	if err != nil {
		return err
	}
	defer f.Close()

	d := yaml.NewDecoder(f)
	if err = d.Decode(&Config); err != nil {
		return err
	}

	return os.MkdirAll(Config.AssetDir, 0700)
}

func defaultSettings(path string) (*Settings, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	const cacheDirKey = "CACHE_AUDIO_DIR"
	assetDir := DefaultAssertDir
	if assetDirPath, ok := os.LookupEnv(cacheDirKey); ok && assetDirPath != "" {
		assetDir = assetDirPath
	}

	defaultSettings := Settings{
		AssetDir: assetDir,
		Upload: UploadSettings{
			MaxFileCount: 5,
			MaxFileSize:  defaultMaxFileSize,
		},
	}

	e := yaml.NewEncoder(f)
	defer e.Close()
	if err = e.Encode(&defaultSettings); err != nil {
		return nil, err
	}

	return &defaultSettings, nil
}
