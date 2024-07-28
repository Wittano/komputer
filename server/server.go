package server

import (
	"context"
	pb "github.com/wittano/komputer/api/proto"
	"github.com/wittano/komputer/internal/mongodb"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
)

type Server struct {
	l          net.Listener
	cancelFunc context.CancelFunc
	serv       *grpc.Server
}

func (s *Server) Start() error {
	log.Printf("Server listing on port %serv\n", s.l.Addr())

	return s.serv.Serve(s.l)
}

func (s *Server) Close() error {
	s.serv.Stop()
	s.cancelFunc()
	return s.l.Close()
}

func New(port uint64) (*Server, error) {
	l, err := net.Listen("tcp", "0.0.0.0:"+strconv.FormatUint(port, 10))
	if err != nil {
		return nil, err
	}

	s := grpc.NewServer()
	ctx, cancel := context.WithCancel(context.Background())
	server := &Server{l, cancel, s}

	database := mongodb.NewMongodb(ctx)

	pb.RegisterJokeServiceServer(s, &jokeServer{Db: mongodb.Service{Db: database}})
	pb.RegisterAudioServiceServer(s, &audioServer{})
	pb.RegisterAudioFileServiceServer(s, &fileServer{})

	return server, nil
}
