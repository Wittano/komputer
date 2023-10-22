package assets

import (
	"errors"
	"os"
	"path"
)

const assertDir = "assets"

func GetAudioPath(filename string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	p := path.Join(dir, assertDir, filename)
	if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
		return "", err
	}

	return p, nil
}
