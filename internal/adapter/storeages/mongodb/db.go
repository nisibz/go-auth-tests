package mongodb

import (
	"context"
	"time"

	"github.com/nisibz/go-auth-tests/internal/adapter/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// NewMongoClient creates a new MongoDB client and pings the database
func NewMongoClient(config *config.Mongo) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	client, err := mongo.Connect(options.Client().ApplyURI(config.URI))
	if err != nil {
		return nil, err
	}

	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}
