package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AppException struct {
	log *zap.Logger
}

func New(logger *zap.Logger) *AppException {
	return &AppException{
		log: logger,
	}
}

func (app *AppException) Handle() gin.HandlerFunc {
	DefaultErrorWriter := &PanicException{
		logger: app.log,
	}
	return gin.RecoveryWithWriter(DefaultErrorWriter, func(c *gin.Context, err any) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("%s", err),
		})
	})
}

// PanicException  panic等异常记录
type PanicException struct {
	logger *zap.Logger
}

func (p *PanicException) Write(b []byte) (n int, err error) {
	errStr := string(b)
	err = errors.New(errStr)
	if p.logger != nil {
		p.logger.Error("Internal Server Error: ", zap.String("strace", errStr))
	}
	return len(errStr), err
}
