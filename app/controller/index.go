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
	var param request.Foo
	if err := variable.Event.Dispatch(&entity.FooEvent{
		Name: "hello",
	}); err != nil {
		fmt.Println(err)
	}

	// (&producer.FooProducer{}).SendMessage([]byte("foo message"))

	if err := i.base.Validate(ctx, &param); err == nil {
		fmt.Println(param)
		ctx.JSON(http.StatusOK, gin.H{"data": "Hello World"})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
	}
}
