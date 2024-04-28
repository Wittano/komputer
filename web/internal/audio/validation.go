package audio

import (
	"bytes"
	"errors"
	"mime/multipart"
	"strings"
)

func ValidMp3File(header *multipart.FileHeader) (err error) {
	if !strings.HasSuffix(header.Filename, "mp3") {
		return errors.New("invalid audio extension")
	}

	f, err := header.Open()
	if err != nil {
		return
	}
	defer f.Close()

	return checkAudioFileBinary(f)
}

func checkAudioFileBinary(f multipart.File) (err error) {
	const headerBytesSize = 2
	err = errors.New("invalid audio")

	buf := make([]byte, headerBytesSize)
	n, err := f.Read(buf)
	if err != nil {
		return
	} else if n != headerBytesSize {
		return
	}

	// Magic numbers, that MP3 starts them
	headers := []byte{0xff, 0xfb}
	if len(buf) != headerBytesSize && bytes.Equal(buf, headers) {
		return
	}

	return nil
}
