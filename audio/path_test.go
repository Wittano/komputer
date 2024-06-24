package audio

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestAudioIDs_AssetDirHasEmptyDirs(t *testing.T) {
	dir := t.TempDir()

	for i := 0; i < 5; i++ {
		os.Mkdir(filepath.Join(dir, strconv.Itoa(i)), 0700)
	}

	if err := os.Setenv(assetsDirKey, dir); err != nil {
		t.Fatal(err)
	}

	paths, err := Paths()
	if err != nil {
		t.Fatal(err)
	}

	if len(paths) != 0 {
		t.Fatal("something was found in empty directory")
	}
}

func TestAudioIDs(t *testing.T) {
	dir := t.TempDir()

	const expectedFilesNumber = 5
	for i := 0; i < expectedFilesNumber; i++ {
		f, err := os.CreateTemp(dir, fmt.Sprintf("test-%d.*.mp3", i))
		if err != nil {
			t.Fatal(err)
		}
		f.Close()
	}

	if err := os.Setenv(assetsDirKey, dir); err != nil {
		t.Fatal(err)
	}

	paths, err := Paths()
	if err != nil {
		t.Fatal(err)
	}

	if len(paths) != expectedFilesNumber {
		t.Fatalf("missing audios IDs. Expected '%d', Result: '%d'", expectedFilesNumber, len(paths))
	}

	for i, id := range paths {
		fileID := strings.Split(id, ".")[0]

		if fileID != "test-"+strconv.Itoa(i) {
			t.Fatalf("invalid ID. Expected: '%d', Result: '%s'", i, fileID)
		}
	}
}
