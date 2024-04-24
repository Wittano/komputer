package voice

import (
	"context"
	"errors"
	"fmt"
	"github.com/wittano/komputer/api"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type AudioQueryType byte

type AudioSearch struct {
	Type  AudioQueryType
	Value string
}

const (
	NameType AudioQueryType = 0x0
	IDType   AudioQueryType = 0x1
)

type BotLocalStorage struct {
	storagePath string
}

func (b BotLocalStorage) Get(ctx context.Context, query AudioSearch) (path string, err error) {
	dir, err := os.ReadDir(b.storagePath)
	if err != nil {
		return "", err
	}

	for _, f := range dir {
		select {
		case <-ctx.Done():
			return "", context.Canceled
		default:
		}

		name := strings.Split(f.Name(), "-")
		if len(name) >= 2 && strings.HasPrefix(name[query.Type], query.Value) {
			return filepath.Join(b.storagePath, f.Name()), nil
		}
	}

	return "", fmt.Errorf("song with search value '%s' wasn't found", query.Value)
}

func (b BotLocalStorage) Add(ctx context.Context, file io.Reader, id, name string) (string, error) {
	select {
	case <-ctx.Done():
		return "", context.Canceled
	default:
	}

	if name == "" {
		name = "file"
	}

	if id == "" {
		return "", errors.New("audio ID is empty")
	}

	destPath := filepath.Join(b.storagePath, fmt.Sprintf("%s-%s.mp3", name, id))
	if _, err := os.Stat(destPath); err == nil {
		return "", os.ErrExist
	}

	destFile, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer destFile.Close()

	for {
		select {
		case <-ctx.Done():
			os.Remove(destFile.Name())
			return "", context.Canceled
		default:
			const oneMegaByteBuf = 1 << 20
			if n, err := io.CopyN(destFile, file, oneMegaByteBuf); errors.Is(err, io.EOF) && n != 0 {
				return destFile.Name(), nil
			} else if err != nil {
				os.Remove(destFile.Name())
				return "", err
			}
		}
	}
}

func (b BotLocalStorage) Remove(ctx context.Context, query AudioSearch) error {
	path, err := b.Get(ctx, query)
	if err != nil {
		return err
	}

	return os.Remove(path)
}

func (b BotLocalStorage) SearchAudio(ctx context.Context, option AudioSearch, page uint) ([]api.AudioFileInfo, error) {
	select {
	case <-ctx.Done():
		return nil, context.Canceled
	default:
	}

	dirs, err := os.ReadDir(b.storagePath)
	if err != nil {
		return nil, err
	}

	index := page * 10
	if int(index) > len(dirs) {
		return []api.AudioFileInfo{}, nil
	}

	result := make([]api.AudioFileInfo, 0, 10)
	i := 0

	for _, dir := range dirs[index:] {
		if i == 10 {
			break
		}

		split := strings.Split(dir.Name(), "-")
		if len(split) != 2 {
			continue
		}

		value := split[option.Type]
		if option.Type == IDType {
			value = strings.TrimSuffix(value, ".mp3")
		}

		if strings.Contains(value, option.Value) {
			result = append(result, api.AudioFileInfo{Filename: split[NameType], ID: split[IDType]})
			i++
		}
	}

	return result, nil
}

func NewBotLocalStorage() BotLocalStorage {
	dirPath := assertDir()

	if err := os.MkdirAll(dirPath, 0700); err != nil {
		log.Fatalln("failed create directory for bot cache. " + err.Error())
		return BotLocalStorage{}
	}

	return BotLocalStorage{dirPath}
}
