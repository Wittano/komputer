package audio

import (
	"errors"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

const (
	defaultCacheDirAudio = "assets"
	assetsDirKey         = "ASSETS_DIR"
)

func Path(name string) string {
	return filepath.Join(AssertDir(), name)
}

func AssertDir() (path string) {
	path = defaultCacheDirAudio
	if cacheDir, ok := os.LookupEnv(assetsDirKey); ok && cacheDir != "" {
		path = cacheDir
	}

	return
}

func Paths() (paths []string, err error) {
	dirs, err := os.ReadDir(AssertDir())
	if err != nil {
		return nil, err
	}

	if len(dirs) <= 0 {
		return nil, errors.New("assert directory is empty")
	}

	paths = make([]string, 0, len(dirs))
	for _, dir := range dirs {
		if dir.Type() != os.ModeDir {
			paths = append(paths, dir.Name())
		}
	}

	return
}

// PathsWithPagination get fixed-sized list of audio paths from assert dictionary
func PathsWithPagination(page uint32, size uint32) (paths []string, err error) {
	dirs, err := os.ReadDir(AssertDir())
	if err != nil {
		return nil, err
	}

	skipFiles := int(page * size)
	if len(dirs) < skipFiles {
		return []string{}, nil
	} else if len(dirs) <= 0 {
		return nil, errors.New("assert directory is empty")
	}

	paths = make([]string, 0, size)
	for _, dir := range dirs[skipFiles:] {
		if dir.Type() != os.ModeDir {
			paths = append(paths, dir.Name())
		}

		if len(paths) >= int(size) {
			break
		}
	}

	return
}

func RandomPath() (string, error) {
	paths, err := Paths()
	if err != nil {
		return "", err
	}

	return paths[rand.Int()%len(paths)], nil
}

func Duration(path string) (duration time.Duration, err error) {
	cmd := exec.Command("ffprobe", "-i", path, "-show_entries", "format=duration", "-v", "quiet", "-of", "csv='p=0'")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	if err = cmd.Start(); err != nil {
		return
	}

	out, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	defer out.Close()

	rawTime, err := io.ReadAll(out)
	if err != nil {
		return
	}

	return time.ParseDuration(string(rawTime) + "s")
}
