package test

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"testing"
)

const assetsDirKey = "ASSETS_DIR"

func CreateAssertDir(t *testing.T, n int) (err error) {
	dir := t.TempDir()

	if err := os.Setenv(assetsDirKey, dir); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < n; i++ {
		f, err := os.Create(filepath.Join(dir, fmt.Sprintf("test-%s.mp3", uuid.NewString())))
		if err != nil {
			return err
		}
		err = f.Close()
	}

	return
}
