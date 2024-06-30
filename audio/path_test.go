package audio

import (
	"errors"
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
	const expectedFilesNumber = 5
	if err := createTestAssertDir(t, expectedFilesNumber); err != nil {
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

func TestPathsWithPagination_ButEmptyDictionary(t *testing.T) {
	if err := os.Setenv(assetsDirKey, t.TempDir()); err != nil {
		t.Fatal(err)
	}

	res, err := PathsWithPagination(0, 10)
	if err == nil {
		t.Fatalf("Assert dictionary was found: %s", os.Getenv(assetsDirKey))
	}

	if len(res) != 0 {
		t.Fatalf("Something was found in assert dictionary, but it doesn't expect. %v", res)
	}
}

func TestPathsWithPagination_PageIsOverANumberOfFiles(t *testing.T) {
	const expectedFilesNumber = 5
	if err := createTestAssertDir(t, expectedFilesNumber); err != nil {
		t.Fatal(err)
	}

	res, err := PathsWithPagination(10, 10)
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 0 {
		t.Fatalf("Something was added to list, but page and size are over number of files. %v", res)
	}
}

func TestPathsWithPagination(t *testing.T) {
	const expectedFilesNumber = 50
	if err := createTestAssertDir(t, expectedFilesNumber); err != nil {
		t.Fatal(err)
	}

	res, err := PathsWithPagination(2, 10)
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 10 {
		t.Fatalf("Files wasn't found but shoule be")
	}

	const startFileSuffix = 20
	for i := startFileSuffix; i < 30; i++ {
		f := res[i-startFileSuffix]
		if _, err = os.Stat(f); !errors.Is(err, os.ErrNotExist) {
			t.Fatal(err)
		}
	}
}

func TestPathsWithPagination_AssertDirNotFound(t *testing.T) {
	if _, err := PathsWithPagination(0, 10); err == nil {
		t.Fatalf("Assert dictionary was found")
	}
}

func createTestAssertDir(t *testing.T, n int) (err error) {
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
