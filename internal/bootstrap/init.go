package bootstrap

import (
	"github.com/czx-lab/skeleton/internal/config"
	"github.com/czx-lab/skeleton/internal/config/driver"
	"github.com/czx-lab/skeleton/internal/logger"
	"github.com/czx-lab/skeleton/internal/variable"
	"github.com/czx-lab/skeleton/internal/variable/consts"
	"log"
)

func init() {
	var err error
	if variable.Config, err = config.New(driver.New(), config.Options{
		BasePath: variable.BasePath,
	}); err != nil {
		log.Fatal(consts.ErrorInitConfig)
	}
	if variable.Log, err = logger.New(logger.WithDebug(false), logger.WithEncode("json")); err != nil {
		log.Fatal(consts.ErrorInitLogger)
	}
	if variable.DB, err = InitMysql(); err != nil {
		log.Fatal(consts.ErrorInitDb)
	}
}
