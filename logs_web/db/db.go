package db

import (
	"context"
	"fmt"
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

// Get fetches log entries from the database within a specified time range.
func (handler *DBHandler) Get(from int, to int) ([]entity.LogEntry, error) {
	fromTime, err := time.Parse("20060102", fmt.Sprintf("%d", from))
	if err != nil {
		return nil, err
	}
	toTime, err := time.Parse("20060102", fmt.Sprintf("%d", to))
	if err != nil {
		return nil, err
	}
	toTime = toTime.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	filter := bson.M{
		"time": bson.M{
			"$gte": fromTime.Format(time.RFC3339),
			"$lte": toTime.Format(time.RFC3339),
		},
	}
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
