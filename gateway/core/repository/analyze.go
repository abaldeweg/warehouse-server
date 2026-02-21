package repository

import (
	"context"

	"github.com/abaldeweg/warehouse-server/gateway/core/models"
	"github.com/abaldeweg/warehouse-server/gateway/db/mdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// AnalyzeRepository struct for analyze repository.
type AnalyzeRepository struct {
	client *mongo.Collection
}

// NewAnalyzeRepository creates a new analyze repository.
func NewAnalyzeRepository(db *mdb.MDBClient) *AnalyzeRepository {
	return &AnalyzeRepository{
		client: db.Client.Database("analyze").Collection("shop_search"),
	}
}

// Add stores an AnalyzeData into MongoDB.
func (a *AnalyzeRepository) Add(data models.AnalyzeShopSearch) (*mongo.InsertOneResult, error) {
	return a.client.InsertOne(context.TODO(), data)
}

// FindShopSearchByDateRange returns analyze entries where the date string is between start and end.
func (a *AnalyzeRepository) FindShopSearchByDateRange(start, end string) ([]models.AnalyzeShopSearch, error) {
	filter := bson.M{"date": bson.M{"$gte": start, "$lte": end}}

	cur, err := a.client.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	var results []models.AnalyzeShopSearch
	for cur.Next(context.TODO()) {
		var item models.AnalyzeShopSearch
		if err := cur.Decode(&item); err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
