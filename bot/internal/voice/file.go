package voice

import (
	"errors"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultCacheDirAudio = "assets"
	CacheDirAudioKey     = "CACHE_AUDIO_DIR"
)

// FIXME added name of files or searching file in local storage by ID
func Path(id string) string {
	return filepath.Join(assertDir(), id+".mp3")
}

func assertDir() (path string) {
	path = defaultCacheDirAudio
	if cacheDir, ok := os.LookupEnv(CacheDirAudioKey); ok && cacheDir != "" {
		path = cacheDir
	}

	return
}

func AudioIDs() ([]string, error) {
	dirs, err := os.ReadDir(assertDir())
	if err != nil {
		return nil, err
	}

	if len(dirs) <= 0 {
		return nil, errors.New("assert directory is empty")
	}

	ids := make([]string, 0, len(dirs))
	for _, dir := range dirs {
		const suffix = ".mp3"

		if dir.Type() == fs.ModeDir || !strings.HasSuffix(dir.Name(), suffix) {
			continue
		}

		filename := strings.Split(dir.Name(), "-")
		ids = append(ids, strings.TrimSuffix(filename[1], suffix))
	}

	return ids, nil
}

func RandomAudioID() (string, error) {
	ids, err := AudioIDs()
	if err != nil {
		return "", err
	}

	return ids[rand.Int()%len(ids)], nil
}
