package audio

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/wittano/komputer/db"
	"github.com/wittano/komputer/web/settings"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type UploadService struct {
	Db db.MongodbService
}

func (u UploadService) Upload(ctx context.Context, files []*multipart.FileHeader) error {
	var (
		errCh     = make(chan error)
		successCh = make(chan struct{})
	)
	defer close(errCh)
	defer close(successCh)

	filesCount := len(files)
	for _, f := range files {
		if err := validRequestedFile(*f); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid '%s' audio", f.Filename))
		}

		go u.save(ctx, f, errCh, successCh)
	}

	for {
		select {
		case <-ctx.Done():
			return context.Canceled
		case err := <-errCh:
			return err
		case <-successCh:
			filesCount -= 1
			break
		}

		if filesCount <= 0 {
			break
		}
	}

	return nil
}

func (u UploadService) save(ctx context.Context, file *multipart.FileHeader, errCh chan<- error, successSig chan<- struct{}) {
	select {
	case <-ctx.Done():
		errCh <- context.Canceled
		return
	default:
	}

	src, err := file.Open()
	if err != nil {
		errCh <- err

		return
	}
	defer src.Close()

	destPath := filepath.Join(settings.Config.AssetDir, file.Filename)
	dest, err := os.Create(destPath)
	if err != nil {
		errCh <- err

		return
	}
	defer dest.Close()

	for {
		select {
		case <-ctx.Done():
			errCh <- context.Canceled

			return
		default:
			const bufSize = 1 << 20 // 1MB buffer size

			_, err := io.CopyN(dest, src, bufSize)
			if errors.Is(err, io.EOF) {
				audioService := DatabaseService{u.Db}

				err = audioService.save(ctx, dest.Name())
				if err != nil {
					errCh <- err
				} else {
					successSig <- struct{}{}
				}

				return
			} else if err != nil {
				errCh <- err

				os.Remove(destPath)

				return
			}
		}
	}
}

func validRequestedFile(file multipart.FileHeader) error {
	if file.Size >= settings.Config.Upload.MaxFileSize {
		return fmt.Errorf("audio '%s' is too big", file.Filename)
	}

	if err := ValidMp3File(&file); err != nil {
		return err
	}

	destFile := filepath.Join(settings.Config.AssetDir, file.Filename)
	if _, err := os.Stat(destFile); err == nil {
		return os.ErrExist
	}

	return nil
}
