package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	db         *mongo.Database
	collection *mongo.Collection
)

func StartDabase() error {
	if err := connectDatabase(); err != nil {
		return err
	}
	return nil
}

func connectDatabase() error {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}

	// Create or get the database
	db = client.Database("newsDB")

	// Create or get the collection
	collection = db.Collection("articles")

	return nil
}
