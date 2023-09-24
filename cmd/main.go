package main

import (
	"fmt"

	_ "skeleton/internal/bootstrap"
	"skeleton/internal/server"
	"skeleton/internal/variable"
	"skeleton/router"
)

func main() {
	port := variable.Config.GetString("HttpServer.Port")
	mode := variable.Config.GetString("HttpServer.Mode")
	http := server.New(server.WithPort(port), server.WithMode(mode))
	http.SetRouters(router.New(http))
	if err := http.Run(); err != nil {
		fmt.Println(err)
	}
}
