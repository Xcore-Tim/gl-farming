package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Collections struct {
	AccountRequests *mongo.Collection
	AccountTypes    *mongo.Collection
	Currency        *mongo.Collection
	Locations       *mongo.Collection
	FarmerAccess    *mongo.Collection
}

func Init(ctx context.Context) (*mongo.Client, error) {

	connectionAddress := "mongodb://localhost:27017"
	mongoConnection := options.Client().ApplyURI(connectionAddress)
	mongoClient, err := mongo.Connect(ctx, mongoConnection)

	if err != nil {
		return nil, err
	}

	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	fmt.Println("Connection established")
	return mongoClient, nil

}
