package server

import (
	"context"
	"errors"
	"fmt"
	komputer "github.com/wittano/komputer/api/proto"
	"github.com/wittano/komputer/audio"
	"io"
	"log/slog"
	"os"
)

type fileServer struct {
	komputer.UnimplementedAudioFileServiceServer
}

func (fs fileServer) Download(request *komputer.DownloadRequest, server komputer.AudioFileService_DownloadServer) error {
	if request == nil {
		return errors.New("download: missing request data")
	}

	path := audio.Path(filename(request))
	f, err := os.Open(path)
	if err != nil {
		slog.Error("failed find f "+path, err)
		return err
	}
	defer f.Close()

	ctx := server.Context()

	buf := make([]byte, 1024)
	for {
		select {
		case <-ctx.Done():
			return context.Canceled
		default:
		}

		n, err := f.Read(buf)
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err
		}

		if err = server.Send(&komputer.FileBuffer{Content: buf, Size: uint64(n)}); err != nil {
			return err
		}
	}

	return nil
}

func filename(request *komputer.DownloadRequest) (name string) {
	if request == nil {
		return
	}

	var uuid *komputer.UUID
	uuid, name = request.GetUuid(), request.GetName()

	if uuid == nil {
		return
	}

	return fmt.Sprintf("%s-%s", name, uuid)
}
