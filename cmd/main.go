package main

import (
	"fmt"
	_ "github.com/czx-lab/skeleton/internal/bootstrap"
	"github.com/czx-lab/skeleton/internal/server"
	"github.com/czx-lab/skeleton/internal/variable"
	"github.com/czx-lab/skeleton/router"
)

func main() {
	port := variable.Config.GetString("HttpServer.Port")
	mode := variable.Config.GetString("HttpServer.Mode")
	http := server.New(server.WithPort(port), server.WithMode(mode))
	http.SetRouters(router.New())
	if err := http.Run(); err != nil {
		fmt.Println(err)
	}
}
