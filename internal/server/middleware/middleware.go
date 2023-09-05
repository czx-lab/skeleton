package middleware

import "github.com/gin-gonic/gin"

type Interface interface {
	Handle() gin.HandlerFunc
}
