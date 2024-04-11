package main

import (
	"flag"
	"github.com/wittano/komputer/web"
	"log"
)

const defaultConfigPath = "config.yml"

func main() {
	configPath := flag.String("config", defaultConfigPath, "Path to web console configuration file")
	flag.Parse()

	s, err := web.NewWebConsoleServer(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	log.Println("Server listening on 8080")
	log.Fatal(s.ListenAndServe())
}
