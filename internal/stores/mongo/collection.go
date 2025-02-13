package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ICollection interface {
	SelectPage(ctx context.Context, filter any, sort any, skip, limit int64) (int64, []any, error)
	SelectList(ctx context.Context, filter any, sort interface{}) ([]any, error)
	SelectOne(ctx context.Context, filter any) (interface{}, error)
	Count(ctx context.Context, filter any) (int64, error)
	UpdateOne(ctx context.Context, filter, update any) (int64, error)
	UpdateMany(ctx context.Context, filter, update any) (int64, error)
	Delete(ctx context.Context, filter any) (int64, error)
	InsertOne(ctx context.Context, model any) (any, error)
	InsertMany(ctx context.Context, models []any) ([]any, error)
	Aggregate(ctx context.Context, pipeline any, result any) error
	CreateIndexes(ctx context.Context, indexes []mongo.IndexModel) error
	GetCollection() *mongo.Collection
}

type Collection struct {
	coll *mongo.Collection
}

// Aggregate implements ICollection.
func (c *Collection) Aggregate(ctx context.Context, pipeline any, result any) error {
	finder, err := c.coll.Aggregate(ctx, pipeline, options.Aggregate())
	if err != nil {
		return err
	}
	if err := finder.All(ctx, &result); err != nil {
		return err
	}
	return nil
}

// CreateIndexes implements ICollection.
func (c *Collection) CreateIndexes(ctx context.Context, indexes []mongo.IndexModel) error {
	_, err := c.coll.Indexes().CreateMany(ctx, indexes, options.CreateIndexes())
	return err
}

// Delete implements ICollection.
func (c *Collection) Delete(ctx context.Context, filter any) (int64, error) {
	result, err := c.coll.DeleteMany(ctx, filter, options.Delete())
	if err != nil {
		return 0, err
	}
	if result.DeletedCount == 0 {
		return 0, fmt.Errorf("DeleteOne result: %s ", "document not found")
	}
	return result.DeletedCount, nil
}

// GetCollection implements ICollection.
func (c *Collection) GetCollection() *mongo.Collection {
	return c.coll
}

// InsertMany implements ICollection.
func (c *Collection) InsertMany(ctx context.Context, models []any) ([]any, error) {
	result, err := c.coll.InsertMany(ctx, models, options.InsertMany())
	if err != nil {
		return nil, err
	}
	return result.InsertedIDs, err
}

// InsertOne implements ICollection.
func (c *Collection) InsertOne(ctx context.Context, model any) (any, error) {
	result, err := c.coll.InsertOne(ctx, model, options.InsertOne())
	if err != nil {
		return nil, err
	}
	return result.InsertedID, err
}

// Count implements ICollection.
func (c *Collection) Count(ctx context.Context, filter any) (int64, error) {
	return c.coll.CountDocuments(ctx, filter)
}

// SelectList implements ICollection.
func (c *Collection) SelectList(ctx context.Context, filter any, sort interface{}) ([]any, error) {
	var err error

	opts := options.Find().SetSort(sort)
	finder, err := c.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, 0)
	if err := finder.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, err
}

// SelectOne implements ICollection.
func (c *Collection) SelectOne(ctx context.Context, filter any) (interface{}, error) {
	result := new(interface{})
	err := c.coll.FindOne(ctx, filter, options.FindOne()).Decode(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// SelectPage implements ICollection.
func (c *Collection) SelectPage(ctx context.Context, filter any, sort any, skip int64, limit int64) (int64, []any, error) {
	var err error

	resultCount, err := c.coll.CountDocuments(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	opts := options.Find().SetSort(sort).SetSkip(skip).SetLimit(limit)
	finder, err := c.coll.Find(ctx, filter, opts)
	if err != nil {
		return resultCount, nil, err
	}

	result := make([]interface{}, 0)
	if err := finder.All(ctx, &result); err != nil {
		return resultCount, nil, err
	}
	return resultCount, result, nil
}

// UpdateMany implements ICollection.
func (c *Collection) UpdateMany(ctx context.Context, filter any, update any) (int64, error) {
	result, err := c.coll.UpdateMany(ctx, filter, update, options.Update())
	if err != nil {
		return 0, err
	}
	if result.MatchedCount == 0 {
		return 0, fmt.Errorf("update result: %s ", "document not found")
	}
	return result.MatchedCount, nil
}

// UpdateOne implements ICollection.
func (c *Collection) UpdateOne(ctx context.Context, filter any, update any) (int64, error) {
	result, err := c.coll.UpdateOne(ctx, filter, update, options.Update())
	if err != nil {
		return 0, err
	}
	if result.MatchedCount == 0 {
		return 0, fmt.Errorf("update result: %s ", "document not found")
	}
	return result.MatchedCount, nil
}

var _ ICollection = (*Collection)(nil)
