package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Location struct {
	Name string `json:"name" bson:"name"`
	ISO  string `json:"iso" bson:"iso"`
}

type LocationDTO struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name"`
	ISO  string             `json:"iso" bson:"iso"`
}

func (l Location) ToDTO() LocationDTO {
	return LocationDTO{
		Name: l.Name,
		ISO:  l.ISO,
	}
}
