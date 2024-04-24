package voice

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestAudioIDsButDirectoryHasEmptyDirs(t *testing.T) {
	dir := t.TempDir()

	for i := 0; i < 5; i++ {
		os.Mkdir(filepath.Join(dir, strconv.Itoa(i)), 0700)
	}

	if err := os.Setenv(CacheDirAudioKey, dir); err != nil {
		t.Fatal(err)
	}

	ids, err := AudioIDs()
	if err != nil {
		t.Fatal(err)
	}

	if len(ids) != 0 {
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

	if err := os.Setenv(CacheDirAudioKey, dir); err != nil {
		t.Fatal(err)
	}

	ids, err := AudioIDs()
	if err != nil {
		t.Fatal(err)
	}

	if len(ids) != expectedFilesNumber {
		t.Fatalf("missing audios IDs. Expected '%d', Result: '%d'", expectedFilesNumber, len(ids))
	}

	for i, id := range ids {
		fixedID := strings.Split(id, ".")[0]

		if fixedID != strconv.Itoa(i) {
			t.Fatalf("invalid ID. Expected: '%d', Result: '%s'", i, id)
		}
	}
}
