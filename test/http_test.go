package test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"skeleton/internal/server"
	"skeleton/router"
)

func TestHttp(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Log(err)
		}
	}()
	http := server.New(server.WithPort(":8888"), server.WithMode(gin.DebugMode))
	http.SetRouters(router.New(http))
	if err := http.Run(); err != nil {
		t.Log(err)
	}
	t.Log("success")
}
