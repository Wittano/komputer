package server

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/wittano/komputer/api/proto"
	"github.com/wittano/komputer/audio"
	"github.com/wittano/komputer/test"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func createDownloadClient() (client pb.AudioFileServiceClient, server io.Closer, err error) {
	s, err := New(port)
	if err != nil {
		return nil, nil, err
	}

	go func() {
		if err := s.Start(); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	conn, err := grpc.NewClient("localhost:"+strconv.Itoa(port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		s.Close()
		return
	}

	server = closers{conn, s}
	client = pb.NewAudioFileServiceClient(conn)
	return
}

func TestFileServer_Download_ButFileDoesNotExists(t *testing.T) {
	client, closer, err := createDownloadClient()
	if err != nil {
		t.Fatal(err)
	}
	defer closer.Close()

	if err = test.CreateAssertDir(t, 1); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	name := "invalid_name"
	stream, err := client.Download(ctx, &pb.DownloadFile{Name: &name})
	if err != nil {
		t.Fatal(err)
	}
	defer stream.CloseSend()

	if _, err := stream.Recv(); status.Code(err) != codes.NotFound {
		t.Fatal("file invalid_name wasn't found in assets dir. Status code: " + strconv.Itoa(int(status.Code(err))))
	}
}

func TestFileServer_Download(t *testing.T) {
	t.Parallel()

	client, closer, err := createDownloadClient()
	if err != nil {
		t.Fatal(err)
	}
	defer closer.Close()

	if err = test.CreateAssertDir(t, 1); err != nil {
		t.Fatal(err)
	}

	assertDir := audio.AssertDir()
	dir, err := os.ReadDir(assertDir)
	if err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(assertDir, dir[0].Name())
	if err := fillTempFile(t, path); err != nil {
		t.Fatal(err)
	}

	base := filepath.Base(path)
	split := strings.Split(base, "-")
	if len(split) < 2 {
		t.Fatal("invalid split " + base)
	}

	uuid := pb.UUID{Uuid: []byte(strings.Join(split[1:], "-"))}
	data := []*pb.DownloadFile{
		{
			Name: &split[0],
		},
		{
			Uuid: &uuid,
		},
		{
			Name: &split[0],
			Uuid: &uuid,
		},
	}

	for _, d := range data {
		var name string
		if d.Name != nil {
			name = *d.Name
		}

		t.Run(fmt.Sprintf("download file with name: %s and uuid: %s", name, d.Uuid), func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			stream, err := client.Download(ctx, d)
			if err != nil {
				t.Fatal(err)
			}
			defer stream.CloseSend()

			for {
				select {
				case <-ctx.Done():
					return
				default:
				}

				recv, err := stream.Recv()
				if err != nil && !errors.Is(err, io.EOF) {
					t.Fatal(err)
				} else if errors.Is(err, io.EOF) || (recv != nil && recv.Size < downloadBufSize) {
					return
				}

				if recv.Size == 0 {
					t.Fatal("server didn't send chunk of file")
				}

				if len(recv.Content) == 0 {
					t.Fatalf("invalid size and number of bytes in content. Want: %d, Got: %d", recv.Size, len(recv.Content))
				}
			}
		})
	}
}
