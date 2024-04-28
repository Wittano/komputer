package audio

import (
	"context"
	"errors"
	"github.com/wittano/komputer/db"
	"github.com/wittano/komputer/web/settings"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
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

	for _, f := range files {
		go u.save(ctx, f, errCh, successCh)
	}

	count := len(files)
	for {
		select {
		case <-ctx.Done():
			return context.Canceled
		case err := <-errCh:
			return err
		case <-successCh:
			count -= 1
			break
		}

		if count <= 0 {
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

	service := DatabaseService{u.Db}
	for {
		select {
		case <-ctx.Done():
			errCh <- context.Canceled

			return
		default:
			const bufSize = 1 << 20 // 1MB buffer size

			_, err := io.CopyN(dest, src, bufSize)
			if errors.Is(err, io.EOF) {
				id, err := service.save(ctx, dest.Name())
				if err != nil {
					errCh <- err
				} else {
					newPath := strings.ReplaceAll(dest.Name(), file.Filename, id.Hex()+".mp3")
					if err = os.Rename(destPath, newPath); err != nil {
						errCh <- err
						os.Remove(destPath)

						return
					}

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
