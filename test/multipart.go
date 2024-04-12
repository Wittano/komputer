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

	formWriter := multipart.NewWriter(&buf)
	filename := filepath.Base(path)
	formPart, err := formWriter.CreateFormFile(filename, filepath.Base(path))
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(formPart, f); err != nil {
		return nil, err
	}

	err = formWriter.Close()
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(buf.Bytes())
	formReader := multipart.NewReader(reader, formWriter.Boundary())

	multipartForm, err := formReader.ReadForm(1 << 20)
	if err != nil {
		return nil, err
	}

	if file, ok := multipartForm.File[filename]; !ok || len(file) <= 0 {
		return nil, errors.New("failed create multipart audio")
	} else {
		return file[0], nil
	}
}
