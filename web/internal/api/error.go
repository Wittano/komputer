package api

import "github.com/labstack/echo/v4"

type Error struct {
	HttpErr *echo.HTTPError
	Err     error
}

func (a Error) Error() string {
	return a.Err.Error()
}
