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

// TODO Verification of service structe
func (fs fileServer) Download(req *komputer.DownloadRequest, server komputer.AudioFileService_DownloadServer) error {
	if req == nil {
		return errors.New("download: missing req data")
	}

	path := audio.Path(filename(req))
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

func filename(req *komputer.DownloadRequest) (name string) {
	if req == nil {
		return
	}

	var uuid *komputer.UUID
	uuid, name = req.GetUuid(), req.GetName()

	if uuid == nil {
		return
	}

	return fmt.Sprintf("%s-%s", name, uuid)
}
