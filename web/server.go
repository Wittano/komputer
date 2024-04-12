package web

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/wittano/komputer/web/internal/handler"
	"github.com/wittano/komputer/web/settings"
)

func NewWebConsoleServer(configPath string) (*echo.Echo, error) {
	if err := settings.Load(configPath); err != nil {
		return nil, err
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/api/v1/audio", handler.UploadNewAudio)
	e.GET("/api/v1/audio/:id", handler.GetAudio)

	e.GET("/api/v1/setting", handler.GetSettings)
	e.PUT("/api/v1/setting", handler.UpdateSettings)

	return e, nil
}
