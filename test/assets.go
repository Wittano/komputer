package test

import (
	"fmt"
	"os"
	"testing"
)

const assetsDirKey = "ASSETS_DIR"

func CreateAssertDir(t *testing.T, n int) (err error) {
	dir := t.TempDir()

	if err := os.Setenv(assetsDirKey, dir); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < n; i++ {
		f, err := os.CreateTemp(dir, fmt.Sprintf("test-%d.*.mp3", i))
		if err != nil {
			t.Fatal(err)
		}
		err = f.Close()
	}

	return
}
