package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/wittano/komputer/web/settings"
	"net/http"
)

func UpdateSettings(c echo.Context) error {
	var setting settings.Settings

	if err := c.Bind(&setting); err != nil {
		return err
	}

	return settings.Config.Update(setting)
}

func Settings(c echo.Context) error {
	return c.JSON(http.StatusOK, settings.Config)
}
