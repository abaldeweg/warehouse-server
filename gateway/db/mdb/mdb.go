package mdb

import (
	"context"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MDBClient handles MongoDB connections.
type MDBClient struct {
	Client *mongo.Client
}

// NewMDBClient creates a new MongoDB client.
func NewMDBClient() (*MDBClient, error) {
	viper.SetDefault("MONGODB_URI", "mongodb://localhost:27017")
	clientOptions := options.Client().ApplyURI(viper.Get("MONGODB_URI").(string))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	return &MDBClient{Client: client}, nil
}
