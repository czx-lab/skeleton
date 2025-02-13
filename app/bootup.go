package app

import (
	"czx/app/cron"
	"czx/app/event/listen"
	"czx/app/router"
	"czx/app/svc"
	"czx/internal/constants"
	"czx/internal/event"
	"czx/internal/server/xhttp"
	"czx/internal/xcron"
	"fmt"
	"log"
)

type Bootup struct {
	server xhttp.IServer
	svcCtx *svc.AppContext
}

func New(path string) constants.IService {
	bootup := &Bootup{}

	ctx := svc.NewAppContext(path)
	bootup.server = xhttp.New(ctx.Config.HttpConf).RegisterHandlers(router.New(ctx))
	bootup.svcCtx = ctx

	return bootup
}

// Start implements constants.IService.
func (b *Bootup) Start() error {
	b.service()

	fmt.Printf("Starting server at %s...\n", b.svcCtx.Config.HttpConf.Addr)

	if err := b.server.Start(); err != nil {
		log.Fatalf("error: app bootup %s", err.Error())
		return err
	}
	return nil
}

func (b *Bootup) service() {
	// events
	events := []event.IListen{new(listen.Foo)}
	if err := b.svcCtx.Event.Register(events...); err != nil {
		log.Fatalf("error: app event register %s", err.Error())
		return
	}

	// crons
	tasks := []xcron.ICron{new(cron.FooTask)}
	xcron.New().AddFunc(tasks...)

	// mq consumers
	// TODO::
}

var _ constants.IService = (*Bootup)(nil)
