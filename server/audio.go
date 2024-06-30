package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	pb "github.com/wittano/komputer/api/proto"
	"github.com/wittano/komputer/audio"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type audioServer struct {
	pb.UnimplementedAudioServiceServer
}

func (a audioServer) List(pagination *pb.Pagination, server pb.AudioService_ListServer) (err error) {
	page := paginationOrDefault(pagination)

	paths, err := audio.PathsWithPagination(page.Page, page.Size)
	if err != nil {
		return status.Error(codes.NotFound, err.Error())
	}

	for _, path := range paths {
		err = errors.Join(err, server.Send(&pb.AudioInfo{Name: path, Type: pb.FileFormat_MP3}))
	}

	return
}

// TODO Add file validation
func (a audioServer) Add(server pb.AudioService_AddServer) error {
	id := uuid.NewString()
	var path string

	au, err := server.Recv()
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	if _, err := audio.Path(au.Info.Name); err == nil {
		return status.Error(codes.AlreadyExists, fmt.Sprintf("file %s already exists", au.Info.Name))
	}

	if path == "" {
		path = filepath.Join(audio.AssertDir(), fmt.Sprintf("%s-%s.%s", au.Info.Name, id, strings.ToLower(au.Info.Type.String())))
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return status.Error(codes.FailedPrecondition, err.Error())
	}
	defer f.Close()

	for {
		select {
		case <-server.Context().Done():
			return status.Error(codes.Canceled, context.Canceled.Error())
		default:
		}

		if len(au.Chunk) == 0 {
			break
		}

		if _, err = f.Write(au.Chunk); err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		if au, err = server.Recv(); err != nil && errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}

	stat, err := os.Stat(path)
	if err != nil {
		return status.Error(codes.NotFound, "uploaded file wasn't found")
	}

	return server.SendAndClose(&pb.UploadAudioResponse{Size: uint64(stat.Size()), Filename: filepath.Base(path)})
}

func (a audioServer) Remove(_ context.Context, req *pb.RemoveAudio) (e *emptypb.Empty, err error) {
	e = &emptypb.Empty{}
	if req == nil {
		return
	}

	for _, query := range req.Name {
		path, err := audio.Path(query)
		if err != nil {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		if err = os.Remove(path); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return
}
