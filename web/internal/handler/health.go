package handler

import (
	"github.com/labstack/echo/v4"
	globalApi "github.com/wittano/komputer/api"
	"net/http"
	"time"
)

func HealthChecker(c echo.Context) error {
	return c.JSON(http.StatusOK, globalApi.HealthCheck{
		Status:    "ok",
		Timestamp: time.Now(),
	})
}
