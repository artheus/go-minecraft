package core

import (
	"github.com/artheus/go-minecraft/core/ctx"
	"github.com/artheus/go-minecraft/core/game"
	"github.com/artheus/go-minecraft/core/game/rpc"
	"github.com/artheus/go-minecraft/core/game/store"
	"github.com/artheus/go-minecraft/core/item"
	"log"
	"time"
)

func Run() {
	var err error
	var gameApp *game.Application

	err = item.LoadTextureDesc()
	if err != nil {
		log.Fatal(err)
	}

	err = store.InitStore()
	if err != nil {
		log.Panic(err)
	}
	defer store.Storage.Close()

	gameApp, err = game.NewGame(800, 600)
	if err != nil {
		log.Panic(err)
	}

	appCtx, err := ctx.NewContext(gameApp)
	if err != nil {
		log.Fatal(err)
	}

	err = gameApp.Init(appCtx)
	if err != nil {
		log.Panic(err)
	}

	err = rpc.InitClient(appCtx)
	if err != nil {
		log.Panic(err)
	}

	if rpc.Client != nil {
		defer rpc.Client.Close()
	}

	gameApp.Camera().Restore(store.Storage.GetPlayerState())
	tick := time.Tick(time.Second / 60)
	for !gameApp.ShouldClose() {
		<-tick
		gameApp.Update()
	}

	if err = store.Storage.UpdatePlayerState(gameApp.Camera().State()); err != nil {
		log.Panic(err)
	}
}
