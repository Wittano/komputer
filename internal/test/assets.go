package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

const assetsDirKey = "ASSETS_DIR"

func CreateAssertDir(t *testing.T, n int) (err error) {
	dir := t.TempDir()

	if err = os.Setenv(assetsDirKey, dir); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < n; i++ {
		var f *os.File
		f, err = os.Create(filepath.Join(dir, fmt.Sprintf("test-%d.mp3", i)))
		if err != nil {
			return
		}
		err = f.Close()
	}

	return
}
