package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

type SocketClient struct {
	Conn   *websocket.Conn
	option SocketOption
	sync.RWMutex
	Send               chan []byte
	HeartbeatFailTimes int
}

type SocketOption struct {
	writeReadBufferSize   int
	heartbeatFailMaxTimes int
	writeDeadline         time.Duration
	readDeadline          time.Duration
	pingPeriod            time.Duration
	pingMsg               string
}

func NewSocket(context *gin.Context, opts ...SocketOptionFunc) (*SocketClient, error) {
	sOpt := &SocketOption{}
	for _, opt := range opts {
		opt.apply(sOpt)
	}
	client := &SocketClient{
		option: *sOpt,
	}
	if err := client.upGrader(context); err != nil {
		return nil, err
	}
	return client, nil
}

type MessageHandler interface {
	OnMessage(messageType int, data []byte)
	OnError(err error)
	OnClose()
}

// ReadPump 消息处理
func (s *SocketClient) ReadPump(handler MessageHandler) {
	defer func() {
		handler.OnClose()
	}()
	for {
		if mt, data, err := s.Conn.ReadMessage(); err != nil {
			handler.OnError(err)
			break
		} else {
			handler.OnMessage(mt, data)
		}
	}
}

// SendMessage 发送消息
func (s *SocketClient) SendMessage(messageType int, message string) error {
	s.Lock()
	defer func() {
		s.Unlock()
	}()
	if err := s.Conn.SetWriteDeadline(time.Now().Add(s.option.writeDeadline)); err != nil {
		return err
	}
	if err := s.Conn.WriteMessage(messageType, []byte(message)); err != nil {
		return err
	}
	return nil
}

func (s *SocketClient) upGrader(context *gin.Context) error {
	upGrader := websocket.Upgrader{
		ReadBufferSize:  s.option.writeReadBufferSize,
		WriteBufferSize: s.option.writeReadBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	wsConn, err := upGrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		return err
	}
	s.Conn = wsConn
	s.Send = make(chan []byte, s.option.writeReadBufferSize)
	go s.heartbeat()
	return nil
}

func (s *SocketClient) heartbeat() {
	ticker := time.NewTicker(s.option.pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	_ = s.Conn.SetReadDeadline(time.Now().Add(s.option.readDeadline))
	s.Conn.SetPongHandler(func(receivedPong string) error {
		if s.option.readDeadline > time.Nanosecond {
			_ = s.Conn.SetReadDeadline(time.Now().Add(s.option.readDeadline))
		} else {
			_ = s.Conn.SetReadDeadline(time.Time{})
		}
		return nil
	})
	for {
		select {
		case <-ticker.C:
			if err := s.SendMessage(websocket.PingMessage, s.option.pingMsg); err != nil {
				s.HeartbeatFailTimes++
				if s.HeartbeatFailTimes > s.option.heartbeatFailMaxTimes {
					return
				}
			} else {
				if s.HeartbeatFailTimes > 0 {
					s.HeartbeatFailTimes--
				}
			}
		}
	}
}

type SocketOptionInterface interface {
	apply(*SocketOption)
}

type SocketOptionFunc func(opt *SocketOption)

func (f SocketOptionFunc) apply(opt *SocketOption) {
	f(opt)
}

func WithWriteReadBufferSize(size int) SocketOptionFunc {
	return func(opt *SocketOption) {
		opt.writeReadBufferSize = size
	}
}

func WithReadDeadline(deadline time.Duration) SocketOptionFunc {
	return func(opt *SocketOption) {
		opt.readDeadline = deadline
	}
}

func WithHeartbeatFailMaxTimes(heartbeatFailMaxTimes int) SocketOptionFunc {
	return func(opt *SocketOption) {
		opt.heartbeatFailMaxTimes = heartbeatFailMaxTimes
	}
}

func WithWriteDeadline(writeDeadline time.Duration) SocketOptionFunc {
	return func(opt *SocketOption) {
		opt.writeDeadline = writeDeadline
	}
}

func WithPingPeriod(pingPeriod time.Duration) SocketOptionFunc {
	return func(opt *SocketOption) {
		opt.pingPeriod = pingPeriod
	}
}

func WithPingMsg(pingMsg string) SocketOptionFunc {
	return func(opt *SocketOption) {
		opt.pingMsg = pingMsg
	}
}
