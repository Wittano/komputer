package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	pb "github.com/wittano/komputer/api/proto"
	"github.com/wittano/komputer/audio"
	"google.golang.org/protobuf/types/known/emptypb"
	"os"
	"sync"
)

type audioServer struct {
	m *sync.Mutex
	pb.UnimplementedAudioServiceServer
}

func (a audioServer) List(pagination *pb.Pagination, server pb.AudioService_ListServer) (err error) {
	page := paginationOrDefault(pagination)

	paths, err := audio.PathsWithPagination(page.Page, page.Size)
	if err != nil {
		return
	}

	for _, path := range paths {
		err = errors.Join(err, server.Send(&pb.AudioInfo{Name: path, Type: pb.FileFormat_MP3}))
	}

	return
}

// TODO Verification upload request structure
func (a audioServer) Add(server pb.AudioService_AddServer) error {
	for {
		au, err := server.Recv()
		if err != nil {
			return err
		}

		path := audio.Path(fmt.Sprintf("%s-%s.%s", au.Info.Name, uuid.NewString(), au.Info.Type.String()))

		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return err
		}
		f.Close()

		a.m.Lock()
		if _, err := f.Write(au.Chunk); err != nil {
			a.m.Unlock()
			return err
		}
		a.m.Unlock()
	}
}

func (a audioServer) Remove(_ context.Context, req *pb.RemoveAudioRequest) (e *emptypb.Empty, err error) {
	e = &emptypb.Empty{}
	if req == nil {
		return
	}

	for _, query := range req.Name {
		err = errors.Join(err, os.Remove(query))
	}

	return
}
