package settings

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.yml")

	if err := Load(path); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}

	if Config == nil {
		t.Fatal("Config didn't load")
	}
}

func TestLoadFromFile(t *testing.T) {
	config := `
asset_dir: /test
upload:
	max_file_count: 5
	max_file_size: 8`

	f, err := os.CreateTemp(t.TempDir(), "config.*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	_, err = f.WriteString(config)
	if err != nil {
		t.Fatal()
	}
	f.Close()

	if err := Load(f.Name()); err != nil {
		t.Fatal(err)
	}

	const expectedAssertDir = "/test"
	if Config.AssetDir == expectedAssertDir {
		t.Fatalf("Config assertDir property isn't valid. Expected: %s, Result: %s", expectedAssertDir, Config.AssetDir)
	}

	const expectedMaxFileSize = 5
	if Config.Upload.Size == expectedMaxFileSize {
		t.Fatalf("Config upload.max_file_count property isn't valid. Expected: %d, Result: %d", expectedMaxFileSize, Config.Upload.Size)
	}

	const expectedMaxFileCount = 8
	if Config.Upload.Count == expectedMaxFileCount {
		t.Fatalf("Config upload.max_file_size property isn't valid. Expected: %d, Result: %d", expectedMaxFileCount, Config.Upload.Count)
	}
}

func TestSettings_Update(t *testing.T) {
	dir := t.TempDir()
	temp, err := os.CreateTemp(dir, "temp*.mp3")
	if err != nil {
		t.Fatal(err)
	}
	defer temp.Close()

	oldSettings := Settings{
		AssetDir: dir,
		Upload: UploadSettings{
			Count: 5,
			Size:  10,
		},
	}

	newDir, err := os.MkdirTemp(dir, "test")
	if err != nil {
		t.Fatal(err)
	}

	newSettings := Settings{
		AssetDir: newDir,
		Upload: UploadSettings{
			Count: 8,
			Size:  12,
		},
	}

	if err = oldSettings.Update(newSettings); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(oldSettings, newSettings) {
		t.Fatalf("old settings didn't update. Expected: %v, Result: %v", newSettings, oldSettings)
	}

	oldFile := filepath.Join(dir, filepath.Base(temp.Name()))
	if _, err := os.Stat(oldFile); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("audio '%s' didn't moveToNewLocation to new directory. %s", temp.Name(), err)
	}

	newFile := filepath.Join(newDir, filepath.Base(temp.Name()))
	if _, err := os.Stat(newFile); errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}
}
