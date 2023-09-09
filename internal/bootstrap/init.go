package bootstrap

import (
	AppEvent "github.com/czx-lab/skeleton/internal/event"
	"log"
	"time"

	"github.com/czx-lab/skeleton/app/event"
	"github.com/czx-lab/skeleton/app/task"
	"github.com/czx-lab/skeleton/internal/config"
	"github.com/czx-lab/skeleton/internal/config/driver"
	"github.com/czx-lab/skeleton/internal/crontab"
	"github.com/czx-lab/skeleton/internal/logger"
	"github.com/czx-lab/skeleton/internal/mq"
	"github.com/czx-lab/skeleton/internal/redis"
	"github.com/czx-lab/skeleton/internal/variable"
	"github.com/czx-lab/skeleton/internal/variable/consts"
)

func init() {
	var err error
	if variable.Config, err = config.New(driver.New(), config.Options{
		BasePath: variable.BasePath,
	}); err != nil {
		log.Fatal(consts.ErrorInitConfig)
	}
	if variable.Log, err = logger.New(
		logger.WithDebug(true),
		logger.WithEncode("json"),
		logger.WithFilename(variable.BasePath+"/storage/logs/system.log"),
	); err != nil {
		log.Fatal(consts.ErrorInitLogger)
	}
	if variable.DB, err = InitMysql(); err != nil {
		log.Fatal(consts.ErrorInitDb)
	}
	if err = InitMongo(); err != nil {
		log.Fatal(consts.ErrorInitMongoDb)
	}
	redisConfig := variable.Config.Get("Redis").(map[string]any)
	if redisConfig != nil && !redisConfig["disabled"].(bool) {
		variable.Redis = redis.New(
			redis.WithAddr(redisConfig["addr"].(string)),
			redis.WithPwd(redisConfig["pwd"].(string)),
			redis.WithDb(redisConfig["db"].(int)),
			redis.WithPoolSize(redisConfig["poolsize"].(int)),
			redis.WithMaxIdleConn(redisConfig["maxidleconn"].(int)),
			redis.WithMinIdleConn(redisConfig["minidleconn"].(int)),
			redis.WithMaxLifetime(time.Duration(redisConfig["maxlifetime"].(int))),
			redis.WithMaxIdleTime(time.Duration(redisConfig["maxidletime"].(int))),
		)
	}
	if variable.Config.GetBool("Crontab.Enable") {
		variable.Crontab = crontab.New()
		variable.Crontab.AddFunc(task.New().Tasks()...)
		variable.Crontab.Start()
	}
	if variable.Config.GetBool("MQ.Enable") {
		if variable.MQ, err = mq.New(
			mq.WithNameServers(variable.Config.GetStringSlice("MQ.Servers")),
			mq.WithGroupId(variable.Config.GetString("MQ.GroupId")),
			mq.WithConsumptionSize(variable.Config.GetInt("MQ.ConsumptionSize")),
			mq.WithRetries(variable.Config.GetInt("MQ.Retries")),
		); err != nil {
			log.Fatal(consts.ErrorInitMQ)
		}
	}
	variable.Event = AppEvent.New()
	if err = (&event.Event{}).Init(); err != nil {
		log.Fatal(consts.ErrorInitEvent)
	}
}
