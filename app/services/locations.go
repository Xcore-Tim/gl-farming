package services

import (
	"gl-farming/app/models"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LocationService interface {
	Create(echo.Context, models.LocationDTO) error
	Update() error
	Delete() error
	Get() error
	GetAll(echo.Context, *[]models.LocationDTO) error
	DeleteAll(echo.Context) (int, error)
}

type LocationsServiceImpl struct {
	collection *mongo.Collection
}

func NewLocationService(collection *mongo.Collection) LocationService {
	return &LocationsServiceImpl{
		collection: collection,
	}
}

func (s LocationsServiceImpl) Create(c echo.Context, location models.LocationDTO) error {

	if _, err := s.collection.InsertOne(c.Request().Context(), location); err != nil {
		return err
	}

	return nil
}

func (s LocationsServiceImpl) Update() error {
	return nil
}

func (s LocationsServiceImpl) Delete() error {
	return nil
}

func (s LocationsServiceImpl) Get() error {
	return nil
}

func (s LocationsServiceImpl) GetAll(c echo.Context, locations *[]models.LocationDTO) error {

	cursor, err := s.collection.Find(c.Request().Context(), bson.D{{}})

	if err != nil {
		return err
	}

	defer cursor.Close(c.Request().Context())

	for cursor.Next(c.Request().Context()) {
		var location models.LocationDTO
		err := cursor.Decode(&location)
		if err != nil {
			return err
		}
		*locations = append(*locations, location)
	}

	if err := cursor.Err(); err != nil {
		return err
	}

	return nil
}

func (s LocationsServiceImpl) DeleteAll(c echo.Context) (int, error) {

	filter := bson.D{bson.E{}}
	result, err := s.collection.DeleteMany(c.Request().Context(), filter)

	if err != nil {
		return 0, err
	}

	return int(result.DeletedCount), nil
}
