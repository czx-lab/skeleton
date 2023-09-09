package controller

import (
	"fmt"
	"github.com/czx-lab/skeleton/app/event/entity"
	"github.com/czx-lab/skeleton/app/request"
	"github.com/czx-lab/skeleton/internal/variable"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Index struct {
	base
}

func (i *Index) Hello(ctx *gin.Context) {
	var param request.Foo
	data, err := variable.Event.Dispatch(&entity.FooEvent{
		Name: "hello",
	})
	fmt.Println(data, err)
	if err := i.base.Validate(ctx, &param); err == nil {
		fmt.Println(param)
		ctx.JSON(http.StatusOK, gin.H{"data": "Hello World"})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
	}
}
