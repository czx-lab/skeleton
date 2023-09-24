package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"skeleton/app/request"
	"skeleton/internal/variable/consts"
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
