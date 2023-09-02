package variable

import (
	"github.com/czx-lab/skeleton/internal/config"
	"github.com/czx-lab/skeleton/internal/variable/consts"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log"
	"os"
	"strings"
)

var (
	BasePath string
	Log      *zap.Logger
	Config   *config.Config
	DB       *gorm.DB
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
