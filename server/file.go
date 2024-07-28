package server

import (
	"context"
	"errors"
	"fmt"
	komputer "github.com/wittano/komputer/api/proto"
	"github.com/wittano/komputer/internal/audio"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log/slog"
	"os"
)

const downloadBufSize = 1024 * 1024

type fileServer struct {
	komputer.UnimplementedAudioFileServiceServer
}

func (fs fileServer) Download(req *komputer.DownloadFile, server komputer.AudioFileService_DownloadServer) error {
	if req == nil {
		return status.Error(codes.InvalidArgument, "download: missing required data")
	}

	path, err := audio.Path(filename(req))
	if err != nil {
		return status.Error(codes.NotFound, err.Error())
	}

	f, err := os.Open(path)
	if err != nil {
		slog.Error("failed find f "+path, err)
		return status.Error(codes.NotFound, err.Error())
	}
	defer f.Close()

	buf := make([]byte, downloadBufSize)
	for {
		select {
		case <-server.Context().Done():
			return status.Error(codes.Canceled, context.Canceled.Error())
		default:
		}

		n, err := f.Read(buf)
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		if err = server.Send(&komputer.FileBuffer{Content: buf, Size: uint64(n)}); err != nil {
			slog.Error("failed send chunk of file", err)
			return status.Error(codes.Internal, err.Error())
		}
	}

	return nil
}

func filename(req *komputer.DownloadFile) (name string) {
	if req == nil {
		return
	}

	var uuid *komputer.UUID
	uuid, name = req.GetUuid(), req.GetName()

	if uuid == nil {
		return
	}

	if name == "" {
		return string(uuid.Uuid)
	}

	return fmt.Sprintf("%s-%s", name, uuid.Uuid)
}
