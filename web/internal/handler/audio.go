package handler

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	globalApi "github.com/wittano/komputer/api"
	"github.com/wittano/komputer/db"
	"github.com/wittano/komputer/web/internal/api"
	"github.com/wittano/komputer/web/internal/audio"
	"github.com/wittano/komputer/web/settings"
	"mime/multipart"
	"net/http"
	"strconv"
)

func GetAudio(c echo.Context) error {
	ctx := c.Request().Context()
	service := audio.DatabaseService{Database: db.Mongodb(ctx)}

	info, err := service.Get(ctx, c.Param("id"))
	if err != nil {
		return err
	}

	c.Response().Header().Add("filename", info.Original)

	return c.File(info.Path)
}

func GetAudioFilesInfo(c echo.Context) error {
	ctx := c.Request().Context()
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		return err
	}

	service := audio.DatabaseService{Database: db.Mongodb(ctx)}

	ids, err := service.AudioFilesInfo(ctx, c.Param("type"), c.Param("value"), page)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, globalApi.GetAudioIdsResponse{Files: ids})
}

func RemoveAudio(c echo.Context) error {
	ctx := c.Request().Context()
	service := audio.DatabaseService{Database: db.Mongodb(ctx)}

	id := c.Param("id")
	if err := service.Delete(ctx, id); errors.Is(err, audio.NotFoundErr) {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("audio with id '%s' wasn't found", id))
	} else if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
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
		return errors.Join(echo.NewHTTPError(http.StatusBadRequest, "invalid number of uploaded files"), err)
	}

	ctx := c.Request().Context()
	service := audio.UploadService{Db: db.Mongodb(ctx)}

	if err := service.Upload(ctx, files); err != nil {
		var apiError api.Error

		if errors.As(err, &apiError) {
			c.Logger().Error(apiError.Err)
			err = apiError.HttpErr
		}

		return err
	}

	return c.NoContent(http.StatusCreated)
}
