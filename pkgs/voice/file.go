package voice

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"strings"
)

func UploadFile(ctx context.Context, src io.Reader, dest *os.File) error {
	for {
		select {
		case <-ctx.Done():
			return context.Canceled
		default:
			const bufSize = 1 << 20 // 1MB buffer size

			_, err := io.CopyN(dest, src, bufSize)
			if errors.Is(err, io.EOF) {
				return nil
			} else if err != nil {
				return err
			}
		}
	}
}

func ValidMp3File(file *multipart.FileHeader) (err error) {
	if !strings.HasSuffix(file.Filename, "mp3") {
		return errors.New("invalid file extension")
	}

	f, err := file.Open()
	if err != nil {
		return
	}
	defer f.Close()

	if err = checkAudioFileBinary(f); err != nil {
		return
	}
	return nil
}

func checkAudioFileBinary(f multipart.File) (err error) {
	const headerBytesSize = 2
	err = errors.New("invalid file")

	buf := make([]byte, headerBytesSize)
	n, err := f.Read(buf)
	if err != nil {
		return
	} else if n != headerBytesSize {
		return
	}

	mp3MagicNumbersHeader := []byte{0xff, 0xfb}
	if len(buf) != headerBytesSize && bytes.Equal(buf, mp3MagicNumbersHeader) {
		return
	}

	return nil
}
