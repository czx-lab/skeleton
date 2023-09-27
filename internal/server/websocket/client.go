package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ClientState int

const (
	_ ClientState = iota
	OnlineState
	OffLineState
)

type SocketClient struct {
	key                string
	conn               *websocket.Conn
	send               chan []byte
	heartbeatFailTimes int
	socket             *Socket
	state              ClientState
}

func NewSocketClient(ctx *gin.Context, key string, socket *Socket) *SocketClient {
	client := &SocketClient{
		key:    key,
		socket: socket,
		state:  OnlineState,
	}
	client.upGrader(ctx, socket.opts)
	return client
}

func (s *SocketClient) readPump() {
	defer func() {
		if err := recover(); err != nil {
			s.socket.opts.handler.OnError(s.key, errors.New(fmt.Sprintf("%v", err)))
		}
		s.close()
	}()
	_ = s.conn.SetReadDeadline(time.Now().Add(s.socket.opts.readDeadline))
	s.conn.SetPongHandler(func(receivedPong string) error {
		if s.socket.opts.readDeadline > time.Nanosecond {
			_ = s.conn.SetReadDeadline(time.Now().Add(s.socket.opts.readDeadline))
		} else {
			_ = s.conn.SetReadDeadline(time.Time{})
		}
		return nil
	})
	for {
		if mt, data, err := s.conn.ReadMessage(); err != nil {
			s.socket.opts.handler.OnError(s.key, err)
			break
		} else {
			message := Message{
				MessageType: mt,
				Data:        data,
				Subkeys:     []string{s.key},
			}
			s.socket.opts.handler.OnMessage(message)
		}
	}
}

func (s *SocketClient) writePump() {
	ticker := time.NewTicker(s.socket.opts.pingPeriod)
	defer func() {
		if err := recover(); err != nil {
			s.socket.opts.handler.OnError(s.key, errors.New(fmt.Sprintf("%v", err)))
		}
		s.close()
	}()
	for {
		select {
		case message, ok := <-s.send:
			s.conn.SetWriteDeadline(time.Now().Add(s.socket.opts.readDeadline))
			if !ok {
				s.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := s.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			if err := s.conn.SetWriteDeadline(time.Now().Add(s.socket.opts.writeDeadline)); err != nil {
				return
			}
			if err := s.conn.WriteMessage(websocket.PingMessage, []byte(s.socket.opts.pingMsg)); err != nil {
				s.heartbeatFailTimes++
				if s.heartbeatFailTimes > s.socket.opts.heartbeatFailMaxTimes {
					return
				}
			} else {
				if s.heartbeatFailTimes > 0 {
					s.heartbeatFailTimes--
				}
			}
		}
	}
}

func (s *SocketClient) close() {
	if s.state == OnlineState {
		s.state = OffLineState
		s.socket.unregister <- s.key
		s.conn.Close()
		s.socket.opts.handler.OnClose(s.key)
	}
}

func (s *SocketClient) upGrader(context *gin.Context, opts *SocketOption) {
	upGrader := websocket.Upgrader{
		ReadBufferSize:  opts.writeReadBufferSize,
		WriteBufferSize: opts.writeReadBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	wsConn, err := upGrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		if opts.logger != nil {
			opts.logger.Error(err.Error())
		} else {
			log.Panicln(err)
		}
		return
	}
	s.conn = wsConn
	s.send = make(chan []byte, opts.writeReadBufferSize)
	go s.readPump()
	go s.writePump()
}
