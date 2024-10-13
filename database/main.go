package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	db         *mongo.Database
	collection *mongo.Collection
)

func StartDatabase() error {
	if err := connectDatabase(); err != nil {
		return err
	}
	return nil
}

func connectDatabase() error {
	// connection string for the mongodb container
	mongoURI := "mongodb://admin:password@localhost:27017/?authSource=admin"

	// Set client options
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB
	client, err := mongo.Connect(clientOptions)
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

func CloseDatabase() {
	// Close the connection with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10)
	defer cancel()
	if err := db.Client().Disconnect(ctx); err != nil {
		log.Fatalf("Failed to close database: %v\n", err)
	}
}
