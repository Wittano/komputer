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

func TestNewBotLocalStorage_Get(t *testing.T) {
	dir := t.TempDir()
	if err := os.Setenv(CacheDirAudioKey, dir); err != nil {
		t.Fatal(err)
	}

	s := NewBotLocalStorage()
	const (
		fileID   = "1"
		filename = "file"
	)

	f, err := os.CreateTemp(dir, fmt.Sprintf("%s-%s.mp3", filename, fileID))
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	args := []SearchParams{
		{IDType, fileID},
		{NameType, filename},
	}

	for _, data := range args {
		t.Run("get file from s with query value "+data.Value, func(t *testing.T) {
			path, err := s.Get(context.Background(), data)
			if err != nil {
				t.Fatal(err)
			}

			if _, err := os.Stat(path); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestNewBotLocalStorage_GetWithFileMissing(t *testing.T) {
	dir := t.TempDir()
	if err := os.Setenv(CacheDirAudioKey, dir); err != nil {
		t.Fatal(err)
	}

	s := NewBotLocalStorage()
	if _, err := s.Get(context.Background(), (SearchParams{Type: IDType, Value: "1123412"})); err == nil {
		t.Fatal(err)
	}
}

func TestNewBotLocalStorage_RemoveWithFileMissing(t *testing.T) {
	dir := t.TempDir()
	if err := os.Setenv(CacheDirAudioKey, dir); err != nil {
		t.Fatal(err)
	}

	s := NewBotLocalStorage()
	params := SearchParams{Type: IDType, Value: "142308579as"}

	if err := s.Remove(context.Background(), params); err == nil {
		t.Fatal("random file was removed, but s service shouldn't remove anything")
	}
}

func TestNewBotLocalStorage_Remove(t *testing.T) {
	dir := t.TempDir()
	if err := os.Setenv(CacheDirAudioKey, dir); err != nil {
		t.Fatal(err)
	}

	s := NewBotLocalStorage()
	const (
		fileID   = "1"
		filename = "file"
	)

	params := []SearchParams{
		{IDType, fileID},
		{NameType, filename},
	}

	for _, data := range params {
		t.Run("remove file from s with query value "+data.Value, func(t *testing.T) {
			f, err := os.CreateTemp(dir, fmt.Sprintf("%s-%s.mp3", filename, fileID))
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			if err := s.Remove(context.Background(), data); err != nil {
				t.Fatal(err)
			}

			if _, err := os.Stat(f.Name()); err == nil {
				t.Fatalf("file %s wasn't removed", f.Name())
			}
		})
	}
}

func TestNewBotLocalStorage_Add(t *testing.T) {
	if err := os.Setenv(CacheDirAudioKey, filepath.Join(t.TempDir(), "assets")); err != nil {
		t.Fatal(err)
	}

	s := NewBotLocalStorage()
	const (
		fileID      = "1"
		filename    = "file"
		fileContent = "test data"
	)
	buf := strings.NewReader(fileContent)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	path, err := s.Add(ctx, buf, fileID, filename)
	if err != nil {
		t.Fatal(err)
	}

	if stat, err := os.Stat(path); err != nil || stat.Size() != int64(len(fileContent)) {
		t.Fatalf("file %s wasn't saved correctly. %s", path, err)
	}
}

func TestNewBotLocalStorage_AddWithEmptySource(t *testing.T) {
	if err := os.Setenv(CacheDirAudioKey, filepath.Join(t.TempDir(), "assets")); err != nil {
		t.Fatal(err)
	}

	s := NewBotLocalStorage()
	const (
		fileID      = "1"
		fileContent = "test data"
	)
	buf := strings.NewReader(fileContent)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	path, err := s.Add(ctx, buf, fileID, "")
	if err != nil {
		t.Fatal(err)
	}

	if stat, err := os.Stat(path); err != nil || stat.Size() != int64(len(fileContent)) {
		t.Fatalf("file %s wasn't saved correctly. %s", path, err)
	}
}

func TestNewBotLocalStorage_AddWithMissingFileID(t *testing.T) {
	if err := os.Setenv(CacheDirAudioKey, filepath.Join(t.TempDir(), "assets")); err != nil {
		t.Fatal(err)
	}

	s := NewBotLocalStorage()
	const (
		filename    = "file"
		fileContent = "test data"
	)
	buf := strings.NewReader(fileContent)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	path, err := s.Add(ctx, buf, "", filename)
	if err == nil {
		t.Fatalf("file %s was created, but ID is missing", path)
	}
}

func TestNewBotLocalStorage_AddWithContextCanceled(t *testing.T) {
	if err := os.Setenv(CacheDirAudioKey, filepath.Join(t.TempDir(), "assets")); err != nil {
		t.Fatal(err)
	}

	s := NewBotLocalStorage()
	const (
		fileID   = "1"
		filename = "file"
	)

	const fileContent = "test data"
	buf := strings.NewReader(fileContent)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	path, err := s.Add(ctx, buf, fileID, filename)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("file %s was created, but context was canceled. %s", path, err)
	}
}

func TestNewBotLocalStorageAddButFileExists(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "assets")
	if err := os.Setenv(CacheDirAudioKey, dir); err != nil {
		t.Fatal(err)
	}

	s := NewBotLocalStorage()
	const (
		filename    = "file"
		fileID      = "1"
		fileContent = "test data"
	)
	buf := strings.NewReader(fileContent)
	f, err := os.Create(filepath.Join(dir, fmt.Sprintf("%s-%s.mp3", filename, fileID)))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	path, err := s.Add(ctx, buf, fileID, filename)
	if !errors.Is(err, os.ErrExist) {
		t.Fatalf("file %s was created, but file was existed previous. %s", path, err)
	}
}
