package assets

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

const assertDir = "assets"

func GetAudioPaths(filename string) ([]string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	basePath := path.Join(dir, assertDir, filename)
	dir = filepath.Dir(basePath)
	name := filepath.Base(basePath)

	ls, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var paths = make([]string, len(ls))

	for _, d := range ls {
		fName := d.Name()
		if !d.IsDir() && strings.Contains(fName, name) {
			paths = append(paths, filepath.Join(dir, fName))
		}
	}

	return paths, nil
}
