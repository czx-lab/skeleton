package router

import (
	"net/http"
	"czx/app/logic"
	"czx/app/svc"
	"czx/internal/server/router"

	"github.com/gin-gonic/gin"
)

type ApiRouter struct {
	ctx *svc.AppContext
}

var _ router.IRouter = (*ApiRouter)(nil)

func New(ctx *svc.AppContext) *ApiRouter {
	return &ApiRouter{ctx}
}

func (a *ApiRouter) Add(server *gin.Engine) {
	server.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello word!")
	})

	index := logic.NewIndex(a.ctx)
	server.GET("/demo", index.Demo)
}
