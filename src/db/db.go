package db

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/rohan-luthra/microservice-grpc-go/service-restaurants-go/src/kafka"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Database base type struct
type Database struct {
	Client                    *mongo.Client
	RestaurantCollection      *mongo.Collection
	TaxCollection             *mongo.Collection
	CurrencyCollection        *mongo.Collection
	RestaurantGroupCollection *mongo.Collection
	Kafka                     *kafka.Kafka
}

// New mongodb database instance
func New(config *Config) (error, *Database) {

	// Set client options
	clientOptions := options.Client().ApplyURI(config.DatabaseURI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return errors.Wrap(err, "unable to connect to database"), nil
	}

	// Check the connection
	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		return errors.Wrap(err, "unable to connect to database"), nil
	}

	k, err := kafka.New()
	if err != nil {
		return err, nil
	}

	restaurantCollection := client.Database(config.Database).Collection("restaurants")
	taxCollection := client.Database(config.Database).Collection("taxes")
	currencyCollection := client.Database(config.Database).Collection("currencies")
	restaurantGroupCollection := client.Database(config.Database).Collection("restaurant_group")

	return nil, &Database{
		Client:                    client,
		RestaurantCollection:      restaurantCollection,
		TaxCollection:             taxCollection,
		CurrencyCollection:        currencyCollection,
		RestaurantGroupCollection: restaurantGroupCollection,
		Kafka:                     k,
	}
}
