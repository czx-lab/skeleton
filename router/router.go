package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppRouter struct{}

func New() *AppRouter {
	return &AppRouter{}
}

func (*AppRouter) Add(server *gin.Engine) {
	server.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello word!")
	})
}
