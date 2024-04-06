package main

import (
	"github.com/wittano/komputer/api/handler"
	"log"
	"net/http"
)

func main() {
	v1 := http.NewServeMux()
	v1.HandleFunc("POST /v1/audio", handler.MakeHttpHandler(handler.UploadNewAudio))

	server := http.Server{
		Addr:    ":8080",
		Handler: v1,
	}
	defer server.Close()

	log.Fatal(server.ListenAndServe())
}
