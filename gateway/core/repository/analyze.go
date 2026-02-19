package repository

import (
	"context"

	"github.com/abaldeweg/warehouse-server/gateway/core/models"
	"github.com/abaldeweg/warehouse-server/gateway/db/mdb"
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
