package analyze

import (
	"context"

	"github.com/abaldeweg/warehouse-server/gateway/db/mdb"
	"go.mongodb.org/mongo-driver/mongo"
)

type AnalyzeClient struct {
	Collection *mongo.Collection
}

type AnalyzeShopSearch struct {
	Branch int
	Date   string
	Term   string
	Page   int
	Genre  int
}

const DBName string = "analyze"

// NewAnalyze creates a new AnalyzeClient with the given MongoDB client and collection name.
func NewAnalyze(db *mdb.MDBClient, col string) *AnalyzeClient {
	collection := db.Client.Database(DBName).Collection(col)
	return &AnalyzeClient{
		Collection: collection,
	}
}

// Add stores an AnalyzeData into MongoDB.
func (a AnalyzeClient) Add(data AnalyzeShopSearch) (*mongo.InsertOneResult, error) {
	result, err := a.Collection.InsertOne(context.TODO(), data)
	return result, err
}
