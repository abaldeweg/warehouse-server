package db

import (
	"context"
	"encoding/json"

	"github.com/abaldeweg/warehouse-server/logs_import/entity"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DBHandler handles database operations for logs.
type DBHandler struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewDBHandler creates a new DBHandler.
func NewDBHandler() (*DBHandler, error) {
	clientOptions := options.Client().ApplyURI(viper.Get("MONGODB_URI").(string))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}

	collection := client.Database("logs").Collection("events")

	return &DBHandler{client: client, collection: collection}, nil
}

// Close closes the database connection.
func (handler *DBHandler) Close() error {
	return handler.client.Disconnect(context.TODO())
}

// Write inserts a log entry into the database.
func (handler *DBHandler) Write(date int, data entity.LogEntry) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// Convert JSON data to BSON
	var bsonData bson.M
	if err := bson.UnmarshalExtJSON(jsonData, true, &bsonData); err != nil {
		return err
	}
	_, err = handler.collection.InsertOne(context.TODO(), bsonData)
	return err
}

// Exists checks if a log entry already exists in the database.
func (handler *DBHandler) Exists(date int, data entity.LogEntry) (bool, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return false, err
	}
	var bsonData bson.M
	if err := bson.UnmarshalExtJSON(jsonData, true, &bsonData); err != nil {
		return false, err
	}
	filter := bsonData
	count, err := handler.collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
