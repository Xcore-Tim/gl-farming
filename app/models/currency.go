package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Currency struct {
	ID      string             `json:"_id,omitempty"`
	MongoID primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Name    string             `json:"name" bson:"name"`
	ISO     string             `json:"iso" bson:"iso"`
	Symbol  string             `json:"symbol" bson:"symbol"`
}

func (c *Currency) ConvertID() {
	c.ID = c.MongoID.Hex()
}
