package logic

import (
	"net/http"
	"czx/app/svc"

	"github.com/gin-gonic/gin"
)

type Index struct {
	svcCtx *svc.AppContext
}

func NewIndex(ctx *svc.AppContext) *Index {
	return &Index{
		svcCtx: ctx,
	}
}

func (i *Index) Demo(ctx *gin.Context) {
	ctx.String(http.StatusOK, "hello word!")
}
