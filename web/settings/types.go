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
	Count int64 `yaml:"max_file_count" json:"max_file_count"`
	Size  int64 `yaml:"max_file_size" json:"max_file_size"`
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

	if new.Upload.Count != 0 && s.Upload.Count != new.Upload.Count {
		s.Upload.Count = new.Upload.Count
	}

	if new.Upload.Size != 0 && s.Upload.Size != new.Upload.Size {
		s.Upload.Size = new.Upload.Size
	}

	return nil
}

func (s Settings) CheckFilesLimit(c int) bool {
	return c >= 1 && int64(c) <= s.Upload.Count
}

var Config *Settings

// TODO Remove global variable/singleton
func Load(path string) error {
	if Config != nil {
		return nil
	}

	destPath := DefaultSettingsPath
	if path != "" {
		destPath = path
	}

	if _, err := os.Stat(destPath); errors.Is(err, os.ErrNotExist) {
		Config, err = defaultSettings(destPath)
		return err
	}

	f, err := os.Open(destPath)
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
	dir := DefaultAssertDir
	if assetDirPath, ok := os.LookupEnv(cacheDirKey); ok && assetDirPath != "" {
		dir = assetDirPath
	}

	def := Settings{
		AssetDir: dir,
		Upload: UploadSettings{
			Count: 5,
			Size:  defaultMaxFileSize,
		},
	}

	e := yaml.NewEncoder(f)
	defer e.Close()
	if err = e.Encode(&def); err != nil {
		return nil, err
	}

	return &def, nil
}
