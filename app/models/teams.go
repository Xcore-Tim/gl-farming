package models

type Team struct {
	ID       int      `json:"id" bson:"id"`
	Teamlead Employee `json:"teamlead" bson:"teamlead"`
}
