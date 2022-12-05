package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccountRequest struct {
	ID       string             `json:"_id,omitempty"`
	MongoID  primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Type     AccountType        `json:"type" bson:"type"`
	Location Location           `json:"location" bson:"location"`
	Status   int                `json:"status" bson:"status"`
	Team     Team               `json:"team" bson:"team"`

	Quantity uint `json:"quantity" bson:"quantity"`
	Valid    uint `json:"valid" bson:"valid"`

	Currency Currency `json:"currency" bson:"currency"`
	Rate     float64  `json:"rate" bson:"rate"`
	Price    float64  `json:"price" bson:"price"`
	Total    float64  `json:"total" bson:"total"`

	CrossRate float64 `json:"crossRate" bson:"crossRate"`

	BaseCurrency Currency `json:"baseCurrency" bson:"baseCurrency"`
	BaseRate     float64  `json:"baseRate" bson:"baseRate"`
	BasePrice    float64  `json:"basePrice" bson:"basePrice"`
	BaseTotal    float64  `json:"baseTotal" bson:"baseTotal"`

	Buyer  Employee `json:"buyer" bson:"buyer"`
	Farmer Employee `json:"farmer" bson:"farmer"`

	TakenBy     Employee `json:"takenBy" bson:"takenBy"`
	UpdatedBy   Employee `json:"updatedBy" bson:"updatedBy"`
	CancelledBy Employee `json:"cancelledBy" bson:"cancelledBy"`
	CompletedBy Employee `json:"completedBy" bson:"completedBy"`
	ReturnedBy  Employee `json:"returnedBy" bson:"returnedBy"`

	Description       string `json:"description" bson:"description"`
	CancellationCause string `json:"cancellationCause" bson:"cancellationCause"`

	FileName string `json:"fileName" bson:"fileName"`

	DateCreated   int64 `json:"dateCreated" bson:"dateCreated"`
	DateTaken     int64 `json:"dateTaken" bson:"dateTaken"`
	DateUpdated   int64 `json:"dateUpdated" bson:"dateUpdated"`
	DateCancelled int64 `json:"dateCancelled" bson:"dateCancelled"`
	DateCompleted int64 `json:"dateCompleted" bson:"dateCompleted"`
	DateReturned  int64 `json:"dateReturned" bson:"dateReturned"`
}

type CreateAccountRequest struct {
	TypeID      string  `json:"typeID"`
	LocationID  string  `json:"locationID"`
	CurrencyID  string  `json:"currencyID"`
	Quantity    uint    `json:"quantity"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}

type CancelAccountRequest struct {
	RequestID         string   `json:"requestID"`
	CancelledBy       Employee `json:"-"`
	CancellationCause string   `json:"cancellationCause"`
	DateCancelled     int64    `json:"-"`
}

type CompleteAccountRequest struct {
	RequestID  string  `json:"requestID"`
	CurrencyID string  `json:"currencyID"`
	Price      float64 `json:"price"`
	Valid      uint    `json:"valid"`
}

type TakeAccountRequest struct {
	RequestID string   `json:"requestID"`
	Farmer    Employee `json:"-"`
}

type ReturnAccountRequest struct {
	RequestID  string   `json:"requestID"`
	Farmer     Employee `json:"_"`
	ReturnedBy Employee `json:"-"`
}

type UpdateAccountRequest struct {
	RequestID   string  `json:"requestID"`
	TypeID      string  `json:"typeID"`
	Description string  `json:"description"`
	CurrencyID  string  `json:"currencyID"`
	LocationID  string  `json:"locationID"`
	Price       float64 `json:"price"`
	Quantity    uint    `json:"quantity"`
}

func (a *AccountRequest) ConvertID() {
	a.ID = a.MongoID.Hex()
}
