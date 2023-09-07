package router

import (
	"github.com/czx-lab/skeleton/app/controller"
	"github.com/czx-lab/skeleton/app/middleware"
	"net/http"

	"github.com/czx-lab/skeleton/internal/server"
	"github.com/gin-gonic/gin"
)

type AppRouter struct {
	server server.HttpServer
}

func New(server server.HttpServer) *AppRouter {
	server.SetMiddleware(&middleware.Foo{})
	return &AppRouter{
		server,
	}
}

func (*AppRouter) Add(server *gin.Engine) {
	server.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello word!")
	})

	index := &controller.Index{}
	server.GET("/hello", index.Hello)
}
