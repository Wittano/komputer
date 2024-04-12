package audio

import (
	"bytes"
	"errors"
	"mime/multipart"
	"strings"
)

func ValidMp3File(file *multipart.FileHeader) (err error) {
	if !strings.HasSuffix(file.Filename, "mp3") {
		return errors.New("invalid audio extension")
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
	err = errors.New("invalid audio")

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
