package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type healthCheck struct {
	status    string    `json:"status"`
	timestamp time.Time `json:"timestamp"`
}

func HealthChecker(c echo.Context) error {
	status := healthCheck{
		status:    "ok",
		timestamp: time.Now(),
	}

	return c.JSON(http.StatusOK, status)
}
