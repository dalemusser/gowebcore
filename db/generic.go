package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// FindOne decodes a single document into v.
func FindOne[T any](ctx context.Context, coll *mongo.Collection, filter any) (*T, error) {
	var out T
	if err := coll.FindOne(ctx, filter).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

// InsertMany inserts docs and returns inserted IDs.
func InsertMany[T any](ctx context.Context, coll *mongo.Collection, docs []T) ([]interface{}, error) {
	var anyDocs []interface{}
	for _, d := range docs {
		anyDocs = append(anyDocs, d)
	}
	res, err := coll.InsertMany(ctx, anyDocs)
	if err != nil {
		return nil, err
	}
	return res.InsertedIDs, nil
}

// UpdateByID performs {$set: updates} on the given ID.
func UpdateByID(ctx context.Context, coll *mongo.Collection, id any, updates any) error {
	_, err := coll.UpdateByID(ctx, id, bson.M{"$set": updates})
	return err
}
