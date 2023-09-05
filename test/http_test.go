package test

import (
	"testing"

	"github.com/czx-lab/skeleton/internal/server"
)

func TestHttp(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Log(err)
		}
	}()
	http := server.New(server.WithAddress(":8888"))
	if err := http.Run(); err != nil {
		t.Log(err)
	}
	t.Log("success")
}
