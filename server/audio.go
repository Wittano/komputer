package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	komputer "github.com/wittano/komputer/api/proto"
	"github.com/wittano/komputer/audio"
	"google.golang.org/protobuf/types/known/emptypb"
	"os"
	"sync"
)

type audioServer struct {
	komputer.UnimplementedAudioServiceServer
}

func (a audioServer) List(pagination *komputer.Pagination, server komputer.AudioService_ListServer) (err error) {
	page := paginationOrDefault(pagination)

	paths, err := audio.PathsWithPagination(page.Page, page.Size)
	if err != nil {
		return
	}

	for _, path := range paths {
		err = errors.Join(err, server.Send(&komputer.AudioInfo{Name: path, Type: komputer.FileFormat_MP3}))
	}

	return
}

func (a audioServer) Add(server komputer.AudioService_AddServer) error {
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

		m := sync.Mutex{}

		m.Lock()
		if _, err := f.Write(au.Chunk); err != nil {
			m.Unlock()
			return err
		}
		m.Unlock()
	}
}

func (a audioServer) Remove(ctx context.Context, request *komputer.NameOrIdAudioRequest) (e *emptypb.Empty, err error) {
	e = &emptypb.Empty{}
	if request == nil || len(request.Query) == 0 {
		return e, nil
	}

	var wg sync.WaitGroup
	for _, query := range request.Query {
		if query == nil {
			continue
		}

		wg.Add(1)
		go func(q *komputer.FileQuery) {
			err = remove(ctx, &wg, q)
		}(query)
	}
	wg.Wait()

	return
}

func remove(ctx context.Context, wg *sync.WaitGroup, query *komputer.FileQuery) error {
	defer wg.Done()
	select {
	case <-ctx.Done():
		return context.Canceled
	default:
	}

	uid, name := query.GetUuid(), query.GetName()

	var path string
	if uid != nil && name != "" {
		path = audio.Path(fmt.Sprintf("%s-%s", name, uid.Uuid))
	}

	if path == "" {
		return os.ErrNotExist
	}

	return os.Remove(path)
}
