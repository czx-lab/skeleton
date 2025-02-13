package test

import (
	"skeleton/internal/logx"
	"testing"

	"go.uber.org/zap"
)

func TestLogger(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("TestConfig filed:%v", err)
		}
	}()
	log := logx.NewLogx(
		logx.WithLevel(zap.ErrorLevel),
		logx.WithEncoding(logx.EncodingPlain),
		logx.WithMod(logx.ModFile),
	).Zap()
	log.Error("test")
}
