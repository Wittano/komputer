package audio

import (
	"github.com/wittano/komputer/test"
	"testing"
)

func TestValidRequestedFile(t *testing.T) {
	if err := test.LoadDefaultConfig(t); err != nil {
		t.Fatal(err)
	}

	path, err := test.CreateTempAudioFiles(t)
	if err != nil {
		t.Fatal(err)
	}

	multipartFileHeader, err := test.CreateMultipartFileHeader(path)
	if err != nil {
		t.Fatal(err)
	}

	err = validRequestedFile(*multipartFileHeader)
	if err != nil {
		t.Fatal(err)
	}
}
