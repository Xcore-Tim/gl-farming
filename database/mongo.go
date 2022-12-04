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

func Init() (*Collections, error) {

	client, err := Connect()

	if err != nil {
		return nil, err
	}

	var Collections = Collections{
		AccountRequests: client.Database("gypsyland").Collection("accountRequests"),
		AccountTypes:    client.Database("gypsyland").Collection("accountTypes"),
		Locations:       client.Database("gypsyland").Collection("locations"),
		Currency:        client.Database("gypsyland").Collection("currency"),
		FarmerAccess:    client.Database("gypsyland").Collection("farmerAccess"),
	}

	return &Collections, nil
}

func Connect() (*mongo.Client, error) {

	ctx := context.TODO()
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
