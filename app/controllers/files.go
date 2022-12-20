package controllers

import (
	"gl-farming/app/constants/files"
	googleDriveAPI "gl-farming/app/google-drive"
	"gl-farming/app/models"
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

	// contentRange := c.Request().Header.Get("Content-Range")
	contentRange := c.Request().Header.Get("C-Range")
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

		driveID, driveLink := googleDriveAPI.Upload(fileName, filePath)

		oid := c.QueryParam("oid")

		ctrl.Services.Files.DeleteFile(filePath)

		// oldFile, isFound := ctrl.Services.Files.CheckFile(c, oid)

		// if isFound {
		// 	if err := ctrl.Services.Files.DeletePreviousFile(oldFile); err != nil {
		// 		return c.String(http.StatusConflict, err.Error())
		// 	}
		// }

		var fileData = models.FileData{
			FileName:  fileName,
			DriveID:   driveID,
			DriveLink: driveLink,
		}

		if err := ctrl.Services.Files.UpdateDownloadLink(c, &fileData, oid); err != nil {
			response.RequestMessage = "Error updating account request"
		}

		response.Success = true
		response.Name = filePath
		response.FileMessage = "File uploaded successfuly"
		response.RequestMessage = "Account request updated successfuly"

		return c.JSON(http.StatusOK, response)

	}

	return c.JSON(http.StatusOK, response)
}

// DownloadInline godoc
// @Summary      Download inline
// @Description  downloads inline
// @Tags         Files
// @Accept       json
// @Produce      json
// @Param        fileName    query     string  false  "file name"
// @Success      200  {string}  string
// @Router       /v2/files/download/inline [get]
func (ctrl FileController) DownloadInline(c echo.Context) error {
	fileName := c.QueryParam("fileName")
	filePath := files.Static + "/" + fileName
	return c.Inline(filePath, fileName)
}

// DownloadAttachment godoc
// @Summary      Download attachment
// @Description  downloads attachment
// @Tags         Files
// @Accept       json
// @Produce      json
// @Param        fileName    query     string  false  "file name"
// @Success      200  {string}  string
// @Router       /v2/files/download/attachment [get]
func (ctrl FileController) DownloadAttachment(c echo.Context) error {
	fileName := c.QueryParam("fileName")
	filePath := files.Static + "/" + fileName
	if err := c.Attachment(filePath, fileName); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.Attachment(filePath, fileName)
}

// DownloadFile godoc
// @Summary      Download file
// @Description  downloads file
// @Tags         Files
// @Accept       json
// @Produce      json
// @Param        fileName    query     string  false  "file name"
// @Success      200  {string}  string
// @Router       /v2/files/download/file [get]
func (ctrl FileController) Download(c echo.Context) error {
	fileName := c.QueryParam("fileName")
	filePath := files.Static + "/" + fileName
	if err := c.File(filePath); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.File(filePath)
}

func (ctrl FileController) DownloadPage(c echo.Context) error {
	return c.File("index.html")
}
