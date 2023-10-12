package client

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Product struct {
	Id         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name" bson:"name"`
	Attributes []interface{}      `json:"attributes" bson:"attributes"`
	Variants   []interface{}      `json:"variants" bson:"variants"`
}

type ClientConfig struct {
	Ctx        context.Context
	Cancel     context.CancelFunc
	Client     *mongo.Client
	Collection *mongo.Collection
}

func NewClient() (*ClientConfig, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URL")))
	if err != nil {
		cancel()
		return nil, err
	}

	collection := client.Database("warehouse").Collection("products")

	return &ClientConfig{ctx, cancel, client, collection}, nil
}

func (config *ClientConfig) List() ([]Product, error) {
	if config.Cancel != nil {
		defer config.Cancel()
	}

	cur, err := config.Collection.Find(config.Ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(config.Ctx)

	var results []Product

	for cur.Next(config.Ctx) {
		var result Product
		if err := cur.Decode(&result); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (config *ClientConfig) Create(data Product) (*mongo.InsertOneResult, error) {
	if config.Cancel != nil {
		defer config.Cancel()
	}

	res, err := config.Collection.InsertOne(config.Ctx, data)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (config *ClientConfig) Update(id, key, value string) (*mongo.UpdateResult, error) {
	if config.Cancel != nil {
		defer config.Cancel()
	}

	d, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{"_id", d}}
	update := bson.D{{"$set", bson.D{{key, value}}}}

	res, err := config.Collection.UpdateOne(config.Ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (config *ClientConfig) Delete(id string) (*mongo.DeleteResult, error) {
	if config.Cancel != nil {
		defer config.Cancel()
	}

	d, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	document := bson.D{{"_id", d}}

	res, err := config.Collection.DeleteOne(config.Ctx, document)
	if err != nil {
		return nil, err
	}

	return res, nil
}
