package db

import (
	"context"
	"time"

	"github.com/abaldeweg/warehouse-server/logs_web/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DBHandler handles database operations for logs.
type DBHandler struct {
	client *mongo.Client
	collection *mongo.Collection
}

// NewDBHandler creates a new DBHandler.
func NewDBHandler() (*DBHandler, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
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

// Get fetches all log entries from the database.
func (handler *DBHandler) Get(from int, to int) ([]entity.LogEntry, error) {
	filter := bson.M{"date": bson.M{"$gte": time.Unix(int64(from), 0), "$lte": time.Unix(int64(to), 0)}}
	cursor, err := handler.collection.Find(context.TODO(), filter)
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
