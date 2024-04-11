package web

import (
	"github.com/labstack/echo/v4"
	"github.com/wittano/komputer/web/internal/handler"
	"github.com/wittano/komputer/web/internal/settings"
)

func NewWebConsoleServer(configPath string) (*echo.Echo, error) {
	if err := settings.Load(configPath); err != nil {
		return nil, err
	}

	e := echo.New()

	e.POST("/api/v1/audio", handler.UploadNewAudio)
	e.GET("/api/v1/setting", handler.GetSettings)
	e.PUT("PUT /api/v1/setting", handler.UpdateSettings)

	return e, nil
}
