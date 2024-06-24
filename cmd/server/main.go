package main

import (
	"flag"
	"github.com/wittano/komputer/server"
	"log"
)

func main() {
	port := flag.Uint64("port", 8080, "Server TCP port")

	s, err := server.New(*port)
	if err != nil {
		log.Fatalf("failed initialized server: %s", err)
	}
	defer s.Close()

	if err = s.Start(); err != nil {
		log.Fatal(err)
	}
}
