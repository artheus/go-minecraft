package main

import (
	"flag"
	"github.com/artheus/go-minecraft/core"
	_ "image/png"
	"log"

	"net/http"
	_ "net/http/pprof"

	"github.com/faiface/mainthread"
)

var (
	pprofPort = flag.String("pprof", "", "http pprof port")
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	flag.Parse()
	go func() {
		if *pprofPort != "" {
			log.Fatal(http.ListenAndServe(*pprofPort, nil))
		}
	}()
	mainthread.Run(core.Run)
}
