package core

import (
	"log"
	"time"
)

var (
	game *Game
)

func Run() {
	err := LoadTextureDesc()
	if err != nil {
		log.Fatal(err)
	}

	err = InitStore()
	if err != nil {
		log.Panic(err)
	}
	defer store.Close()

	err = InitClient()
	if err != nil {
		log.Panic(err)
	}
	if client != nil {
		defer client.Close()
	}

	game, err = NewGame(800, 600)
	if err != nil {
		log.Panic(err)
	}

	game.Camera().Restore(store.GetPlayerState())
	tick := time.Tick(time.Second / 60)
	for !game.ShouldClose() {
		<-tick
		game.Update()
	}

	if err = store.UpdatePlayerState(game.Camera().State()); err != nil {
		log.Panic(err)
	}
}
