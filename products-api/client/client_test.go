package client

import (
	"context"
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestList(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Test List",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		mt.Run(tt.name, func(mt *mtest.T) {
			res1 := mtest.CreateCursorResponse(1, "test.documents", mtest.FirstBatch, bson.D{})
			res2 := mtest.CreateCursorResponse(1, "test.documents", mtest.NextBatch, bson.D{})
			res3 := mtest.CreateCursorResponse(0, "test.documents", mtest.NextBatch)
			mt.AddMockResponses(res1, res2, res3)

            config := &ClientConfig{context.TODO(),nil,nil,mt.Coll}

			got, err := config.List()
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(got).String() == "*mongo.Collection" {
				t.Errorf("List() = %v, want %v", reflect.TypeOf(got).String(), "*mongo.Collection")
			}
		})
	}
}

func TestCreate(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	type args struct {
		data Product
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Test Create",
			args:    args{Product{}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		mt.Run(tt.name, func(mt *mtest.T) {
			mt.AddMockResponses(mtest.CreateSuccessResponse())

            config := &ClientConfig{context.TODO(),nil,nil,mt.Coll}

			got, err := config.Create(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(got).String() != "*mongo.InsertOneResult" {
				t.Errorf("Create() = %v, want %v", reflect.TypeOf(got).String(), "*mongo.InsertOneResult")
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	type args struct {
		id    string
		key   string
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    *mongo.UpdateResult
		wantErr bool
	}{
		{
			name:    "Test Update",
			args:    args{"65083e341084d75d12a1f969", "name", "Test"},
			want:    &mongo.UpdateResult{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		mt.Run(tt.name, func(mt *mtest.T) {
			mt.AddMockResponses(mtest.CreateSuccessResponse())

            config := &ClientConfig{context.TODO(),nil,nil,mt.Coll}

			got, err := config.Update(tt.args.id, tt.args.key, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    *mongo.DeleteResult
		wantErr bool
	}{
		{
			name:    "Test Delete",
			args:    args{"65083e341084d75d12a1f969"},
			want:    &mongo.DeleteResult{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		mt.Run(tt.name, func(mt *mtest.T) {
			mt.AddMockResponses(mtest.CreateSuccessResponse())

            config := &ClientConfig{context.TODO(),nil,nil, mt.Coll}

			got, err := config.Delete(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}
