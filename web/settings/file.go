package settings

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func moveToNewLocation(oldSrc string, destDir string) (err error) {
	dirs, err := os.ReadDir(oldSrc)
	if err != nil {
		return nil
	}

	var wg sync.WaitGroup

	for _, dir := range dirs {
		if dir.IsDir() || !strings.HasSuffix(dir.Name(), "mp3") {
			continue
		}

		wg.Add(1)
		go func(wg *sync.WaitGroup, oldSrc string, file os.DirEntry) {
			defer wg.Done()

			if err != nil {
				return
			}

			oldPath, newPath := filepath.Join(oldSrc, file.Name()), filepath.Join(destDir, file.Name())
			if err = os.Rename(oldPath, newPath); err != nil {
				return
			}
		}(&wg, oldSrc, dir)
	}

	wg.Wait()

	return
}
