package main

import (
	"czx/app"
	"flag"
)

var configFile = flag.String("f", "config/config.yaml", "the config file")

func main() {
	flag.Parse()

	_ = app.New(*configFile).Start()
}
