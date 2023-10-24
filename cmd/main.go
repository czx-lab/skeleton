package main

import (
	_ "skeleton/internal/bootstrap"
	"skeleton/internal/server"
	"skeleton/internal/variable"
	"skeleton/router"
)

func main() {
	port := variable.Config.GetString("HttpServer.Port")
	mode := variable.Config.GetString("HttpServer.Mode")
	http := server.New(
		server.WithPort(port),
		server.WithMode(mode),
		server.WithLogger(variable.Log),
	)
	http.SetRouters(router.New(http)).Run()
}
