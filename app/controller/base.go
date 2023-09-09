package controller

import (
	"github.com/czx-lab/skeleton/app/request"
	"github.com/czx-lab/skeleton/internal/variable/consts"
	"github.com/gin-gonic/gin"
	"log"
)

var validator *request.Request

type base struct {
}

func init() {
	var err error
	validator, err = request.New()
	if err != nil {
		log.Fatal(consts.ErrorInitConfig)
	}
}

func (base) Validate(ctx *gin.Context, param any) map[string]string {
	return validator.Validator(ctx, param)
}
