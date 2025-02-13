package router

import "github.com/gin-gonic/gin"

type IRouter interface {
	Add(server *gin.Engine)
}

type Router struct {
	server *gin.Engine
}

func New(server *gin.Engine) *Router {
	return &Router{
		server: server,
	}
}

func (r *Router) AddRouter(routers IRouter) {
	routers.Add(r.server)
}
