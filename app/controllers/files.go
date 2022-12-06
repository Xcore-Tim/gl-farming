package controllers

import (
	"gl-farming/app/constants/files"
	"gl-farming/app/services"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type FileController struct {
	Services services.AppServices
}

type Response struct {
	Success        bool   `json:"success"`
	FileMessage    string `json:"fileMessage"`
	RequestMessage string `json:"requestMessage"`
	Name           string `json:"name"`
}

func NewFileController(appServices services.AppServices) TableController {
	return TableController{
		Services: appServices,
	}
}

func (ctrl FileController) Upload(c echo.Context) error {

	var file *os.File

	response := Response{
		Success:     false,
		FileMessage: "keep sending chunks",
	}

	uplFile, upload, err := c.Request().FormFile("file")

	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	contentRange := c.Request().Header.Get("Content-Range")
	rangeAndSize := strings.Split(contentRange, "/")
	rangeParts := strings.Split(rangeAndSize[0], "-")

	maxRange, err := strconv.Atoi(rangeParts[1])
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	fileSize, err := strconv.Atoi(rangeAndSize[1])
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if fileSize > files.MaxSize {
		return c.String(http.StatusBadRequest, "error: file size must be less than 100 MB")
	}

	if file == nil {
		filePath := files.Static + "/" + upload.Filename
		file, err = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
	}

	if _, err := io.Copy(file, uplFile); err != nil {
		return c.String(http.StatusBadRequest, "error: error writing to a file")
	}

	file.Close()

	response.Success = true

	if maxRange >= fileSize-1 {

		fileExt := filepath.Ext(c.Request().FormValue("fileName"))

		fileName, filePath, err := ctrl.Services.Files.CreateNewFile(upload, fileExt)

		if err != nil {
			return c.String(http.StatusConflict, err.Error())

		}

		oid := c.QueryParam("oid")

		oldFile, isFound := ctrl.Services.Files.CheckPreviousFile(c, oid)

		if isFound {
			if err := ctrl.Services.Files.DeletePreviousFile(oldFile); err != nil {
				return c.String(http.StatusConflict, err.Error())
			}
		}

		if err := ctrl.Services.Files.UpdateDownloadLink(c, fileName, oid); err != nil {
			response.RequestMessage = "Error updating account request"
		}

		response.Success = true
		response.Name = filePath
		response.FileMessage = "File uploaded successfuly"
		response.RequestMessage = "Account request updated successfuly"

		return c.JSON(http.StatusOK, response)

	}

	return c.String(http.StatusOK, "success")
}
