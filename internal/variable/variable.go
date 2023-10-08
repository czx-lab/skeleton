package variable

import (
	"log"
	"os"
	"skeleton/internal/event"
	"skeleton/internal/mongo"
	"strings"

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
	DB       *gorm.DB
	MongoDB  *mongo.MongoDB
	Redis    *redis.Client
	Crontab  *crontab.Crontab
	MQ       mq.Interface
	Event    *event.Event
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
