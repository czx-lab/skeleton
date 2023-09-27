package server

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SocketOption struct {
	writeReadBufferSize   int
	heartbeatFailMaxTimes int
	writeDeadline         time.Duration
	readDeadline          time.Duration
	pingPeriod            time.Duration
	pingMsg               string
	handler               MessageHandler
	logger                *zap.Logger
}

type MessageHandler interface {
	OnMessage(message Message)
	OnError(key string, err error)
	OnClose(key string)
}

type SocketClientInterface interface {
	WriteMessage(message Message) error
	GetAllKeys() []string
	GetClientState(key string) ClientState
	Connect(ctx *gin.Context, subkey string)
}

type Message struct {
	MessageType int
	Subkeys     []string
	Data        []byte
}

type Socket struct {
	clients    map[string]*SocketClient
	unregister chan string
	opts       *SocketOption
}

func NewSocket(opts ...SocketOptionFunc) (SocketClientInterface, error) {
	sOpt := &SocketOption{}
	socket := &Socket{
		clients:    make(map[string]*SocketClient),
		unregister: make(chan string),
	}
	for _, opt := range opts {
		opt.apply(sOpt)
	}
	defaultOption(sOpt)
	socket.opts = sOpt
	go socket.listen()
	return socket, nil
}

func (s *Socket) listen() {
	for {
		select {
		case key := <-s.unregister:
			if client, ok := s.clients[key]; ok {
				delete(s.clients, key)
				close(client.send)
			}
		}
	}
}

func (s *Socket) Connect(ctx *gin.Context, subkey string) {
	client, ok := s.clients[subkey]
	if ok && client.state == OnlineState {
		return
	}
	s.clients[subkey] = NewSocketClient(ctx, subkey, s)
}

func (s *Socket) GetAllKeys() []string {
	clients := s.clients
	if len(clients) == 0 {
		return []string{}
	}
	keys := make([]string, 0, len(clients))
	for k, c := range clients {
		if c.state == OnlineState {
			keys = append(keys, k)
		}
	}
	return keys
}

func (s *Socket) GetClientState(key string) ClientState {
	client, ok := s.clients[key]
	if !ok {
		return OffLineState
	}
	return client.state
}

func (s *Socket) WriteMessage(message Message) error {
	if len(message.Subkeys) == 0 {
		for _, client := range s.clients {
			if client.state == OnlineState {
				client.send <- message.Data
			}
		}
	} else {
		for _, key := range message.Subkeys {
			client, ok := s.clients[key]
			if !ok && client.state == OffLineState {
				return errors.New("Connect does not exist")
			}
			client.send <- message.Data
		}
	}
	return nil
}

func defaultOption(opts *SocketOption) {
	if opts.pingPeriod == 0 {
		opts.pingPeriod = 20 * time.Second
	}
	if opts.writeDeadline == 0 {
		opts.writeDeadline = 35 * time.Second
	}
	if opts.writeReadBufferSize == 0 {
		opts.writeReadBufferSize = 20480
	}
	if opts.heartbeatFailMaxTimes == 0 {
		opts.heartbeatFailMaxTimes = 4
	}
	if opts.readDeadline == 0 {
		opts.readDeadline = 30 * time.Second
	}
}

type SocketOptionInterface interface {
	apply(*SocketOption)
}

type SocketOptionFunc func(opt *SocketOption)

func (f SocketOptionFunc) apply(opt *SocketOption) {
	f(opt)
}

func WithHandler(handler MessageHandler) SocketOptionFunc {
	return func(opt *SocketOption) {
		opt.handler = handler
	}
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
