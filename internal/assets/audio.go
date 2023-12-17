package assets

import (
	"errors"
	"math/rand"
	"os"
	"path"
	"path/filepath"
)

const assertDir = "assets"

func Audios() ([]string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(path.Join(dir, assertDir))
	if err != nil {
		return nil, err
	}

	if len(files) <= 0 {
		return nil, errors.New("assert directory is empty")
	}

	paths := make([]string, len(files))
	for i, f := range files {
		paths[i] = filepath.Join(dir, assertDir, f.Name())
	}

	return paths, nil
}

func RandomAudio() (string, error) {
	paths, err := Audios()
	if err != nil {
		return "", err
	}

	return paths[rand.Int()%len(paths)], nil
}
