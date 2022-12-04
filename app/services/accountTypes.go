package services

import (
	"errors"

	"gl-farming/app/models"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccountTypeService interface {
	Create(echo.Context, models.AccountType) error
	Update() error
	Delete(echo.Context, string) error
	Get(echo.Context, string) (models.AccountType, error)
	GetAll(echo.Context, *[]models.AccountType) error
	DeleteAll(echo.Context) (int, error)
}

type AccountTypeServiceImpl struct {
	collection *mongo.Collection
}

func NewAccountTypeService(collection *mongo.Collection) AccountTypeService {
	return &AccountTypeServiceImpl{
		collection: collection,
	}
}

func (s AccountTypeServiceImpl) Create(c echo.Context, accountType models.AccountType) error {

	filter := bson.D{bson.E{Key: "name", Value: accountType.Name}}

	if result := s.collection.FindOne(c.Request().Context(), filter); result.Err() == nil {
		return errors.New("such account type already exists")
	}

	if _, err := s.collection.InsertOne(c.Request().Context(), accountType); err != nil {
		return err
	}

	return nil
}

func (s AccountTypeServiceImpl) Update() error {
	return nil
}

func (s AccountTypeServiceImpl) Delete(c echo.Context, id string) error {

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
		return errors.New("found no account type with provided id")
	}

	return nil
}

func (s AccountTypeServiceImpl) Get(c echo.Context, id string) (models.AccountType, error) {

	var accountType models.AccountType

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return accountType, err
	}

	filter := bson.D{bson.E{Key: "_id", Value: oid}}
	result := s.collection.FindOne(c.Request().Context(), filter)

	if err := result.Decode(&accountType); err != nil {
		return accountType, err
	}

	accountType.ConvertID()
	return accountType, err
}

func (s AccountTypeServiceImpl) GetAll(c echo.Context, accountTypes *[]models.AccountType) error {

	cursor, err := s.collection.Find(c.Request().Context(), bson.D{{}})

	if err != nil {
		return err
	}

	defer cursor.Close(c.Request().Context())

	for cursor.Next(c.Request().Context()) {
		var accountType models.AccountType
		if err := cursor.Decode(&accountType); err != nil {
			return err
		}
		accountType.ConvertID()
		*accountTypes = append(*accountTypes, accountType)
	}

	if err := cursor.Err(); err != nil {
		return err
	}

	return nil
}

func (s AccountTypeServiceImpl) DeleteAll(c echo.Context) (int, error) {

	filter := bson.D{bson.E{}}
	result, err := s.collection.DeleteMany(c.Request().Context(), filter)

	if err != nil {
		return 0, err
	}

	return int(result.DeletedCount), nil
}
