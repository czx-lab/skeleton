package server

import (
	"io"

	"github.com/czx-lab/skeleton/internal/server/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Http struct {
	engine  *gin.Engine
	logger  *zap.Logger
	mode    string
	address []string
}

type Option interface {
	apply(http *Http)
}

type OptionFunc func(http *Http)

func (f OptionFunc) apply(http *Http) {
	f(http)
}

func New(opts ...Option) *Http {
	httpClass := &Http{}
	for _, opt := range opts {
		opt.apply(httpClass)
	}
	httpClass.defaultOption()
	httpClass.engine = httpClass.setServerEngine()
	return httpClass
}

func (h *Http) GetServerEngine() *gin.Engine {
	return h.engine
}

func (h *Http) setServerEngine() (engine *gin.Engine) {
	switch h.mode {
	case gin.DebugMode:
		engine = gin.Default()
	default:
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
		engine = gin.New()
		engine.Use(gin.Logger(), middleware.CustomRecovery(h.logger))
	}
	return
}

func (h *Http) defaultOption() {
	if h.mode == "" {
		h.mode = gin.DebugMode
	}
}

func (h *Http) Run() error {
	return h.engine.Run(h.address...)
}

func WithMode(mode string) Option {
	return OptionFunc(func(http *Http) {
		http.mode = mode
	})
}

func WithLogger(logger *zap.Logger) Option {
	return OptionFunc(func(http *Http) {
		http.logger = logger
	})
}

func WithAddress(address ...string) Option {
	return OptionFunc(func(http *Http) {
		http.address = address
	})
}
