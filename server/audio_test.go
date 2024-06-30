package server

import (
	"context"
	"errors"
	pb "github.com/wittano/komputer/api/proto"
	"github.com/wittano/komputer/audio"
	"github.com/wittano/komputer/test"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"strconv"
	"testing"
)

const (
	port = 3000
)

type closers []io.Closer

func (c closers) Close() (err error) {
	for _, f := range c {
		err = errors.Join(err, f.Close())
	}

	return
}

func createClient() (client pb.AudioServiceClient, server io.Closer, err error) {
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
	client = pb.NewAudioServiceClient(conn)
	return
}

func TestRemoveAudio(t *testing.T) {
	client, closer, err := createClient()
	if err != nil {
		t.Fatal(err)
	}
	defer closer.Close()

	if err = test.CreateAssertDir(t, 5); err != nil {
		t.Fatal(err)
	}

	paths, err := audio.Paths()
	if err != nil {
		t.Fatal(err)
	}

	if _, err := client.Remove(context.Background(), &pb.RemoveAudio{Name: paths}); err != nil {
		t.Fatal(err)
	}
}
