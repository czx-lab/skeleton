package test

import (
	"czx/internal/server/xhttp"
	"testing"

	"github.com/spf13/viper"
)

func xhttpConf() (conf xhttp.HttpConf) {
	viper.SetConfigFile("../config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return
	}
	if err := viper.Unmarshal(&conf); err != nil {
		return
	}
	return
}

func TestXHttp(t *testing.T) {
	xhp := xhttp.New(xhttpConf())
	xhp.Start()
}
