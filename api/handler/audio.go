package handler

import (
	"context"
	"fmt"
	"github.com/wittano/komputer/pkgs/settings"
	"github.com/wittano/komputer/pkgs/voice"
	"net/http"
	"os"
	"path/filepath"
)

const oneMegaByte = 1 << 20 // 8MB in bytes

func UploadNewAudio(res http.ResponseWriter, req *http.Request) (err error) {
	err = req.ParseMultipartForm(settings.Config.Upload.MaxFileSize * oneMegaByte)
	if err != nil {
		return newInternalApiError(err)
	}

	filesCount := len(req.MultipartForm.File)
	if !settings.Config.CheckFileCountLimit(filesCount) {
		return apiError{
			Status: http.StatusBadRequest,
			Msg:    "illegal uploaded files count",
		}
	}

	var (
		errCh     = make(chan error)
		successCh = make(chan struct{}, filesCount)
	)
	defer close(errCh)
	defer close(successCh)

	for k := range req.MultipartForm.File {
		err = validRequestedFile(k, *req)
		if err != nil {
			return apiError{
				Status: http.StatusBadRequest,
				Msg:    fmt.Sprintf("invalid '%s' file", k),
				Err:    err,
			}
		}

		go uploadRequestedFile(req.Context(), k, req, errCh, successCh)
	}

	successCounter := filesCount

	for {
		select {
		case <-req.Context().Done():
			return context.Canceled
		case err = <-errCh:
			return err
		case <-successCh:
			successCounter -= 1
			break
		}

		if successCounter <= 0 {
			res.WriteHeader(http.StatusCreated)

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

	if err = voice.ValidMp3File(fileHeader); err != nil {
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
		errCh <- newInternalApiError(err)

		return
	}
	defer f.Close()

	dest, err := os.Create(filepath.Join(settings.Config.AssetDir, filename))
	if err != nil {
		errCh <- newInternalApiError(err)

		return
	}
	defer dest.Close()

	if err = voice.UploadFile(ctx, f, dest); err != nil {
		errCh <- newInternalApiError(err)
		os.Remove(dest.Name())

		return
	}

	successSig <- struct{}{}
}
