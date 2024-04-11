package handler

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/wittano/komputer/web/internal/file"
	"github.com/wittano/komputer/web/internal/settings"
	"net/http"
	"os"
	"path/filepath"
)

func UploadNewAudio(c echo.Context) (err error) {
	multipartForm, err := c.MultipartForm()
	if err != nil {
		return err
	}

	filesCount := len(multipartForm.File)
	if !settings.Config.CheckFileCountLimit(filesCount) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid number of uploaded files")
	}

	var (
		errCh     = make(chan error)
		successCh = make(chan struct{})
	)
	defer close(errCh)
	defer close(successCh)

	for k := range multipartForm.File {
		if err = validRequestedFile(k, *c.Request()); err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid '%s' file", k))
		}

		go uploadRequestedFile(c.Request().Context(), k, c.Request(), errCh, successCh)
	}

	for {
		select {
		case <-c.Request().Context().Done():
			return context.Canceled
		case err = <-errCh:
			return err
		case <-successCh:
			filesCount -= 1
			break
		}

		if filesCount <= 0 {
			c.Response().WriteHeader(http.StatusCreated)

			break
		}
	}

	return nil
}

func validRequestedFile(filename string, req http.Request) error {
	_, fileHeader, err := req.FormFile(filename)
	if err != nil {
		return err
	}

	if fileHeader.Size >= settings.Config.Upload.MaxFileSize {
		return fmt.Errorf("file '%s' is too big", filename)
	}

	if err = file.ValidMp3File(fileHeader); err != nil {
		return err
	}

	destFile := filepath.Join(settings.Config.AssetDir, filename)
	if _, err = os.Stat(destFile); err == nil {
		return os.ErrExist
	}

	return nil
}

func uploadRequestedFile(ctx context.Context, filename string, req *http.Request, errCh chan<- error, successSig chan<- struct{}) {
	select {
	case <-ctx.Done():
		errCh <- context.Canceled
		return
	default:
	}

	f, _, err := req.FormFile(filename)
	if err != nil {
		errCh <- err

		return
	}
	defer f.Close()

	dest, err := os.Create(filepath.Join(settings.Config.AssetDir, filename))
	if err != nil {
		errCh <- err

		return
	}
	defer dest.Close()

	if err = file.UploadFile(ctx, f, dest); err != nil {
		errCh <- err
		os.Remove(dest.Name())

		return
	}

	successSig <- struct{}{}
}
