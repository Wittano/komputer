package handler

import (
	"bytes"
	"context"
	"errors"
	"github.com/wittano/komputer/web/internal/settings"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const testFileName = "test.mp3"

func createTempAudioFiles(t *testing.T) (string, error) {
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

func createMultipartFileHeader(filename string) (*multipart.FileHeader, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var buf bytes.Buffer

	formWriter := multipart.NewWriter(&buf)
	formPart, err := formWriter.CreateFormFile(testFileName, filepath.Base(filename))
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

	if file, ok := multipartForm.File[testFileName]; !ok || len(file) <= 0 {
		return nil, errors.New("failed create multipart audio")
	} else {
		return file[0], nil
	}
}

func loadDefaultConfig(t *testing.T) error {
	configFile := filepath.Join(t.TempDir(), "config.yml")

	if err := settings.Load(configFile); err != nil {
		return err
	}

	return settings.Config.Update(settings.Settings{AssetDir: filepath.Join(t.TempDir(), "assets")})
}

func TestValidRequestedFile(t *testing.T) {
	if err := loadDefaultConfig(t); err != nil {
		t.Fatal(err)
	}

	filePath, err := createTempAudioFiles(t)
	if err != nil {
		t.Fatal(err)
	}

	multipartFileHeader, err := createMultipartFileHeader(filePath)
	if err != nil {
		t.Fatal(err)
	}

	req := http.Request{
		MultipartForm: &multipart.Form{
			File: map[string][]*multipart.FileHeader{
				testFileName: {
					multipartFileHeader,
				},
			},
		},
	}

	err = validRequestedFile(testFileName, req)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUploadRequestedFile(t *testing.T) {
	if err := loadDefaultConfig(t); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	filePath, err := createTempAudioFiles(t)
	if err != nil {
		t.Fatal(err)
	}

	multipartFileHeader, err := createMultipartFileHeader(filePath)
	if err != nil {
		t.Fatal(err)
	}

	req := &http.Request{
		MultipartForm: &multipart.Form{
			File: map[string][]*multipart.FileHeader{
				testFileName: {
					multipartFileHeader,
				},
			},
		},
	}

	successCh := make(chan struct{})
	errCh := make(chan error)
	defer close(successCh)
	defer close(errCh)

	go uploadRequestedFile(ctx, testFileName, req, errCh, successCh)

	for {
		select {
		case <-ctx.Done():
			t.Fatal(context.Canceled)
		case err = <-errCh:
			t.Fatal(err)
		case <-successCh:
			return
		}
	}
}
