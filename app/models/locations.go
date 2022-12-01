package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Location struct {
	ID      string             `json:"_id,omitempty"`
	MongoID primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Name    string             `json:"name" bson:"name,omitempty"`
	ISO     string             `json:"iso" bson:"iso,omitempty"`
}

func (l *Location) ConvertID() {
	l.ID = l.MongoID.Hex()
}
