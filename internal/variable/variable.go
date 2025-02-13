package variable

import (
	"log"
	"os"
	"skeleton/internal/elasticsearch"
	"skeleton/internal/event"
	"skeleton/internal/mongo"
	"strings"

	conf "skeleton/config"
	"skeleton/internal/config"
	"skeleton/internal/crontab"
	"skeleton/internal/mq"
	"skeleton/internal/variable/consts"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	BasePath string
	Log      *zap.Logger
	Config   *config.Config
	AppConf  *conf.Config
	DB       *gorm.DB
	MongoDB  *mongo.MongoDB
	Redis    *redis.Client
	Crontab  *crontab.Crontab
	Amqp     mq.RabbitMQInterface
	Event    *event.Event
	Elastic  *elasticsearch.Elasticsearch

	// RocketMQ 目前官方RocketMQ Golang SDK一些功能尚未完善，暂时不可用
	RocketMQ mq.Interface
)

func init() {
	if curPath, err := os.Getwd(); err == nil {
		if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-test") {
			BasePath = strings.Replace(strings.Replace(curPath, `\test`, "", 1), `/test`, "", 1)
		} else {
			BasePath = curPath
		}
	} else {
		log.Fatal(consts.ErrorsBasePath)
	}
}
