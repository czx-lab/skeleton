package driver

import (
	"github.com/czx-lab/skeleton/internal/config"
	constants "github.com/czx-lab/skeleton/internal/constants/config"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"time"
)

type ViperConfig struct {
	viper  *viper.Viper
	option config.Options
}

var _ constants.DriverInterface = (*ViperConfig)(nil)

// 由于 vipver 包本身对于文件的变化事件有一个bug，相关事件会被回调两次
// 常年未彻底解决，相关的 issue 清单：https://github.com/spf13/viper/issues?q=OnConfigChange
// 设置一个内部全局变量，记录配置文件变化时的时间点，如果两次回调事件事件差小于1秒，我们认为是第二次回调事件，而不是人工修改配置文件
// 这样就避免了 viper 包的这个bug
var lastChangeTime time.Time

func init() {
	lastChangeTime = time.Now()
}

func New() *ViperConfig {
	return &ViperConfig{}
}

// Apply 创建实例
func (v *ViperConfig) Apply(option config.Options) error {
	v.option = option
	viperConfig := viper.New()
	viperConfig.AddConfigPath(option.BasePath + "/config")
	if len(option.Filename) == 0 {
		viperConfig.SetConfigName("config")
	} else {
		viperConfig.SetConfigName(option.Filename)
	}
	viperConfig.SetConfigType(option.Cate)
	if err := viperConfig.ReadInConfig(); err != nil {
		return err
	}
	v.viper = viperConfig
	return nil
}

// Listen 监听文件变化
func (v *ViperConfig) Listen() {
	v.viper.OnConfigChange(func(in fsnotify.Event) {
		if time.Now().Sub(lastChangeTime).Seconds() >= 1 {
			if in.Op.String() == "WRITE" {
				// 清除cache内的配置
				v.option.Cache.FuzzyDelete(v.option.CachePrefix)
				lastChangeTime = time.Now()
			}
		}
	})
}

func (v *ViperConfig) Get(key string) any {
	return v.viper.Get(key)
}

func (v *ViperConfig) Set(key string, value any) bool {
	v.viper.Set(key, value)
	return true
}

func (v *ViperConfig) Has(key string) bool {
	return v.viper.IsSet(key)
}
