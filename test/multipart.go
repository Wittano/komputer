package test

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"
)

func CreateTempAudioFiles(t *testing.T) (string, error) {
	dir := t.TempDir()
	f, err := os.CreateTemp(dir, "test.*.mp3")
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = f.Write([]byte{0xff, 0xfb})
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}

func CreateMultipartFileHeader(path string) (*multipart.FileHeader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var buf bytes.Buffer

	w := multipart.NewWriter(&buf)
	filename := filepath.Base(path)
	formWriter, err := w.CreateFormFile(filename, filepath.Base(path))
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(formWriter, f); err != nil {
		return nil, err
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	r := multipart.NewReader(bytes.NewReader(buf.Bytes()), w.Boundary())
	form, err := r.ReadForm(1 << 20)
	if err != nil {
		return nil, err
	}

	if file, ok := form.File[filename]; !ok || len(file) <= 0 {
		return nil, errors.New("failed create multipart audio")
	} else {
		return file[0], nil
	}
}
