package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/wittano/komputer/web/settings"
	"net/http"
)

func UpdateSettings(c echo.Context) error {
	var newSetting settings.Settings

	if err := c.Bind(&newSetting); err != nil {
		return err
	}

	return settings.Config.Update(newSetting)
}

func GetSettings(c echo.Context) error {
	return c.JSON(http.StatusOK, settings.Config)
}
