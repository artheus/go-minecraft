package ctx

import (
	"context"
	"github.com/artheus/go-minecraft/core/types"
)

type Context struct {
	ctx     context.Context
	cancel  context.CancelFunc
	gameApp types.IGameApplication
}

func NewContext(gameApp types.IGameApplication) (appContext *Context, err error) {
	ctx, cancel := context.WithCancel(context.Background())

	appContext = &Context{
		ctx:     ctx,
		cancel:  cancel,
		gameApp: gameApp,
	}

	return appContext, nil
}

func (c *Context) Game() types.IGameApplication {
	return c.gameApp
}

func (c *Context) Context() context.Context {
	return c.ctx
}
