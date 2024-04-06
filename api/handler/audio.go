package handler

import (
	"context"
	"github.com/wittano/komputer/pkgs/voice"
	"net/http"
	"time"
)

// TODO export property to config file/environment variable
const maxFileSize = 8 * 1024 * 1024 // 8MB in bytes

func UploadNewAudio(res http.ResponseWriter, req *http.Request) (err error) {
	err = req.ParseMultipartForm(maxFileSize)
	if err != nil {
		return newInternalApiError(err)
	}

	// TODO export to property how much user can upload files
	if counts := len(req.MultipartForm.Value); counts < 1 || counts > 5 {
		return apiError{
			Status: http.StatusBadRequest,
			Msg:    "illegal uploaded files count",
		}
	}

	// TODO export uploading timeout to external properties
	const uploadingTimeout = time.Second * 2
	ctx, cancel := context.WithTimeout(req.Context(), uploadingTimeout)
	defer cancel()

	filesCount := len(req.MultipartForm.File)
	var (
		errCh        = make(chan error)
		successSigCh = make(chan struct{}, filesCount)
	)
	defer close(errCh)
	defer close(successSigCh)

	for k := range req.MultipartForm.File {
		go uploadFile(ctx, *req, k, errCh, successSigCh)
	}

	var (
		resError       error
		successCounter = filesCount
	)

	for {
		select {
		case <-ctx.Done():
			resError = context.Canceled
			break
		case err = <-errCh:
			resError = err
			break
		case <-successSigCh:
			successCounter -= 1
		}

		if successCounter <= 0 {
			res.WriteHeader(http.StatusCreated)

			break
		}
	}

	return resError
}

func uploadFile(ctx context.Context, req http.Request, name string, errCh chan<- error, successSig chan<- struct{}) {
	file, fileHeader, err := req.FormFile(name)
	if err != nil {
		errCh <- newInternalApiError(err)
		return
	}
	defer file.Close()

	if err := voice.ValidMp3File(fileHeader); err != nil {
		errCh <- newInternalApiError(err)
		return
	}

	if err = voice.UploadFile(ctx, fileHeader.Filename, file); err != nil {
		errCh <- newInternalApiError(err)

		return
	}

	successSig <- struct{}{}
}
