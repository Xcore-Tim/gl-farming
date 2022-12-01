package services

import (
	"errors"
	"gl-farming/app/models"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CurrencyService interface {
	Create(echo.Context, models.Currency) error
	Update() error
	Delete(echo.Context, string) error
	Get(echo.Context, string) (models.Currency, error)
	GetAll(echo.Context, *[]models.Currency) error
	DeleteAll(echo.Context) (int, error)
}

type CurrencyServiceImpl struct {
	collection *mongo.Collection
}

func NewCurrencyService(collection *mongo.Collection) CurrencyService {
	return &CurrencyServiceImpl{
		collection: collection,
	}
}

func (s CurrencyServiceImpl) Create(c echo.Context, currency models.Currency) error {

	filter := bson.D{
		bson.E{Key: "$and", Value: bson.A{
			bson.D{
				bson.E{Key: "name", Value: currency.Name},
				bson.E{Key: "iso", Value: currency.ISO},
				bson.E{Key: "symbol", Value: currency.Symbol},
			}}},
	}

	if result := s.collection.FindOne(c.Request().Context(), filter); result.Err() == nil {
		return errors.New("such account type already exists")
	}

	if _, err := s.collection.InsertOne(c.Request().Context(), currency); err != nil {
		return err
	}

	return nil
}

func (s CurrencyServiceImpl) Update() error {
	return nil
}

func (s CurrencyServiceImpl) Delete(c echo.Context, id string) error {

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

func (s CurrencyServiceImpl) Get(c echo.Context, id string) (models.Currency, error) {

	var currency models.Currency

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return currency, err
	}

	filter := bson.D{bson.E{Key: "_id", Value: oid}}
	result := s.collection.FindOne(c.Request().Context(), filter)

	if err := result.Decode(&currency); err != nil {
		return currency, err
	}

	currency.ConvertID()
	return currency, err
}

func (s CurrencyServiceImpl) GetAll(c echo.Context, currencyList *[]models.Currency) error {

	cursor, err := s.collection.Find(c.Request().Context(), bson.D{{}})

	if err != nil {
		return err
	}

	defer cursor.Close(c.Request().Context())

	for cursor.Next(c.Request().Context()) {
		var currency models.Currency
		err := cursor.Decode(&currency)
		if err != nil {
			return err
		}
		currency.ConvertID()
		*currencyList = append(*currencyList, currency)
	}

	if err := cursor.Err(); err != nil {
		return err
	}

	return nil
}

func (s CurrencyServiceImpl) DeleteAll(c echo.Context) (int, error) {

	filter := bson.D{bson.E{}}
	result, err := s.collection.DeleteMany(c.Request().Context(), filter)

	if err != nil {
		return 0, err
	}

	return int(result.DeletedCount), nil
}
