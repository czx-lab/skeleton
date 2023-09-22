package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Cors struct{}

func (*Cors) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		origin := ctx.Request.Header.Get("Origin")
		if origin != "" {
			//接收客户端发送的origin
			ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			ctx.Header("Access-Control-Allow-Headers", "Authorization, content-type, Content-Length, X-CSRF-Token, Token, session, Access-Control-Allow-Headers, account")
			// 允许浏览器（客户端）可以解析的头部
			ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			ctx.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie
			ctx.Header("Access-Control-Allow-Credentials", "true")
			ctx.Set("Content-Type", "application/json")
		}

		//允许类型校验
		if method == "OPTIONS" {
			ctx.JSON(http.StatusOK, "success")
		}
		ctx.Next()
	}
}
