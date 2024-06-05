package test

import (
	"fmt"
	"skeleton/internal/server"
	"skeleton/router"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHttp(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Log(err)
		}
	}()
	http := server.New(
		server.WithPort(":8888"),
		server.WithMode(gin.DebugMode),
		server.WithAfterFunc(func() {
			fmt.Println("启动之后执行")
		}),
	)
	http.SetRouters(router.New(http)).Run()
}
