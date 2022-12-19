package services

import (
	"errors"
	"gl-farming/app/models"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TableService interface {
	Get(echo.Context, *models.TableDataRequest) error
	GetAll(echo.Context, *models.TableDataRequest) error
	AggregateDataByUID(echo.Context, primitive.D, primitive.D) ([]models.EmployeePipeline, error)
	AggregateDataByTeam(echo.Context, primitive.D, primitive.D) ([]models.TeamPipiline, error)
}

type TableServiceServiceImpl struct {
	collection *mongo.Collection
}

func NewTableService(collection *mongo.Collection) TableService {
	return &TableServiceServiceImpl{
		collection: collection,
	}
}

func (s TableServiceServiceImpl) Get(c echo.Context, tableDataRequest *models.TableDataRequest) error {

	cursor, err := s.collection.Find(c.Request().Context(), tableDataRequest.Filter, options.Find().SetProjection(tableDataRequest.Projection))
	if err != nil {
		return err
	}

	for cursor.Next(c.Request().Context()) {
		var accountRequest models.TableData
		if err := cursor.Decode(&accountRequest); err != nil {
			return err
		}
		tableDataRequest.DataSlice = append(tableDataRequest.DataSlice, accountRequest)
	}

	if err := cursor.Err(); err != nil {
		return err
	}

	cursor.Close(c.Request().Context())

	if len(tableDataRequest.DataSlice) == 0 {
		return errors.New("documents not found")
	}

	return err
}

func (s TableServiceServiceImpl) GetAll(c echo.Context, tableDataRequest *models.TableDataRequest) error {

	cursor, err := s.collection.Find(c.Request().Context(), tableDataRequest.Filter)
	if err != nil {
		return err
	}

	for cursor.Next(c.Request().Context()) {
		var accountRequest models.TableData

		if err := cursor.Decode(&accountRequest); err != nil {
			return err
		}

		tableDataRequest.DataSlice = append(tableDataRequest.DataSlice, accountRequest)
	}

	if err := cursor.Err(); err != nil {
		return err
	}

	return nil
}

func (s TableServiceServiceImpl) AggregateDataByUID(c echo.Context, matchStage primitive.D, groupStage primitive.D) ([]models.EmployeePipeline, error) {

	var pipelineResult []models.EmployeePipeline

	cursor, err := s.collection.Aggregate(c.Request().Context(), mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		return pipelineResult, err
	}

	for cursor.Next(c.Request().Context()) {
		var tableData models.EmployeePipeline
		if err := cursor.Decode(&tableData); err != nil {
			return pipelineResult, err
		}

		pipelineResult = append(pipelineResult, tableData)
	}

	if cursor.Err() != nil {
		return pipelineResult, err
	}

	return pipelineResult, err
}

func (s TableServiceServiceImpl) AggregateDataByTeam(c echo.Context, matchStage primitive.D, groupStage primitive.D) ([]models.TeamPipiline, error) {

	var pipelineResult []models.TeamPipiline

	cursor, err := s.collection.Aggregate(c.Request().Context(), mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		return pipelineResult, err
	}

	for cursor.Next(c.Request().Context()) {
		var tableData models.TeamPipiline
		if err := cursor.Decode(&tableData); err != nil {
			return pipelineResult, err
		}

		pipelineResult = append(pipelineResult, tableData)
	}

	if cursor.Err() != nil {
		return pipelineResult, err
	}

	return pipelineResult, err
}
