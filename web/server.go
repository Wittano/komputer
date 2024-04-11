package web

import (
	"github.com/wittano/komputer/web/internal/handler"
	"github.com/wittano/komputer/web/internal/settings"
	"net/http"
)

func NewWebConsoleServer(configPath string) (*http.Server, error) {
	if err := settings.Load(configPath); err != nil {
		return nil, err
	}

	v1 := http.NewServeMux()
	v1.HandleFunc("POST /api/v1/audio", handler.MakeHttpHandler(handler.UploadNewAudio))
	v1.HandleFunc("GET /api/v1/setting", handler.MakeHttpHandler(handler.GetSettings))
	v1.HandleFunc("PUT /api/v1/setting", handler.MakeHttpHandler(handler.UpdateSettings))

	return &http.Server{
		Addr:    ":8080",
		Handler: v1,
	}, nil
}
