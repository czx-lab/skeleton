package mongo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const timeout = 3 * time.Second

type MonConf struct {
	Uri  string
	Pool struct {
		Max         uint64
		Min         uint64
		MaxConn     uint64
		MaxIdleTime int64
	}
}

type Mon struct {
	conf MonConf

	instance *mongo.Client
}

func New(conf MonConf) *Mon {
	defaultConf(&conf)

	mon := &Mon{
		conf: conf,
	}
	mon.instance = mon.DB()
	return mon
}

func (m *Mon) DB() *mongo.Client {
	if m.instance == nil {
		m.instance = instance(m.conf)
	}
	return m.instance
}

func (m *Mon) Collection(db, coll string) ICollection {
	dataBase := m.instance.Database(db)
	return &Collection{
		coll: dataBase.Collection(coll),
	}
}

func instance(conf MonConf) *mongo.Client {
	idleTime := time.Duration(conf.Pool.MaxIdleTime) * time.Second
	option := options.Client().ApplyURI(conf.Uri).SetMaxPoolSize(conf.Pool.Max).SetMinPoolSize(conf.Pool.Min)
	option = option.SetMaxConnIdleTime(idleTime).SetMaxConnecting(conf.Pool.MaxConn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, option)
	if err != nil {
		log.Fatalf("mongo connect error: %v", err.Error())
		return nil
	}
	return client
}

func defaultConf(conf *MonConf) {
	if conf.Pool.Max == 0 {
		conf.Pool.Max = 10
	}
}
