package voice

import (
	"os"
	"path/filepath"
)

const (
	defaultCacheDirAudio = "assets"
	CacheDirAudioKey     = "CACHE_AUDIO_DIR"
)

func CheckIfAudioIsDownloaded(id string) (err error) {
	_, err = os.Stat(Path(id))
	return
}

func Path(id string) string {
	cachePath := defaultCacheDirAudio
	if cacheDir, ok := os.LookupEnv(CacheDirAudioKey); ok && cacheDir != "" {
		cachePath = cacheDir
	}

	return filepath.Join(cachePath, id+".mp3")
}
