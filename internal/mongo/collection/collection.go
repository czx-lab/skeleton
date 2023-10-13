package collection

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection struct {
	DbName     string
	ColName    string
	DataBase   *mongo.Database
	Collection *mongo.Collection
}

func (c *Collection) SelectPage(ctx context.Context, filter any, sort any, skip, limit int64) (int64, []any, error) {
	var err error

	resultCount, err := c.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	opts := options.Find().SetSort(sort).SetSkip(skip).SetLimit(limit)
	finder, err := c.Collection.Find(ctx, filter, opts)
	if err != nil {
		return resultCount, nil, err
	}

	result := make([]interface{}, 0)
	if err := finder.All(ctx, &result); err != nil {
		return resultCount, nil, err
	}
	return resultCount, result, nil
}

func (c *Collection) SelectList(ctx context.Context, filter any, sort any) ([]any, error) {
	var err error

	opts := options.Find().SetSort(sort)
	finder, err := c.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, 0)
	if err := finder.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, err
}

func (c *Collection) SelectOne(ctx context.Context, filter any) (any, error) {
	result := new(interface{})
	err := c.Collection.FindOne(ctx, filter, options.FindOne()).Decode(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Collection) SelectCount(ctx context.Context, filter any) (int64, error) {
	return c.Collection.CountDocuments(ctx, filter)
}

func (c *Collection) UpdateOne(ctx context.Context, filter, update any) (int64, error) {
	result, err := c.Collection.UpdateOne(ctx, filter, update, options.Update())
	if err != nil {
		return 0, err
	}
	if result.MatchedCount == 0 {
		return 0, fmt.Errorf("Update result: %s ", "document not found")
	}
	return result.MatchedCount, nil
}

func (c *Collection) UpdateMany(ctx context.Context, filter, update any) (int64, error) {
	result, err := c.Collection.UpdateMany(ctx, filter, update, options.Update())
	if err != nil {
		return 0, err
	}
	if result.MatchedCount == 0 {
		return 0, fmt.Errorf("Update result: %s ", "document not found")
	}
	return result.MatchedCount, nil
}

func (c *Collection) Delete(ctx context.Context, filter any) (int64, error) {
	result, err := c.Collection.DeleteMany(ctx, filter, options.Delete())
	if err != nil {
		return 0, err
	}
	if result.DeletedCount == 0 {
		return 0, fmt.Errorf("DeleteOne result: %s ", "document not found")
	}
	return result.DeletedCount, nil
}

func (c *Collection) InsertOne(ctx context.Context, model any) (any, error) {
	result, err := c.Collection.InsertOne(ctx, model, options.InsertOne())
	if err != nil {
		return nil, err
	}
	return result.InsertedID, err
}

func (c *Collection) InsertMany(ctx context.Context, models []any) ([]any, error) {
	result, err := c.Collection.InsertMany(ctx, models, options.InsertMany())
	if err != nil {
		return nil, err
	}
	return result.InsertedIDs, err
}

func (c *Collection) Aggregate(ctx context.Context, pipeline any, result any) error {
	finder, err := c.Collection.Aggregate(ctx, pipeline, options.Aggregate())
	if err != nil {
		return err
	}
	if err := finder.All(ctx, &result); err != nil {
		return err
	}
	return nil
}

func (c *Collection) CreateIndexes(ctx context.Context, indexes []mongo.IndexModel) error {
	_, err := c.Collection.Indexes().CreateMany(ctx, indexes, options.CreateIndexes())
	return err
}

func (c *Collection) GetCollection() *mongo.Collection {
	return c.Collection
}
