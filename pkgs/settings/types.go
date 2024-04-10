package settings

import (
	"errors"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"sync"
)

const (
	DefaultAssertDir    = ".cache/komputer"
	DefaultSettingsPath = ".config/komputer/settings.yml"
)

const maxFileSize = 8 * (1 << 20) // 8MB in bytes

type UploadSettings struct {
	MaxFileCount uint `yaml:"max_file_count" json:"max_file_count"`
	MaxFileSize  uint `yaml:"max_file_size" json:"max_file_size"`
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

		err := moveAssets(s.AssetDir, new.AssetDir)
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

var Config *Settings

func Load(path string) error {
	if Config != nil {
		return nil
	}

	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	settingPath := filepath.Join(home, DefaultSettingsPath)
	if path != "" {
		settingPath = filepath.Join(path)
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

	return nil
}

func defaultSettings(path string) (*Settings, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	defaultSettings := Settings{
		AssetDir: DefaultAssertDir,
		Upload: UploadSettings{
			MaxFileCount: 5,
			MaxFileSize:  maxFileSize,
		},
	}

	e := yaml.NewEncoder(f)
	defer e.Close()
	if err = e.Encode(&defaultSettings); err != nil {
		return nil, err
	}

	return &defaultSettings, nil
}

func moveAssets(oldSrc string, path string) (err error) {
	dirs, err := os.ReadDir(oldSrc)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(dirs))

	for _, dir := range dirs {
		go func(wg *sync.WaitGroup, oldSrc string, file os.DirEntry) {
			defer wg.Done()

			if err != nil {
				return
			}

			filename := filepath.Join(oldSrc, file.Name())
			newPath := filepath.Join(path, filepath.Base(filename))
			if err = os.Rename(filename, newPath); err != nil {
				return
			}
		}(&wg, oldSrc, dir)
	}

	wg.Wait()

	return
}
