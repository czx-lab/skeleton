package mongo

import (
	"context"
	"github.com/czx-lab/skeleton/internal/mongo/collection"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

const (
	defaultTimeout = 50 * time.Second
)

type MongoDB struct {
	client *mongo.Client
}

func New(opts ...*options.ClientOptions) (client *MongoDB, err error) {
	client = &MongoDB{}
	err = client.InitMongo(opts...)
	return
}

func (m *MongoDB) InitMongo(opts ...*options.ClientOptions) error {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	m.client, err = mongo.Connect(ctx, opts...)
	if err != nil {
		return err
	}
	if err := m.client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	return nil
}

func (m *MongoDB) CreateMongoCollection(dbName, colName string) collection.CollectionInterface {
	dataBase := m.client.Database(dbName)
	return &collection.Collection{
		DbName:     dbName,
		ColName:    colName,
		DataBase:   dataBase,
		Collection: dataBase.Collection(colName),
	}
}
