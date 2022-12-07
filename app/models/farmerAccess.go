package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type FarmerAccess struct {
	ID      string             `json:"_id,omitempty"`
	MongoID primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Farmer  Employee           `json:"farmer" bson:"farmer"`
	Team    int                `json:"team" bson:"team"`
}

type FarmerAccessList struct {
	Farmer Employee `json:"farmer" bson:"farmer"`
	Teams  []int    `json:"teams" bson:"teams"`
}

type AccessRequest struct {
	Farmer Employee `json:"farmer" bson:"farmer"`
	TeamID int      `json:"teamID" bson:"team"`
}

type FullAccessRequest struct {
	Farmer Employee `json:"farmer"`
}

func (f *FarmerAccess) ConvertID() {
	f.ID = f.MongoID.Hex()
}
