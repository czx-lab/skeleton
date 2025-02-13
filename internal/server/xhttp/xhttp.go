package xhttp

import (
	"context"
	"czx/internal/constants"
	"czx/internal/server/router"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type IServer interface {
	constants.IService

	Engine() *gin.Engine
	RegisterHandlers(router.IRouter) IServer
}

type HttpConf struct {
	Mode     string
	Addr     string
	TemplDir string
	Timeout  int64
}

type XHttp struct {
	conf HttpConf

	engine *gin.Engine
	router *router.Router
}

func New(conf HttpConf) IServer {
	defaultConf(&conf)

	xhp := &XHttp{conf: conf}
	xhp.router = router.New(xhp.Engine())
	return xhp
}

// RegisterHandlers implements IServer.
func (x *XHttp) RegisterHandlers(routers router.IRouter) IServer {
	x.router.AddRouter(routers)
	return x
}

// Engine implements IServer.
func (x *XHttp) Engine() *gin.Engine {
	if x.engine == nil {
		switch x.conf.Mode {
		case gin.ReleaseMode:
			x.engine = release()
		default:
			x.engine = debug()
		}
	}
	if len(x.conf.TemplDir) > 0 {
		x.engine.LoadHTMLGlob(fmt.Sprintf("%s/*", x.conf.TemplDir))
	}
	return x.engine
}

// Start implements IService.
func (x *XHttp) Start() error {
	srv := &http.Server{
		Addr:         x.conf.Addr,
		Handler:      x.engine,
		WriteTimeout: time.Duration(x.conf.Timeout),
		ReadTimeout:  time.Duration(x.conf.Timeout),
	}
	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	srvSignal(srv)

	select {
	case err := <-errCh:
		return err
	case <-time.After(30 * time.Second):
		return errors.New("server start timeout")
	}
}

func srvSignal(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown error %s", err.Error())
	}
	log.Println("server exiting...")
}

func debug() *gin.Engine {
	return gin.Default()
}

func release() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	return gin.New()
}

func defaultConf(conf *HttpConf) {
	if conf.Timeout == 0 {
		conf.Timeout = 3000
	}
}

var _ IServer = (*XHttp)(nil)
