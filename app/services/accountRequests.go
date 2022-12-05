package services

import (
	"errors"
	"gl-farming/app/constants/requestStatus"
	"gl-farming/app/models"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccountRequestService interface {
	Cancel(echo.Context, *models.CancelAccountRequest) error
	Create(echo.Context, *models.AccountRequest) error
	Take(echo.Context, *models.TakeAccountRequest) error
	Update(echo.Context, *models.AccountRequest) error
	Complete(echo.Context, *models.AccountRequest) error
	Return(echo.Context, *models.ReturnAccountRequest) error
	Get(echo.Context, string) (models.AccountRequest, error)

	DeleteAll(echo.Context) (int, error)
}

type AccountRequestServiceImpl struct {
	collection *mongo.Collection
}

func NewAccountRequestService(collection *mongo.Collection) AccountRequestService {
	return &AccountRequestServiceImpl{
		collection: collection,
	}
}

func (s AccountRequestServiceImpl) Create(c echo.Context, accountRequest *models.AccountRequest) error {

	if _, err := s.collection.InsertOne(c.Request().Context(), accountRequest); err != nil {
		return err
	}

	return nil
}

func (s AccountRequestServiceImpl) Cancel(c echo.Context, cancellRequest *models.CancelAccountRequest) error {

	oid, err := primitive.ObjectIDFromHex(cancellRequest.RequestID)
	if err != nil {
		return err
	}

	filter := bson.D{bson.E{Key: "_id", Value: oid}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{
		bson.E{Key: "cancelledBy", Value: cancellRequest.CancelledBy},
		bson.E{Key: "cancellationCause", Value: cancellRequest.CancellationCause},
		bson.E{Key: "dateCancelled", Value: cancellRequest.DateCancelled},
		bson.E{Key: "status", Value: requestStatus.Cancelled},
	}}}

	result, _ := s.collection.UpdateOne(c.Request().Context(), filter, update)

	if result.MatchedCount != 1 {
		return errors.New("no requests found")
	}

	return nil
}

func (s AccountRequestServiceImpl) Take(c echo.Context, takeRequest *models.TakeAccountRequest) error {

	oid, err := primitive.ObjectIDFromHex(takeRequest.RequestID)
	if err != nil {
		return err
	}

	filter := bson.D{bson.E{Key: "_id", Value: oid}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{
		bson.E{Key: "farmer", Value: takeRequest.Farmer},
		bson.E{Key: "status", Value: requestStatus.Inwork},
		bson.E{Key: "takenBy", Value: takeRequest.Farmer},
		bson.E{Key: "dateTaken", Value: time.Now().Unix()},
	}}}

	result, _ := s.collection.UpdateOne(c.Request().Context(), filter, update)

	if result.MatchedCount != 1 {
		return errors.New("no requests found")
	}

	return nil
}

func (s AccountRequestServiceImpl) Update(c echo.Context, accountRequest *models.AccountRequest) error {

	filter := bson.D{bson.E{Key: "_id", Value: accountRequest.MongoID}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{
		bson.E{Key: "location", Value: accountRequest.Location},
		bson.E{Key: "type", Value: accountRequest.Type},
		bson.E{Key: "quantity", Value: accountRequest.Quantity},
		bson.E{Key: "currency", Value: accountRequest.Currency},
		bson.E{Key: "rate", Value: accountRequest.Rate},
		bson.E{Key: "price", Value: accountRequest.Price},
		bson.E{Key: "total", Value: accountRequest.Total},
		bson.E{Key: "crossRate", Value: accountRequest.CrossRate},
		bson.E{Key: "baseCurrency", Value: accountRequest.BaseCurrency},
		bson.E{Key: "baseRate", Value: accountRequest.BaseRate},
		bson.E{Key: "basePrice", Value: accountRequest.Price},
		bson.E{Key: "baseTotal", Value: accountRequest.Total},
		bson.E{Key: "updateBy", Value: accountRequest.UpdatedBy},
		bson.E{Key: "dateUpdated", Value: accountRequest.DateUpdated},
	}}}

	result, _ := s.collection.UpdateOne(c.Request().Context(), filter, update)

	if result.ModifiedCount != 1 {
		return errors.New("no requests found")
	}

	return nil
}

func (s AccountRequestServiceImpl) Complete(c echo.Context, accountRequest *models.AccountRequest) error {

	filter := bson.D{bson.E{Key: "_id", Value: accountRequest.MongoID}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{
		bson.E{Key: "status", Value: requestStatus.Complete},
		bson.E{Key: "valid", Value: accountRequest.Valid},
		bson.E{Key: "currency", Value: accountRequest.Currency},
		bson.E{Key: "rate", Value: accountRequest.Rate},
		bson.E{Key: "price", Value: accountRequest.Price},
		bson.E{Key: "total", Value: accountRequest.Total},
		bson.E{Key: "crossRate", Value: accountRequest.CrossRate},
		bson.E{Key: "baseCurrency", Value: accountRequest.BaseCurrency},
		bson.E{Key: "baseRate", Value: accountRequest.BaseRate},
		bson.E{Key: "basePrice", Value: accountRequest.Price},
		bson.E{Key: "baseTotal", Value: accountRequest.Total},
		bson.E{Key: "completedBy", Value: accountRequest.CompletedBy},
		bson.E{Key: "dateCompleted", Value: accountRequest.DateCompleted},
	}}}

	result, _ := s.collection.UpdateOne(c.Request().Context(), filter, update)

	if result.ModifiedCount != 1 {
		return errors.New("no requests found")
	}

	return nil
}

func (s AccountRequestServiceImpl) Return(c echo.Context, returnAccountRequest *models.ReturnAccountRequest) error {

	oid, _ := primitive.ObjectIDFromHex(returnAccountRequest.RequestID)
	filter := bson.D{bson.E{Key: "_id", Value: oid}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{
		bson.E{Key: "status", Value: requestStatus.Pending},
		bson.E{Key: "farmer", Value: returnAccountRequest.Farmer},
		bson.E{Key: "dateReturned", Value: time.Now().Unix()},
		bson.E{Key: "returneddBy", Value: returnAccountRequest.ReturnedBy},
	}}}

	result := s.collection.FindOneAndUpdate(c.Request().Context(), filter, update)

	if result.Err() != nil {
		return errors.New("no requests found")
	}
	return nil
}

func (s AccountRequestServiceImpl) Get(c echo.Context, id string) (models.AccountRequest, error) {

	var accountRequest models.AccountRequest

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return accountRequest, err
	}

	filter := bson.D{bson.E{Key: "_id", Value: oid}}
	if err := s.collection.FindOne(c.Request().Context(), filter).Decode(&accountRequest); err != nil {
		return accountRequest, err
	}

	return accountRequest, nil
}

func (s AccountRequestServiceImpl) DeleteAll(c echo.Context) (int, error) {
	filter := bson.D{bson.E{}}
	result, err := s.collection.DeleteMany(c.Request().Context(), filter)

	if err != nil {
		return 0, err
	}

	return int(result.DeletedCount), nil
}
