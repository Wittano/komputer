package settings

import (
	"os"
	"path/filepath"
	"testing"
)

// Create temporary dictionary with one nested directory, that will be a new destination.
// Optional, if temp parameter isn't empty, then creates file in original source directory
func createOldAndNewSourceDir(t *testing.T, temp string) (old string, new string, err error) {
	old = t.TempDir()
	new = filepath.Join(old, "newDir")

	if err = os.Mkdir(new, 0700); err != nil {
		return
	}

	if temp != "" {
		f, err := os.CreateTemp(old, temp)
		if err != nil {
			return "", "", err
		}
		defer f.Close()
	}

	return
}

func TestMoveToNewLocation(t *testing.T) {
	src, dest, err := createOldAndNewSourceDir(t, "test*.mp3")
	if err != nil {
		t.Fatal(err)
	}

	if err = moveToNewLocation(src, dest); err != nil {
		t.Fatal(err)
	}
}

func TestMoveToNewLocation_SrcDoesNotHaveMp3Files(t *testing.T) {
	src, dest, err := createOldAndNewSourceDir(t, "test*.yml")
	if err != nil {
		t.Fatal(err)
	}

	if err = moveToNewLocation(src, dest); err != nil {
		t.Fatal(err)
	}
}

func TestMoveToNewLocation_SourceIsEmpty(t *testing.T) {
	src, dest, err := createOldAndNewSourceDir(t, "")
	if err != nil {
		t.Fatal(err)
	}

	if err = moveToNewLocation(src, dest); err != nil {
		t.Fatal(err)
	}
}
