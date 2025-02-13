package test

import (
	"fmt"
	"czx/app/config"
	"czx/internal/conf"
	"testing"

	"github.com/spf13/viper"
)

func TestConf(t *testing.T) {
	viper.SetConfigFile("../config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		return
	}
	var conf config.Config
	if err := viper.Unmarshal(&conf); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", conf)
}

func TestXConf(t *testing.T) {
	xconf := conf.New("../config/config.yaml")
	tout, err := xconf.Get("Timeout")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tout)
}
