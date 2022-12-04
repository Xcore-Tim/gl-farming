package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccountType struct {
	ID      string             `json:"_id,omitempty"`
	MongoID primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Name    string             `json:"name" bson:"name"`
}

func (a *AccountType) ConvertID() {
	a.ID = a.MongoID.Hex()
}
