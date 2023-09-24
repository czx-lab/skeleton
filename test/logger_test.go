package test

import (
	"skeleton/internal/logger"
	"testing"
)

func TestLogger(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("TestConfig filed:%v", err)
		}
	}()
	log, err := logger.New(logger.WithDebug(false), logger.WithEncode("json"))
	if err != nil {
		t.Errorf("TestConfig filed:%v", err)
	} else {
		log.Info("测试123131")
	}
}
