package controller

import (
	"fmt"
	"net/http"
	"skeleton/internal/server"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Socket struct{}

func (s *Socket) Connect(ctx *gin.Context) {
	client, err := server.NewSocket(ctx)
	if err != nil {
		ctx.JSON(http.StatusAccepted, gin.H{"message": err})
		return
	}
	client.ReadPump(&socketHandler{
		client: client,
	})
}

type socketHandler struct {
	client *server.SocketClient
}

func (s *socketHandler) OnMessage(messageType int, data []byte) {
	fmt.Println(fmt.Sprintf("mt: %v，data: %s", messageType, data))
	s.client.SendMessage(websocket.TextMessage, "服务端回复消息")
}

func (s *socketHandler) OnError(err error) {
	fmt.Println(fmt.Sprintf("socket err: %s", err))
}

func (s *socketHandler) OnClose() {
	fmt.Println("socket closed.")
}
