package server

import (
	"context"
	komputer "github.com/wittano/komputer/api/proto"
	"github.com/wittano/komputer/db"
	"github.com/wittano/komputer/db/joke"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
	"sync"
)

type Server struct {
	l    net.Listener
	serv *grpc.Server
}

func (s Server) Start() error {
	log.Printf("Server listing on port %serv\n", s.l.Addr())

	return s.serv.Serve(s.l)
}

func (s Server) Close() error {
	s.serv.Stop()
	return s.l.Close()
}

func New(port uint64) (Server, error) {
	l, err := net.Listen("tcp", strconv.FormatUint(port, 10))
	if err != nil {
		return Server{}, err
	}

	return Server{l, newGRPGServer()}, nil
}

func newGRPGServer() *grpc.Server {
	s := grpc.NewServer()

	ctx := context.Background()
	mongodb := db.Mongodb(ctx)

	komputer.RegisterJokeServiceServer(s, &jokeServer{Db: joke.Database{Mongodb: mongodb}})
	komputer.RegisterAudioServiceServer(s, &audioServer{m: new(sync.Mutex)})
	komputer.RegisterAudioFileServiceServer(s, &fileServer{})

	return s
}
