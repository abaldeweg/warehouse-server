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

func List() ([]Product, error) {
	ctx, cancel, collection, err := client()
    defer cancel()
    if err != nil {
		return nil, err
	}

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []Product

	for cur.Next(ctx) {
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

func Create(data Product) (*mongo.InsertOneResult, error) {
	ctx, cancel, collection, err := client()
	defer cancel()
    if err != nil {
		return nil, err
	}

	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func Update(id, key, value string) (*mongo.UpdateResult, error) {
	ctx, cancel, collection, err := client()
	defer cancel()
    if err != nil {
		return nil, err
	}

	d, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{"_id", d}}
	update := bson.D{{"$set", bson.D{{key, value}}}}

	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func Delete(id string) (*mongo.DeleteResult, error) {
	ctx, cancel, collection, err := client()
	defer cancel()
    if err != nil {
		return nil, err
	}

	d, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	document := bson.D{{"_id", d}}

	res, err := collection.DeleteOne(ctx, document)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func client() (context.Context, context.CancelFunc, *mongo.Collection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URL")))
	if err != nil {
        cancel()
		return nil, nil, nil, err
	}

	collection := client.Database("warehouse").Collection("products")

	return ctx, cancel, collection, nil
}
