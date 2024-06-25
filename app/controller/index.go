package controller

import (
	"fmt"
	"net/http"
	"skeleton/app/event/entity"
	"skeleton/app/request"
	"skeleton/internal/variable"

	"github.com/gin-gonic/gin"
)

type Index struct {
	base
}

func (i *Index) Hello(ctx *gin.Context) {
	if err := variable.Event.Dispatch(&entity.FooEvent{
		Name: "hello",
	}); err != nil {
		fmt.Println(err)
	}

	// (&producer.FooProducer{}).SendMessage([]byte("foo message"))
	var bodyFooParam request.FooBody
	if err := i.base.Validate(ctx, &bodyFooParam); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	var param request.Foo
	if err := i.base.Validate(ctx, &param); err == nil {
		fmt.Println(param)
		ctx.JSON(http.StatusOK, gin.H{"data": "Hello World"})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
	}
}
