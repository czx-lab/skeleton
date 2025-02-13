package config

import (
	"czx/internal/queue/rocket"
	"czx/internal/server/xhttp"
	"czx/internal/stores/mongo"
	"czx/internal/stores/xmysql"
	"czx/internal/stores/xredis"
	"czx/internal/xlog"
)

type Config struct {
	HttpConf     xhttp.HttpConf
	MysqlConf    xmysql.MySqlConf
	MongoConf    mongo.MonConf
	RedisConf    xredis.RedisConf
	LogConf      xlog.LogConf
	RocketMqConf rocket.RocketConf
}
