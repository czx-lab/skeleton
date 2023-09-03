package test

import (
	"context"
	_ "github.com/czx-lab/skeleton/internal/bootstrap"
	"github.com/czx-lab/skeleton/internal/variable"
	"testing"
)

func TestRedis(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Error(err)
		}
	}()
	redisClient := variable.Redis
	t.Log(redisClient.Get(context.Background(), "test").Result())
}
