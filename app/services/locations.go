package services

import (
	"errors"
	"gl-farming/app/models"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LocationService interface {
	Create(echo.Context, models.Location) error
	Update(echo.Context, *models.Location) error
	Delete(echo.Context, string) error
	Get(echo.Context, string) (models.Location, error)
	GetAll(echo.Context, *[]models.Location) error
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

func (s LocationsServiceImpl) Create(c echo.Context, location models.Location) error {

	filter := bson.D{
		bson.E{Key: "$and", Value: bson.A{
			bson.D{
				bson.E{Key: "name", Value: location.Name},
				bson.E{Key: "iso", Value: location.ISO},
			}}},
	}

	if result := s.collection.FindOne(c.Request().Context(), filter); result.Err() == nil {
		return errors.New("such location already exists")
	}

	if _, err := s.collection.InsertOne(c.Request().Context(), location); err != nil {
		return err
	}

	return nil
}

func (s LocationsServiceImpl) Update(c echo.Context, location *models.Location) error {

	filter := bson.D{bson.E{Key: "_id", Value: location.MongoID}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{
		bson.E{Key: "name", Value: location.Name},
		bson.E{Key: "iso", Value: location.ISO},
	}}}
	result, _ := s.collection.UpdateOne(c.Request().Context(), filter, update)

	if result.MatchedCount != 1 {
		return errors.New("no locations found")
	}
	return nil
}

func (s LocationsServiceImpl) Delete(c echo.Context, id string) error {

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	filter := bson.D{bson.E{Key: "_id", Value: oid}}
	result, err := s.collection.DeleteOne(c.Request().Context(), filter)

	if err != nil {
		return err
	}

	if result.DeletedCount != 1 {
		return errors.New("found no location with provided id")
	}

	return nil
}

func (s LocationsServiceImpl) Get(c echo.Context, id string) (models.Location, error) {

	var location models.Location

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return location, err
	}

	filter := bson.D{bson.E{Key: "_id", Value: oid}}
	result := s.collection.FindOne(c.Request().Context(), filter)

	if err := result.Decode(&location); err != nil {
		return location, err
	}

	location.ConvertID()
	return location, err
}

func (s LocationsServiceImpl) GetAll(c echo.Context, locations *[]models.Location) error {

	cursor, err := s.collection.Find(c.Request().Context(), bson.D{{}})

	if err != nil {
		return err
	}

	defer cursor.Close(c.Request().Context())

	for cursor.Next(c.Request().Context()) {
		var location models.Location
		err := cursor.Decode(&location)
		if err != nil {
			return err
		}
		location.ConvertID()
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
