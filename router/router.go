package router

import (
	"net/http"

	"skeleton/app/controller"
	"skeleton/app/middleware"

	"github.com/gin-gonic/gin"
	"skeleton/internal/server"
)

type AppRouter struct {
	server server.HttpServer
}

func New(server server.HttpServer) *AppRouter {
	server.SetMiddleware(&middleware.Foo{}, &middleware.Cors{})
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
