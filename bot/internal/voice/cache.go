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

type SearchParams struct {
	Type  AudioQueryType
	Value string
}

const (
	NameType AudioQueryType = 0x0
	IDType   AudioQueryType = 0x1
)

// TODO Added autocleaning cache after a few days
type BotLocalStorage struct {
	path   string
	active bool
}

func (b BotLocalStorage) Get(ctx context.Context, params SearchParams) (path string, err error) {
	dir, err := os.ReadDir(b.path)
	if err != nil {
		return "", err
	}

	for _, f := range dir {
		select {
		case <-ctx.Done():
			return "", context.Canceled
		default:
		}

		data := strings.Split(f.Name(), "-")
		if len(data) >= 2 && strings.HasPrefix(data[params.Type], params.Value) {
			return filepath.Join(b.path, f.Name()), nil
		}
	}

	return "", fmt.Errorf("song with search value '%s' wasn't found", params.Value)
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

	path := filepath.Join(b.path, fmt.Sprintf("%s-%s.mp3", name, id))
	if _, err := os.Stat(path); err == nil {
		return "", os.ErrExist
	}

	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	for {
		select {
		case <-ctx.Done():
			os.Remove(f.Name())
			return "", context.Canceled
		default:
			const oneMegaByteBuf = 1 << 20
			if n, err := io.CopyN(f, file, oneMegaByteBuf); errors.Is(err, io.EOF) && n != 0 {
				return f.Name(), nil
			} else if err != nil {
				os.Remove(f.Name())
				return "", err
			}
		}
	}
}

func (b BotLocalStorage) Remove(ctx context.Context, params SearchParams) error {
	path, err := b.Get(ctx, params)
	if err != nil {
		return err
	}

	return os.Remove(path)
}

func (b BotLocalStorage) Active(_ context.Context) bool {
	return b.active
}

func (b BotLocalStorage) AudioFileInfo(ctx context.Context, params SearchParams, page uint) ([]api.AudioFileInfo, error) {
	select {
	case <-ctx.Done():
		return nil, context.Canceled
	default:
	}

	dirs, err := os.ReadDir(b.path)
	if err != nil {
		return nil, err
	}

	// number of files, that will be skipped
	skip := page * 10
	if int(skip) > len(dirs) {
		return []api.AudioFileInfo{}, nil
	}

	result := make([]api.AudioFileInfo, 0, 10)
	i := 0

	for _, dir := range dirs[skip:] {
		if i == 10 {
			break
		}

		split := strings.Split(dir.Name(), "-")
		if len(split) != 2 {
			continue
		}

		data := split[params.Type]
		if params.Type == IDType {
			data = strings.TrimSuffix(data, ".mp3")
		}

		if strings.Contains(data, params.Value) {
			result = append(result, api.AudioFileInfo{Filename: split[NameType], ID: split[IDType]})
			i++
		}
	}

	return result, nil
}

func NewBotLocalStorage() BotLocalStorage {
	dir := assertDir()

	if err := os.MkdirAll(dir, 0700); err != nil {
		log.Fatalln("failed create directory for bot cache. " + err.Error())
		return BotLocalStorage{}
	}

	return BotLocalStorage{dir, true}
}
