package settings

import (
	"os"
	"path/filepath"
	"testing"
)

func createOldAndNewSourceDir(t *testing.T, tempName string) (old string, new string, err error) {
	old = t.TempDir()
	new = filepath.Join(old, "newDir")

	if err = os.Mkdir(new, 0700); err != nil {
		return
	}

	if tempName != "" {
		f, err := os.CreateTemp(old, tempName)
		if err != nil {
			return "", "", err
		}
		defer f.Close()
	}

	return
}

func TestMoveToNewLocation(t *testing.T) {
	oldSrc, newDest, err := createOldAndNewSourceDir(t, "test*.mp3")
	if err != nil {
		t.Fatal(err)
	}

	if err = moveToNewLocation(oldSrc, newDest); err != nil {
		t.Fatal(err)
	}
}

func TestMoveToNewLocationButSourceHasOnlyNonMp3Files(t *testing.T) {
	oldSrc, newDest, err := createOldAndNewSourceDir(t, "test*.yml")
	if err != nil {
		t.Fatal(err)
	}

	if err = moveToNewLocation(oldSrc, newDest); err != nil {
		t.Fatal(err)
	}
}

func TestMoveToNewLocationButSourceIsEmpty(t *testing.T) {
	oldSrc, newDest, err := createOldAndNewSourceDir(t, "")
	if err != nil {
		t.Fatal(err)
	}

	if err = moveToNewLocation(oldSrc, newDest); err != nil {
		t.Fatal(err)
	}
}
