package db

import (
	"context"

	"github.com/abaldeweg/warehouse-server/logs_web/entity"
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

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	collection := client.Database("logs").Collection("events")

	return &DBHandler{client: client, collection: collection}, nil
}

// Close closes the database connection.
func (handler *DBHandler) Close() error {
	return handler.client.Disconnect(context.TODO())
}

// FindDemanded retrieves logs based on the provided filter.
func (handler *DBHandler) FindDemanded(filter map[string]interface{}) ([]entity.LogEntry, error) {
	bsonFilter := bson.M(filter)
	cursor, err := handler.collection.Find(context.TODO(), bsonFilter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var logs []entity.LogEntry
	for cursor.Next(context.TODO()) {
		var logEntry entity.LogEntry
		if err := cursor.Decode(&logEntry); err != nil {
			return nil, err
		}
		logs = append(logs, logEntry)
	}
	return logs, nil
}
