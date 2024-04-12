package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/wittano/komputer/db"
	"github.com/wittano/komputer/web/internal/audio"
	"github.com/wittano/komputer/web/settings"
	"mime/multipart"
	"net/http"
)

func GetAudio(c echo.Context) error {
	id := c.Param("id")

	ctx := c.Request().Context()
	service := audio.DatabaseService{Database: db.Mongodb(ctx)}

	info, err := service.Get(ctx, id)
	if err != nil {
		return err
	}

	return c.File(info.Path)
}

func UploadNewAudio(c echo.Context) (err error) {
	multipartForm, err := c.MultipartForm()
	if err != nil {
		return err
	}

	var files []*multipart.FileHeader

	for _, v := range multipartForm.File {
		files = append(files, v...)
	}

	filesCount := len(files)
	if !settings.Config.CheckFileCountLimit(filesCount) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid number of uploaded files")
	}

	ctx := c.Request().Context()
	service := audio.UploadService{Db: db.Mongodb(ctx)}

	if err := service.Upload(ctx, files); err != nil {
		return err
	}

	return c.String(http.StatusCreated, "")
}
