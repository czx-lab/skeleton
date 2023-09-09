package collection

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type CollectionInterface interface {
	SelectPage(ctx context.Context, filter any, sort any, skip, limit int64) (int64, []any, error)
	SelectList(ctx context.Context, filter any, sort interface{}) ([]any, error)
	SelectOne(ctx context.Context, filter any) (interface{}, error)
	SelectCount(ctx context.Context, filter any) (int64, error)
	UpdateOne(ctx context.Context, filter, update any) (int64, error)
	UpdateMany(ctx context.Context, filter, update any) (int64, error)
	Delete(ctx context.Context, filter any) (int64, error)
	InsertOne(ctx context.Context, model any) (any, error)
	InsertMany(ctx context.Context, models []any) ([]any, error)
	Aggregate(ctx context.Context, pipeline any, result any) error
	CreateIndexes(ctx context.Context, indexes []mongo.IndexModel) error
	GetCollection() *mongo.Collection
}
