package server

import (
	"bufio"
	"context"
	"errors"
	pb "github.com/wittano/komputer/api/proto"
	"github.com/wittano/komputer/audio"
	"github.com/wittano/komputer/test"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
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

func createAudioClient() (client pb.AudioServiceClient, server io.Closer, err error) {
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
	client, closer, err := createAudioClient()
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

	for i, p := range paths {
		paths[i] = filepath.Base(p)
	}

	if _, err := client.Remove(context.Background(), &pb.RemoveAudio{Name: paths}); err != nil {
		t.Fatal(err)
	}
}

func TestUploadFile(t *testing.T) {
	client, closer, err := createAudioClient()
	if err != nil {
		t.Fatal(err)
	}
	defer closer.Close()

	if err = test.CreateAssertDir(t, 0); err != nil {
		t.Fatal(err)
	}

	path := t.TempDir() + "test.mp3"
	if err := fillTempFile(t, path); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := client.Add(ctx)

	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}

	filename := filepath.Base(path)
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		stream.Send(&pb.Audio{Info: &pb.AudioInfo{Name: filename, Type: pb.FileFormat_MP3}, Chunk: scan.Bytes()})
	}

	res, err := stream.CloseAndRecv()
	if err != nil && !errors.Is(err, io.EOF) {
		t.Fatal(err)
	}

	path, err = audio.Path(res.Filename)
	if err != nil {
		t.Fatal(err)
	}

	if s, err := os.Stat(path); err != nil || s.Size() == 0 {
		t.Fatalf("failed upload file: %v", err)
	}
}

func TestUploadFile_FileAlreadyExists(t *testing.T) {
	client, closer, err := createAudioClient()
	if err != nil {
		t.Fatal(err)
	}
	defer closer.Close()

	if err = test.CreateAssertDir(t, 0); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(audio.AssertDir(), "test")
	f, createErr := os.Create(path)
	if errors.Join(err, createErr) != nil {
		t.Fatal(err)
	}
	f.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := client.Add(ctx)

	if err := stream.Send(&pb.Audio{Info: &pb.AudioInfo{Name: filepath.Base(path)}, Chunk: []byte{}}); err != nil {
		t.Fatal(err)
	}

	if _, err = stream.CloseAndRecv(); status.Code(err) != codes.AlreadyExists {
		t.Fatalf("file %s shouldn't be uploaded", path)
	}
}

func fillTempFile(t *testing.T, path string) error {
	if path == "" {
		path = t.TempDir() + "test.mp3"
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	for i := 0; i < 100; i++ {
		if _, err := f.WriteString(strconv.Itoa(rand.Int()) + "\n"); err != nil {
			return err
		}
	}

	return nil
}
