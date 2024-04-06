package voice

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// TODO added external property for uploading audio path
const uploadDir = ""

func UploadFile(ctx context.Context, filename string, file multipart.File) (err error) {
	path := filepath.Join(uploadDir, filename)
	if _, err := os.Stat(path); err == nil {
		return os.ErrExist
	}

	destFile, err := os.Create(path)
	if err != nil {
		return
	}
	defer destFile.Close()

	for {
		select {
		case <-ctx.Done():
			destFile.Close()

			if err = os.Remove(path); err == nil {
				err = context.Canceled
			}

			return
		default:
			const bufSize = 1024 * 1024

			_, err := io.CopyN(destFile, file, bufSize)
			if errors.Is(err, io.EOF) {
				return nil
			} else {
				return
			}
		}
	}
}
