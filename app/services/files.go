package services

import (
	"errors"
	"gl-farming/app/constants/files"
	"gl-farming/app/models"
	"mime/multipart"
	"os"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FileService interface {
	CreateNewFile(*multipart.FileHeader, string) (fileName string, filePath string, err error)
	CheckPreviousFile(echo.Context, string) (oldFile string, isFound bool)
	UpdateDownloadLink(echo.Context, string, string) error
	DeletePreviousFile(string) error
}

type FileServiceImpl struct {
	collection *mongo.Collection
}

func NewFiletService(collection *mongo.Collection) FileService {
	return &FileServiceImpl{
		collection: collection,
	}
}

func (s FileServiceImpl) CreateNewFile(uploadFile *multipart.FileHeader, fileExt string) (fileName string, filePath string, err error) {

	constructedFile := files.Static + "/" + uploadFile.Filename

	finalFile, err := os.Open(constructedFile)

	if err != nil {
		err = errors.New("failed to upload file")
		return
	}

	finalFile.Close()

	fileName = uuid.NewString() + fileExt
	filePath = files.Static + "/" + fileName

	if err = os.Rename(constructedFile, filePath); err != nil {
		filePath = constructedFile
	}

	return fileName, filePath, err
}

func (s FileServiceImpl) CheckPreviousFile(c echo.Context, oid string) (oldFile string, isFound bool) {

	requestID, err := primitive.ObjectIDFromHex(oid)

	if err != nil {
		return
	}

	var accountRequest models.AccountRequest

	filter := bson.D{bson.E{Key: "_id", Value: requestID}}

	if err = s.collection.FindOne(c.Request().Context(), filter).Decode(&accountRequest); err != nil {
		return
	}

	if accountRequest.FileName != "" {
		oldFile = files.Static + "/" + accountRequest.FileName
		isFound = true
		return
	}

	return
}

func (s FileServiceImpl) UpdateDownloadLink(c echo.Context, fileName, oid string) error {

	requestID, err := primitive.ObjectIDFromHex(oid)

	if err != nil {
		return err
	}

	filter := bson.D{bson.E{Key: "_id", Value: requestID}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{
		bson.E{Key: "fileName", Value: fileName},
	}}}

	result := s.collection.FindOneAndUpdate(c.Request().Context(), filter, update)

	if result.Err() != nil {
		return errors.New("no matched documents found for update")
	}

	return nil
}

func (s FileServiceImpl) DeletePreviousFile(filePath string) error {
	err := os.Remove(filePath)
	return err
}
