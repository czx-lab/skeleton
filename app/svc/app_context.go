package svc

import (
	"czx/app/config"
	"czx/internal/conf"
	"czx/internal/event"
	"czx/internal/queue/rocket"
	"czx/internal/stores/mongo"
	"czx/internal/stores/xmysql"
	"czx/internal/stores/xredis"
	"czx/internal/xlog"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AppContext struct {
	Config config.Config

	DB    *gorm.DB
	MDB   *mongo.Mon
	RDB   redis.UniversalClient
	Log   *zap.Logger
	Event *event.Event

	Rocket *rocket.Rocket
}

func NewAppContext(path string) (ctx *AppContext) {
	var c config.Config
	conf.New(path).Load(&c)

	return &AppContext{
		Config: c,

		DB:  xmysql.NewMysql(c.MysqlConf).DB(),
		RDB: xredis.New(c.RedisConf).DB(),
		MDB: mongo.New(c.MongoConf),
		Log: xlog.New(c.LogConf).Logger(),

		Event: event.New(),

		Rocket: rocket.New(c.RocketMqConf),
	}
}
