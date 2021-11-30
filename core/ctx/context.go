package ctx

import (
	"context"
	"github.com/artheus/go-events"
	evttypes "github.com/artheus/go-events/types"
	"github.com/artheus/go-minecraft/core/types"
)

type Context struct {
	ctx          context.Context
	cancel       context.CancelFunc
	gameApp      types.IGameApplication
	eventPipe    evttypes.Pipe
}

const (
	eventPipeSize = 100
)

func NewContext(gameApp types.IGameApplication) (appContext *Context, err error) {
	ctx, cancel := context.WithCancel(context.Background())

	appContext = &Context{
		ctx:          ctx,
		cancel:       cancel,
		gameApp:      gameApp,
	}

	if appContext.eventPipe, err = events.Pipe(eventPipeSize); err != nil {
		appContext.Cancel()
		return nil, err
	}

	return appContext, nil
}

func (c *Context) Game() types.IGameApplication {
	return c.gameApp
}

func (c *Context) Context() context.Context {
	return c.ctx
}

func (c *Context) EventPipe() evttypes.Pipe {
	return c.eventPipe
}

func (c *Context) Cancel() {
	if c.eventPipe != nil {
		_ = c.eventPipe.Close()
	}

	c.cancel()
}