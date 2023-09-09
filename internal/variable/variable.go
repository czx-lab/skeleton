package variable

import (
	"github.com/czx-lab/skeleton/internal/event"
	"log"
	"os"
	"strings"

	"github.com/czx-lab/skeleton/internal/config"
	"github.com/czx-lab/skeleton/internal/crontab"
	"github.com/czx-lab/skeleton/internal/mq"
	"github.com/czx-lab/skeleton/internal/variable/consts"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	BasePath string
	Log      *zap.Logger
	Config   *config.Config
	DB       *gorm.DB
	Redis    *redis.Client
	Crontab  *crontab.Crontab
	MQ       *mq.RocketMQ
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
