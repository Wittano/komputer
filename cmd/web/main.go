package main

import (
	"flag"
	"github.com/wittano/komputer/api/handler"
	"github.com/wittano/komputer/pkgs/settings"
	"log"
	"net/http"
)

func main() {
	configPath := flag.String("config", settings.DefaultSettingsPath, "Path to web console configuration file")
	flag.Parse()

	if err := settings.Load(*configPath); err != nil {
		log.Fatal(err)
	}

	v1 := http.NewServeMux()
	v1.HandleFunc("POST /api/v1/audio", handler.MakeHttpHandler(handler.UploadNewAudio))
	v1.HandleFunc("GET /api/v1/setting", handler.MakeHttpHandler(handler.GetSettings))
	v1.HandleFunc("PUT /api/v1/setting", handler.MakeHttpHandler(handler.UpdateSettings))

	server := http.Server{
		Addr:    ":8080",
		Handler: v1,
	}
	defer server.Close()

	log.Println("Server listening on 8080")
	log.Fatal(server.ListenAndServe())
}
