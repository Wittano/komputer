package voice

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewBotLocalStorageGet(t *testing.T) {
	dir := t.TempDir()
	if err := os.Setenv(CacheDirAudioKey, dir); err != nil {
		t.Fatal(err)
	}

	cache := NewBotLocalStorage()
	const (
		fileID   = "1"
		filename = "file"
	)

	f, err := os.CreateTemp(dir, fmt.Sprintf("%s-%s.mp3", filename, fileID))
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	args := []AudioSearch{
		{IDType, fileID},
		{NameType, filename},
	}

	for _, data := range args {
		t.Run("get file from cache with query value "+data.Value, func(t *testing.T) {
			path, err := cache.Get(context.Background(), data)
			if err != nil {
				t.Fatal(err)
			}

			if _, err := os.Stat(path); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestNewBotLocalStorageGetButFileMissing(t *testing.T) {
	dir := t.TempDir()
	if err := os.Setenv(CacheDirAudioKey, dir); err != nil {
		t.Fatal(err)
	}

	cache := NewBotLocalStorage()
	query := AudioSearch{Type: IDType, Value: "1123412"}

	if _, err := cache.Get(context.Background(), query); err == nil {
		t.Fatal(err)
	}
}

func TestNewBotLocalStorageRemoveButFileMissing(t *testing.T) {
	dir := t.TempDir()
	if err := os.Setenv(CacheDirAudioKey, dir); err != nil {
		t.Fatal(err)
	}

	cache := NewBotLocalStorage()
	query := AudioSearch{Type: IDType, Value: "142308579as"}

	if err := cache.Remove(context.Background(), query); err == nil {
		t.Fatal("random file was removed, but cache service shouldn't remove anything")
	}
}

func TestNewBotLocalStorageRemove(t *testing.T) {
	dir := t.TempDir()
	if err := os.Setenv(CacheDirAudioKey, dir); err != nil {
		t.Fatal(err)
	}

	cache := NewBotLocalStorage()
	const (
		fileID   = "1"
		filename = "file"
	)

	args := []AudioSearch{
		{IDType, fileID},
		{NameType, filename},
	}

	for _, data := range args {
		t.Run("remove file from cache with query value "+data.Value, func(t *testing.T) {
			f, err := os.CreateTemp(dir, fmt.Sprintf("%s-%s.mp3", filename, fileID))
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			if err := cache.Remove(context.Background(), data); err != nil {
				t.Fatal(err)
			}

			if _, err := os.Stat(f.Name()); err == nil {
				t.Fatalf("file %s wasn't removed", f.Name())
			}
		})
	}
}

func TestNewBotLocalStorageAdd(t *testing.T) {
	if err := os.Setenv(CacheDirAudioKey, filepath.Join(t.TempDir(), "assets")); err != nil {
		t.Fatal(err)
	}

	cache := NewBotLocalStorage()
	const (
		fileID      = "1"
		filename    = "file"
		fileContent = "test data"
	)
	buf := strings.NewReader(fileContent)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	path, err := cache.Add(ctx, buf, fileID, filename)
	if err != nil {
		t.Fatal(err)
	}

	if stat, err := os.Stat(path); err != nil || stat.Size() != int64(len(fileContent)) {
		t.Fatalf("file %s wasn't saved correctly. %s", path, err)
	}
}

func TestNewBotLocalStorageAddButFilenameIsMissing(t *testing.T) {
	if err := os.Setenv(CacheDirAudioKey, filepath.Join(t.TempDir(), "assets")); err != nil {
		t.Fatal(err)
	}

	cache := NewBotLocalStorage()
	const (
		fileID      = "1"
		fileContent = "test data"
	)
	buf := strings.NewReader(fileContent)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	path, err := cache.Add(ctx, buf, fileID, "")
	if err != nil {
		t.Fatal(err)
	}

	if stat, err := os.Stat(path); err != nil || stat.Size() != int64(len(fileContent)) {
		t.Fatalf("file %s wasn't saved correctly. %s", path, err)
	}
}

func TestNewBotLocalStorageAddButMissingFileID(t *testing.T) {
	if err := os.Setenv(CacheDirAudioKey, filepath.Join(t.TempDir(), "assets")); err != nil {
		t.Fatal(err)
	}

	cache := NewBotLocalStorage()
	const (
		filename    = "file"
		fileContent = "test data"
	)
	buf := strings.NewReader(fileContent)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	path, err := cache.Add(ctx, buf, "", filename)
	if err == nil {
		t.Fatalf("file %s was created, but ID is missing", path)
	}
}

func TestNewBotLocalStorageAddButContextCanceled(t *testing.T) {
	if err := os.Setenv(CacheDirAudioKey, filepath.Join(t.TempDir(), "assets")); err != nil {
		t.Fatal(err)
	}

	cache := NewBotLocalStorage()
	const (
		fileID   = "1"
		filename = "file"
	)

	const fileContent = "test data"
	buf := strings.NewReader(fileContent)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	path, err := cache.Add(ctx, buf, fileID, filename)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("file %s was created, but context was canceled. %s", path, err)
	}
}

func TestNewBotLocalStorageAddButFileExists(t *testing.T) {
	destDir := filepath.Join(t.TempDir(), "assets")
	if err := os.Setenv(CacheDirAudioKey, destDir); err != nil {
		t.Fatal(err)
	}

	cache := NewBotLocalStorage()
	const (
		filename    = "file"
		fileID      = "1"
		fileContent = "test data"
	)
	buf := strings.NewReader(fileContent)
	f, err := os.Create(filepath.Join(destDir, fmt.Sprintf("%s-%s.mp3", filename, fileID)))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	path, err := cache.Add(ctx, buf, fileID, filename)
	if !errors.Is(err, os.ErrExist) {
		t.Fatalf("file %s was created, but file was existed previous. %s", path, err)
	}
}
