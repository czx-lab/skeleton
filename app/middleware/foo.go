package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Foo struct{}

func New() *Foo {
	return &Foo{}
}

func (f *Foo) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Println("Foo middleware exec...")
	}
}
