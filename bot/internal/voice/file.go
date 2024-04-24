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

func Path(id string) string {
	return filepath.Join(assertDir(), id+".mp3")
}

func assertDir() (cachePath string) {
	cachePath = defaultCacheDirAudio
	if cacheDir, ok := os.LookupEnv(CacheDirAudioKey); ok && cacheDir != "" {
		cachePath = cacheDir
	}

	return
}

func AudioIDs() ([]string, error) {
	assertsPath := assertDir()
	files, err := os.ReadDir(assertsPath)
	if err != nil {
		return nil, err
	}

	if len(files) <= 0 {
		return nil, errors.New("assert directory is empty")
	}

	ids := make([]string, 0, len(files))
	for _, f := range files {
		const suffix = ".mp3"

		if f.Type() == fs.ModeDir || !strings.HasSuffix(f.Name(), suffix) {
			continue
		}

		filename := strings.Split(f.Name(), "-")
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
