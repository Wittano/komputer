package settings

import (
	"os"
	"path/filepath"
	"sync"
)

func moveToNewLocation(oldSrc string, path string) (err error) {
	dirs, err := os.ReadDir(oldSrc)
	if err != nil {
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(len(dirs))

	for _, dir := range dirs {
		go func(wg *sync.WaitGroup, oldSrc string, file os.DirEntry) {
			defer wg.Done()

			if err != nil {
				return
			}

			filename := filepath.Join(oldSrc, file.Name())
			newPath := filepath.Join(path, filepath.Base(filename))
			if err = os.Rename(filename, newPath); err != nil {
				return
			}
		}(&wg, oldSrc, dir)
	}

	wg.Wait()

	return
}
