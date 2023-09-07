package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Index struct {
}

func (*Index) Hello(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"data": "Hello World"})
}
