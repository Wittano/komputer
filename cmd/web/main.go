package main

import (
	"flag"
	"github.com/wittano/komputer/web"
	"log"
)

const defaultConfigPath = "config.yml"

func main() {
	configPath := flag.String("config", defaultConfigPath, "Path to web console configuration audio")
	flag.Parse()

	e, err := web.NewWebConsoleServer(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	defer e.Close()

	e.Logger.Fatal(e.Start(":8080"))
}
