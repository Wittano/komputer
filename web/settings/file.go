package settings

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func moveToNewLocation(old string, new string) (err error) {
	dirs, err := os.ReadDir(old)
	if err != nil {
		return nil
	}

	var wg sync.WaitGroup

	for _, dir := range dirs {
		if dir.IsDir() || !strings.HasSuffix(dir.Name(), "mp3") {
			continue
		}

		wg.Add(1)
		go func(wg *sync.WaitGroup, old string, dir os.DirEntry) {
			defer wg.Done()

			if err != nil {
				return
			}

			src, dest := filepath.Join(old, dir.Name()), filepath.Join(new, dir.Name())
			if err = os.Rename(src, dest); err != nil {
				return
			}
		}(&wg, old, dir)
	}

	wg.Wait()

	return
}
